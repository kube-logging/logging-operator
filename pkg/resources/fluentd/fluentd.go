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
	"fmt"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
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

	"github.com/kube-logging/logging-operator/pkg/resources"
	"github.com/kube-logging/logging-operator/pkg/resources/configcheck"
	"github.com/kube-logging/logging-operator/pkg/resources/kubetool"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	SecretConfigName        = "fluentd"
	AppSecretConfigName     = "fluentd-app"
	ConfigCheckKey          = "generated.conf"
	ConfigKey               = "fluent.conf"
	AppConfigKey            = "fluentd.conf"
	StatefulSetName         = "fluentd"
	ServiceName             = "fluentd"
	ServicePort             = 24240
	OutputSecretName        = "fluentd-output"
	OutputSecretPath        = "/fluentd/secret"
	PodDisruptionBudgetName = "fluentd"

	bufferPath                     = "/buffers"
	defaultServiceAccountName      = "fluentd"
	roleBindingName                = "fluentd"
	roleName                       = "fluentd"
	clusterRoleBindingName         = "fluentd"
	clusterRoleName                = "fluentd"
	containerName                  = "fluentd"
	defaultBufferVolumeMetricsPort = 9200
	drainerCheckInterval           = "10"
)

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging       *v1beta1.Logging
	fluentdSpec   *v1beta1.FluentdSpec
	fluentdConfig *v1beta1.FluentdConfig
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

func GetFluentd(ctx context.Context, Client client.Client, log logr.Logger, controlNamespace string) *v1beta1.FluentdConfig {
	fluentdList := v1beta1.FluentdConfigList{}
	// Detached fluentd must be in the `control namespace`
	nsOpt := client.InNamespace(controlNamespace)

	if err := Client.List(ctx, &fluentdList, nsOpt); err != nil {
		log.Error(err, "listing fluentd configuration")
		return nil
	}

	if len(fluentdList.Items) > 1 {
		log.Error(errors.New("multiple fluentd configurations found"), fmt.Sprintf("number of configurations: %d", len(fluentdList.Items)))
		return nil
	}

	if len(fluentdList.Items) == 1 {
		return &fluentdList.Items[0]
	}
	return nil
}

func (r *Reconciler) getServiceAccount() string {
	if r.fluentdSpec.Security.ServiceAccount != "" {
		return r.fluentdSpec.Security.ServiceAccount
	}
	return r.Logging.QualifiedName(defaultServiceAccountName)
}

func New(client client.Client, log logr.Logger,
	logging *v1beta1.Logging, fluentdSpec *v1beta1.FluentdSpec, fluentdConfig *v1beta1.FluentdConfig, config *string, secrets *secret.MountSecrets, opts reconciler.ReconcilerOpts) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		fluentdSpec:               fluentdSpec,
		fluentdConfig:             fluentdConfig,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, log, opts),
		config:                    config,
		secrets:                   secrets,
	}
}

// Reconcile reconciles the fluentd resource
func (r *Reconciler) Reconcile(ctx context.Context) (*reconcile.Result, error) {
	patchBase := client.MergeFrom(r.Logging.DeepCopy())

	objects := []resources.Resource{
		r.serviceAccount,
		r.role,
		r.roleBinding,
		r.clusterRole,
		r.clusterRoleBinding,
	}

	for _, res := range objects {
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
			if hasPod, err := r.hasConfigCheckPod(ctx, hash, *r.fluentdSpec); hasPod {
				return nil, errors.WrapIf(err, "current config is invalid")
			}
			// clean the status so that we can rerun the check
			return r.statusUpdate(ctx, patchBase, nil)
		}

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
		r.secretConfig,
		r.appConfigSecret,
		r.statefulset,
		r.service,
		r.headlessService,
		r.serviceMetrics,
		r.serviceBufferMetrics,
		r.pdb,
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

	if res, err := r.reconcileDrain(ctx); res != nil || err != nil {
		return res, err
	}

	return nil, nil
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

func (r *Reconciler) reconcileDrain(ctx context.Context) (*reconcile.Result, error) {
	if r.fluentdSpec.DisablePvc || !r.fluentdSpec.Scaling.Drain.Enabled {
		r.Log.Info("fluentd buffer draining is disabled")
		return nil, nil
	}

	nsOpt := client.InNamespace(r.Logging.Spec.ControlNamespace)
	fluentdLabelSet := r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec)

	var pvcList corev1.PersistentVolumeClaimList
	if err := r.Client.List(ctx, &pvcList, nsOpt,
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(fluentdLabelSet).Add(drainableRequirement),
		}); err != nil {
		return nil, errors.WrapIf(err, "listing PVC resources")
	}

	var stsPods corev1.PodList
	if err := r.Client.List(ctx, &stsPods, nsOpt, client.MatchingLabels(fluentdLabelSet)); err != nil {
		return nil, errors.WrapIf(err, "listing StatefulSet pods")
	}

	bufVolName := r.Logging.QualifiedName(r.fluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName)

	pvcsInUse := make(map[string]bool)
	for _, pod := range stsPods.Items {
		if bufVol := kubetool.FindVolumeByName(pod.Spec.Volumes, bufVolName); bufVol != nil {
			pvcsInUse[bufVol.PersistentVolumeClaim.ClaimName] = true
		}
	}

	replicaCount, err := NewDataProvider(r.Client, r.Logging, r.fluentdSpec, r.fluentdConfig).GetReplicaCount(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "get replica count for fluentd")
	}

	// mark PVCs required for upscaling as in-use
	for i := int32(0); i < utils.PointerToInt32(replicaCount); i++ {
		pvcsInUse[fmt.Sprintf("%s-%s-%d", bufVolName, r.Logging.QualifiedName(StatefulSetName), i)] = true
	}

	var jobList batchv1.JobList
	if err := r.Client.List(ctx, &jobList, nsOpt, client.MatchingLabels(r.Logging.GetFluentdLabels(ComponentDrainer, *r.fluentdSpec))); err != nil {
		return nil, errors.WrapIf(err, "listing buffer drainer jobs")
	}

	jobOfPVC := make(map[string]batchv1.Job)
	for _, job := range jobList.Items {
		if bufVol := kubetool.FindVolumeByName(job.Spec.Template.Spec.Volumes, bufVolName); bufVol != nil {
			jobOfPVC[bufVol.PersistentVolumeClaim.ClaimName] = job
		}
	}

	var cr reconciler.CombinedResult
	for _, pvc := range pvcList.Items {
		pvcLog := r.Log.WithValues("pvc", pvc.Name)

		drained := markedAsDrained(pvc)
		inUse := pvcsInUse[pvc.Name]
		if drained && inUse {
			pvcLog.Info("removing drained label from PVC as it has a matching statefulset pod")

			patch := client.MergeFrom(pvc.DeepCopy())
			delete(pvc.Labels, drainStatusLabelKey)
			if err := client.IgnoreNotFound(r.Client.Patch(ctx, pvc.DeepCopy(), patch)); err != nil {
				cr.CombineErr(errors.WrapIf(err, "removing drained label from pvc"))
			}
			continue
		}

		job, hasJob := jobOfPVC[pvc.Name]
		if hasJob && kubetool.JobSuccessfullyCompleted(&job) {
			pvcLog.Info("drainer job for PVC has completed, adding drained label and deleting job")

			patch := client.MergeFrom(pvc.DeepCopy())
			pvc.Labels[drainStatusLabelKey] = drainStatusLabelValue
			if err := client.IgnoreNotFound(r.Client.Patch(ctx, pvc.DeepCopy(), patch)); err != nil {
				cr.CombineErr(errors.WrapIf(err, "marking pvc as drained"))
				continue
			}

			if err := client.IgnoreNotFound(r.Client.Delete(ctx, &job, client.PropagationPolicy(v1.DeletePropagationBackground))); err != nil {
				cr.CombineErr(errors.WrapIf(err, "deleting completed drainer job"))
				continue
			}

			if r.fluentdSpec.Scaling.Drain.DeleteVolume {
				if err := client.IgnoreNotFound(r.Client.Delete(ctx, &pvc, client.PropagationPolicy(v1.DeletePropagationBackground))); err != nil {
					cr.CombineErr(errors.WrapIfWithDetails(err, "deleting drained PVC", "pvc", pvc.Name))
					continue
				}
			}

			if res, err := r.ReconcileResource(r.placeholderPodFor(pvc), reconciler.StateAbsent); err != nil {
				cr.Combine(res, errors.WrapIfWithDetails(err, "removing placeholder pod for pvc", "pvc", pvc.Name))
				continue
			}

			continue
		}

		if inUse && hasJob {
			pvcLog.Info("deleting drainer job early as PVC is now in use")

			if err := client.IgnoreNotFound(r.Client.Delete(ctx, &job, client.PropagationPolicy(v1.DeletePropagationForeground))); err != nil {
				cr.CombineErr(errors.WrapIf(err, "deleting unnecessary drainer job"))
				continue
			}

			if res, err := r.ReconcileResource(r.placeholderPodFor(pvc), reconciler.StateAbsent); err != nil {
				cr.Combine(res, errors.WrapIfWithDetails(err, "removing placeholder pod for pvc", "pvc", pvc.Name))
				continue
			}
			continue
		}

		if hasJob && !kubetool.JobSuccessfullyCompleted(&job) {
			if job.Status.Failed > 0 {
				cr.CombineErr(errors.NewWithDetails("draining PVC failed", "pvc", pvc.Name, "attempts", job.Status.Failed))
			} else {
				pvcLog.Info("drainer job for PVC has not yet been completed")
			}
			continue
		}

		if !drained && !inUse && !hasJob {
			pvcLog.Info("creating drainer job for PVC")

			if res, err := r.ReconcileResource(r.placeholderPodFor(pvc), reconciler.StatePresent); err != nil {
				cr.Combine(res, errors.WrapIfWithDetails(err, "ensuring placeholder pod is present for pvc", "pvc", pvc.Name))
				continue
			}

			if job, err := r.drainerJobFor(pvc, *r.fluentdSpec); err != nil {
				cr.CombineErr(errors.WrapIf(err, "assembling drainer job"))
			} else {
				cr.Combine(r.ReconcileResource(job, reconciler.StatePresent))
			}
			continue
		}
	}
	var res *reconcile.Result
	if !cr.Result.IsZero() {
		res = &cr.Result
	}
	return res, cr.Err
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
