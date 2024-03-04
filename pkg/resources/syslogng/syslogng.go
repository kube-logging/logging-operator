// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syslogng

import (
	"context"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources"
	"github.com/kube-logging/logging-operator/pkg/resources/configcheck"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	ServiceName                       = "syslog-ng"
	ServicePort                       = 601
	configSecretName                  = "syslog-ng"
	configKey                         = "syslog-ng.conf"
	StatefulSetName                   = "syslog-ng"
	outputSecretName                  = "syslog-ng-output"
	OutputSecretPath                  = "/etc/syslog-ng/secret"
	BufferPath                        = "/buffers"
	serviceAccountName                = "syslog-ng"
	roleBindingName                   = "syslog-ng"
	roleName                          = "syslog-ng"
	clusterRoleBindingName            = "syslog-ng"
	clusterRoleName                   = "syslog-ng"
	ContainerName                     = "syslog-ng"
	defaultBufferVolumeMetricsPort    = 9200
	syslogngImageRepository           = "ghcr.io/axoflow/axosyslog"
	syslogngImageTag                  = "4.5.0"
	prometheusExporterImageRepository = "ghcr.io/axoflow/axosyslog-metrics-exporter"
	prometheusExporterImageTag        = "0.0.2"
	bufferVolumeImageRepository       = "ghcr.io/kube-logging/node-exporter"
	bufferVolumeImageTag              = "v0.7.1"
	configReloaderImageRepository     = "ghcr.io/kube-logging/syslogng-reload"
	configReloaderImageTag            = "v1.3.1"
	socketVolumeName                  = "socket"
	socketPath                        = "/tmp/syslog-ng/syslog-ng.ctl"
	configDir                         = "/etc/syslog-ng/config"
	configVolumeName                  = "config"
	tlsVolumeName                     = "tls"
	metricsPortNumber                 = 9577
	metricsPortName                   = "exporter"
)

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging        *v1beta1.Logging
	syslogNGSpec   *v1beta1.SyslogNGSpec
	syslogNGConfig *v1beta1.SyslogNGConfig
	*reconciler.GenericResourceReconciler
	config  string
	secrets *secret.MountSecrets
}

type Desire struct {
	DesiredObject runtime.Object
	DesiredState  reconciler.DesiredState
	// BeforeUpdateHook has the ability to change the desired object
	// or even to change the desired state in case the object should be recreated
	BeforeUpdateHook func(runtime.Object) (reconciler.DesiredState, error)
}

func New(
	client client.Client,
	log logr.Logger,
	logging *v1beta1.Logging,
	syslogNGSPec *v1beta1.SyslogNGSpec,
	syslogNGCOnfig *v1beta1.SyslogNGConfig,
	config string,
	secrets *secret.MountSecrets,
	opts reconciler.ReconcilerOpts,
) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		syslogNGSpec:              syslogNGSPec,
		syslogNGConfig:            syslogNGCOnfig,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, log, opts),
		config:                    config,
		secrets:                   secrets,
	}
}

// Reconcile reconciles the syslog-ng resource
func (r *Reconciler) Reconcile(ctx context.Context) (*reconcile.Result, error) {
	patchBase := client.MergeFrom(r.Logging.DeepCopy())

	for _, res := range []resources.Resource{
		r.serviceAccount,
		r.role,
		r.roleBinding,
		r.clusterRole,
		r.clusterRoleBinding,
	} {
		o, state, err := res()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", res)
		}
		result, err := r.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
		if result != nil {
			return result, nil
		}
	}
	// Config check and cleanup if enabled
	if !r.Logging.Spec.FlowConfigCheckDisabled { //nolint:nestif
		hash, err := r.configHash()
		if err != nil {
			return nil, err
		}

		// Fail when the current config is invalid
		if result, ok := r.Logging.Status.ConfigCheckResults[hash]; ok && !result {
			if hasPod, err := r.hasConfigCheckPod(ctx, hash); hasPod {
				return nil, errors.WrapIf(err, "current config is invalid")
			}
			// clean the status so that we can rerun the check
			return r.statusUpdate(ctx, patchBase, nil)
		}

		// Cleanup previous configcheck results
		if result, ok := r.Logging.Status.ConfigCheckResults[hash]; ok {
			cleaner := configcheck.NewConfigCheckCleaner(r.Client, ComponentConfigCheck)

			var cleanupErrs error
			cleanupErrs = errors.Append(cleanupErrs, cleaner.SecretCleanup(ctx, hash))
			cleanupErrs = errors.Append(cleanupErrs, cleaner.PodCleanup(ctx, hash))

			if cleanupErrs != nil {
				// Errors with the cleanup should not block the reconciliation, we just note it
				r.Log.Error(err, "issues during configcheck cleanup, moving on")
			} else if len(r.Logging.Status.ConfigCheckResults) > 1 {
				return r.statusUpdate(ctx, patchBase, map[string]bool{
					hash: result,
				})
			}
		} else {
			// We don't have an existing result
			// - let's create what's necessary to have one
			// - if the result is ready write it into the status
			result, err := r.configCheck(ctx)
			if err != nil {
				return nil, errors.WrapIf(err, "failed to validate config")
			}
			if result.Ready {
				r.Logging.Status.ConfigCheckResults[hash] = result.Valid
				if err := r.Client.Status().Patch(ctx, r.Logging, patchBase); err != nil {
					return nil, errors.WrapWithDetails(err, "failed to patch status", "logging", r.Logging)
				} else {
					// explicitly ask for a requeue to short circuit the controller loop after the status update
					return &reconcile.Result{Requeue: true}, nil
				}
			} else {
				if result.Message != "" {
					r.Log.Info(result.Message)
				} else {
					r.Log.Info("still waiting for the configcheck result...")
				}
				return &reconcile.Result{Requeue: true}, nil
			}
		}
	}
	// Prepare output secret
	outputSecret, outputSecretDesiredState, err := r.outputSecret(r.secrets, OutputSecretPath)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create output secret")
	}
	result, err := r.ReconcileResource(outputSecret, outputSecretDesiredState)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to reconcile resource")
	}
	if result != nil {
		return result, nil
	}
	// Mark watched secrets
	secretList, state, err := r.markSecrets(r.secrets)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to mark secrets")
	}
	for _, obj := range secretList {
		result, err := r.ReconcileResource(obj, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
		if result != nil {
			return result, nil
		}
	}
	resourceObjects := []resources.Resource{
		r.configSecret,
		r.statefulset,
		r.service,
		r.headlessService,
		r.serviceMetrics,
		r.serviceBufferMetrics,
	}
	if resources.IsSupported(ctx, resources.ServiceMonitorKey) {
		resourceObjects = append(resourceObjects, r.monitorServiceMetrics, r.monitorBufferServiceMetrics)
	}
	if resources.IsSupported(ctx, resources.PrometheusRuleKey) {
		resourceObjects = append(resourceObjects, r.prometheusRules, r.bufferVolumePrometheusRules)
	}
	for _, res := range resourceObjects {
		o, state, err := res()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", res)
		}
		result, err := r.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (r *Reconciler) hasConfigCheckPod(ctx context.Context, hashKey string) (bool, error) {
	var err error
	pod, err := r.newCheckPod(hashKey)
	if err != nil {
		return false, err
	}

	p := &corev1.Pod{}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(pod), p)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Reconciler) statusUpdate(ctx context.Context, patchBase client.Patch, result map[string]bool) (*reconcile.Result, error) {
	r.Logging.Status.ConfigCheckResults = result
	if err := r.Client.Status().Patch(ctx, r.Logging, patchBase); err != nil {
		return nil, errors.WrapWithDetails(err, "failed to patch status", "logging", r.Logging)
	} else {
		// explicitly ask for a requeue to short circuit the controller loop after the status update
		return &reconcile.Result{Requeue: true}, nil
	}
}

func (r *Reconciler) getServiceAccountName() string {
	return r.Logging.QualifiedName(serviceAccountName)
}

func RegisterWatches(builder *builder.Builder) *builder.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.PersistentVolumeClaim{})
}
