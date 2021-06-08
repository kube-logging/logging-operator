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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	SecretConfigName      = "fluentd"
	AppSecretConfigName   = "fluentd-app"
	ConfigCheckKey        = "generated.conf"
	ConfigKey             = "fluent.conf"
	AppConfigKey          = "fluentd.conf"
	StatefulSetName       = "fluentd"
	PodSecurityPolicyName = "fluentd"
	ServiceName           = "fluentd"
	OutputSecretName      = "fluentd-output"
	OutputSecretPath      = "/fluentd/secret"

	bufferPath                     = "/buffers"
	defaultServiceAccountName      = "fluentd"
	roleBindingName                = "fluentd"
	roleName                       = "fluentd"
	clusterRoleBindingName         = "fluentd"
	clusterRoleName                = "fluentd"
	containerName                  = "fluentd"
	defaultBufferVolumeMetricsPort = 9200
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
	ctx := context.Background()

	for _, res := range []resources.Resource{
		r.serviceAccount,
		r.role,
		r.roleBinding,
		r.clusterRole,
		r.clusterRoleBinding,
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
					if err := r.Client.Status().Update(ctx, r.Logging); err != nil {
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
				if err := r.Client.Status().Update(ctx, r.Logging); err != nil {
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
		r.appConfigSecret,
		r.statefulset,
		r.service,
		r.headlessService,
		r.serviceMetrics,
		r.monitorServiceMetrics,
		r.serviceBufferMetrics,
		r.monitorBufferServiceMetrics,
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

	if !r.Logging.Spec.FluentdSpec.DisablePvc {
		nsOpt := client.InNamespace(r.Logging.Spec.ControlNamespace)
		labelSet := r.getFluentdLabels(ComponentFluentd)

		var pvcList corev1.PersistentVolumeClaimList
		if err := r.Client.List(ctx, &pvcList, nsOpt,
			client.MatchingLabelsSelector{
				Selector: labels.SelectorFromSet(labelSet).Add(drainableRequirement),
			}); err != nil {
			return nil, errors.WrapIf(err, "listing PVC resources")
		}

		var stsPods corev1.PodList
		if err := r.Client.List(ctx, &stsPods, nsOpt, client.MatchingLabels(labelSet)); err != nil {
			return nil, errors.WrapIf(err, "listing StatefulSet pods")
		}

		bufVolName := r.Logging.QualifiedName(r.Logging.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName)

		livePVCs := make(map[string]bool)
		for _, pod := range stsPods.Items {
			if bufVol := findVolumeByName(pod.Spec.Volumes, bufVolName); bufVol != nil {
				livePVCs[bufVol.PersistentVolumeClaim.ClaimName] = true
			}
		}

		var jobList batchv1.JobList
		if err := r.Client.List(ctx, &jobList, nsOpt, client.MatchingLabels(labelSet)); err != nil {
			return nil, errors.WrapIf(err, "listing buffer drain jobs")
		}

		jobOfPVC := make(map[string]batchv1.Job)
		for _, job := range jobList.Items {
			if bufVol := findVolumeByName(job.Spec.Template.Spec.Volumes, bufVolName); bufVol != nil {
				jobOfPVC[bufVol.PersistentVolumeClaim.ClaimName] = job
			}
		}

		var errs error
		for _, pvc := range pvcList.Items {
			drained := markedAsDrained(pvc)
			live := livePVCs[pvc.Name]
			if drained && live {
				patch := client.MergeFrom(pvc.DeepCopy())
				delete(pvc.Labels, drainStatusLabelKey)
				if err := client.IgnoreNotFound(r.Client.Patch(ctx, &pvc, patch)); err != nil {
					errs = errors.Append(errs, errors.WrapIf(err, "removing drained label from pvc"))
				}
				continue
			}
			job, hasJob := jobOfPVC[pvc.Name]
			if hasJob && jobSuccessfullyCompleted(job) {
				patch := client.MergeFrom(pvc.DeepCopy())
				pvc.Labels[drainStatusLabelKey] = drainStatusLabelValue
				if err := client.IgnoreNotFound(r.Client.Patch(ctx, &pvc, patch)); err != nil {
					errs = errors.Append(errs, errors.WrapIf(err, "marking pvc as drained"))
				}

				if err := client.IgnoreNotFound(r.Client.Delete(ctx, &job, client.PropagationPolicy(v1.DeletePropagationBackground))); err != nil {
					errs = errors.Append(errs, errors.WrapIf(err, "deleting completed drain job"))
				}
				continue
			}
			if !drained && !live && !hasJob {
				if job, err := r.drainJobFor(pvc); err != nil {
					errs = errors.Append(errs, errors.WrapIf(err, "assembling drain job"))
				} else if err := r.Client.Create(ctx, job); err != nil {
					errs = errors.Append(errs, errors.WrapIf(err, "creating drain job"))
				}
				continue
			}
		}
		if errs != nil {
			return nil, errs
		}
	}

	return nil, nil
}

func RegisterWatches(builder *builder.Builder) *builder.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&batchv1.Job{})
}

var drainableRequirement = requirementMust(labels.NewRequirement("logging.banzaicloud.io/drain", selection.NotEquals, []string{"no"}))

func requirementMust(req *labels.Requirement, err error) labels.Requirement {
	if err != nil {
		panic(err)
	}
	if req == nil {
		panic("requirement is nil")
	}
	return *req
}

const drainStatusLabelKey = "logging.banzaicloud.io/drain-status"
const drainStatusLabelValue = "drained"

func markedAsDrained(pvc corev1.PersistentVolumeClaim) bool {
	return pvc.Labels[drainStatusLabelKey] == drainStatusLabelValue
}

func findVolumeByName(vols []corev1.Volume, name string) *corev1.Volume {
	for i := range vols {
		vol := &vols[i]
		if vol.Name == name {
			return vol
		}
	}
	return nil
}

func jobSuccessfullyCompleted(job batchv1.Job) bool {
	return job.Status.CompletionTime != nil && job.Status.Succeeded > 0
}
