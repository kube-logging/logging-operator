// Copyright © 2019 Banzai Cloud
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

package fluentd

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/secret"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	SecretConfigName      = "fluentd"
	AppSecretConfigName   = "fluentd-app"
	AppConfigKey          = "fluentd.conf"
	StatefulSetName       = "fluentd"
	PodSecurityPolicyName = "fluentd"
	ServiceName           = "fluentd"
	OutputSecretName      = "fluentd-output"
	OutputSecretPath      = "/fluentd/secret"

	bufferPath                = "/buffers"
	bufferVolumeName          = "fluentd-buffer"
	defaultServiceAccountName = "fluentd"
	roleBindingName           = "fluentd"
	roleName                  = "fluentd"
	containerName             = "fluentd"
)

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging *v1beta1.Logging
	*reconciler.GenericResourceReconciler
	config  *string
	secrets *secret.MountSecrets
}

type Desire struct {
	DesiredObject runtime.Object
	DesiredState  reconciler.DesiredState
	// BeforeUpdateHook has the ability to change the desired object
	// or even to change the desired state in case the object should be recreated
	BeforeUpdateHook func(runtime.Object) (reconciler.DesiredState, error)
}

func (r *Reconciler) getFluentdLabels(component string) map[string]string {
	return util.MergeLabels(
		r.Logging.Spec.FluentdSpec.Labels,
		map[string]string{
			"app.kubernetes.io/name":      "fluentd",
			"app.kubernetes.io/component": component,
		},
		generateLoggingRefLabels(r.Logging.ObjectMeta.GetName()),
	)
}

func (r *Reconciler) getServiceAccount() string {
	if r.Logging.Spec.FluentdSpec.Security.ServiceAccount != "" {
		return r.Logging.Spec.FluentdSpec.Security.ServiceAccount
	}
	return r.Logging.QualifiedName(defaultServiceAccountName)
}

func New(client client.Client, log logr.Logger,
	logging *v1beta1.Logging, config *string, secrets *secret.MountSecrets, opts reconciler.ReconcilerOpts) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, log, opts),
		config:                    config,
		secrets:                   secrets,
	}
}

// Reconcile reconciles the fluentd resource
func (r *Reconciler) Reconcile() (*reconcile.Result, error) {
	for _, res := range []resources.Resource{
		r.serviceAccount,
		r.role,
		r.roleBinding,
		r.clusterPodSecurityPolicy,
		r.pspRole,
		r.pspRoleBinding,
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
		if result, ok := r.Logging.Status.ConfigCheckResults[hash]; ok {
			// We already have an existing configcheck result:
			// - bail out if it was unsuccessful
			// - cleanup previous results if it's successful
			if !result {
				return nil, errors.Errorf("current config is invalid")
			}
			var removedHashes []string
			if removedHashes, err = r.configCheckCleanup(hash); err != nil {
				r.Log.Error(err, "failed to cleanup resources")
			} else {
				if len(removedHashes) > 0 {
					for _, removedHash := range removedHashes {
						delete(r.Logging.Status.ConfigCheckResults, removedHash)
					}
					if err := r.Client.Status().Update(context.TODO(), r.Logging); err != nil {
						return nil, errors.WrapWithDetails(err, "failed to update status", "logging", r.Logging)
					} else {
						// explicitly ask for a requeue to short circuit the controller loop after the status update
						return &reconcile.Result{Requeue: true}, nil
					}
				}
			}
		} else {
			// We don't have an existing result
			// - let's create what's necessary to have one
			// - if the result is ready write it into the status
			result, err := r.configCheck()
			if err != nil {
				return nil, errors.WrapIf(err, "failed to validate config")
			}
			if result.Ready {
				r.Logging.Status.ConfigCheckResults[hash] = result.Valid
				if err := r.Client.Status().Update(context.TODO(), r.Logging); err != nil {
					return nil, errors.WrapWithDetails(err, "failed to update status", "logging", r.Logging)
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
				return &reconcile.Result{RequeueAfter: time.Minute}, nil
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
	for _, res := range []resources.Resource{
		r.secretConfig,
		r.appconfigMap,
		r.statefulset,
		r.service,
		r.headlessService,
		r.serviceMetrics,
		r.monitorServiceMetrics,
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

	return nil, nil
}

func RegisterWatches(builder *builder.Builder) *builder.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{})
}
