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

package controllers

import (
	"bytes"
	"context"
	"regexp"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentbit"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/resources/model"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/render"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	loggingv1alpha2 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
)

// LoggingReconciler reconciles a Logging object
type LoggingReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggings/status,verbs=get;update;patch

// Reconcile logging resources
func (r *LoggingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	log := r.Log.WithValues("logging", req.NamespacedName)

	logging := &loggingv1alpha2.Logging{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, logging)
	if err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logging, err = logging.SetDefaults()
	if err != nil {
		return reconcile.Result{}, err
	}

	fluentdConfig, secretList, err := r.clusterConfiguration(logging)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.V(1).Info("flow configuration", "config", fluentdConfig)

	reconcilers := make([]resources.ComponentReconciler, 0)

	reconcilerOpts := k8sutil.ReconcilerOpts{
		EnableRecreateWorkloadOnImmutableFieldChange: logging.Spec.EnableRecreateWorkloadOnImmutableFieldChange,
		EnableRecreateWorkloadOnImmutableFieldChangeHelp: "Object has to be recreated, but refusing to remove without explicitly being told so. " +
			"Use logging.spec.enableRecreateWorkloadOnImmutableFieldChange to move on but make sure to understand the consequences. " +
			"As of fluentd, to avoid data loss, make sure to use a persistent volume for buffers, which is the default, unless explicitly disabled or configured differently. " +
			"As of fluent-bit, to avoid duplicated logs, make sure to configure a hostPath volume for the positions through `logging.spec.fluentbit.spec.positiondb`. ",
	}

	if logging.Spec.FluentdSpec != nil {
		reconcilers = append(reconcilers,
			fluentd.New(r.Client, r.Log, logging, &fluentdConfig, secretList, reconcilerOpts).Reconcile)
	}

	if logging.Spec.FluentbitSpec != nil {
		reconcilers = append(reconcilers, fluentbit.New(r.Client, r.Log, logging, reconcilerOpts).Reconcile)
	}

	for _, rec := range reconcilers {
		result, err := rec()
		if err != nil {
			return reconcile.Result{}, err
		}
		if result != nil {
			// short circuit if requested explicitly
			return *result, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *LoggingReconciler) clusterConfiguration(logging *loggingv1alpha2.Logging) (string, *secret.MountSecrets, error) {
	if logging.Spec.FlowConfigOverride != "" {
		return logging.Spec.FlowConfigOverride, nil, nil
	}
	loggingResources, err := r.GetResources(logging)
	if err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to get logging resources", "logging", logging)
	}
	builder, err := loggingResources.CreateModel()
	if err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to create model", "logging", logging)
	}
	fluentConfig, err := builder.Build()
	if err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to build model", "logging", logging)
	}
	output := &bytes.Buffer{}
	renderer := render.FluentRender{
		Out:    output,
		Indent: 2,
	}
	err = renderer.Render(fluentConfig)
	if err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to render fluentd config", "logging", logging)
	}
	return output.String(), loggingResources.Secrets, nil
}

// SetupLoggingWithManager setup logging manager
func SetupLoggingWithManager(mgr ctrl.Manager, logger logr.Logger) *ctrl.Builder {
	clusterOutputSource := &source.Kind{Type: &loggingv1alpha2.ClusterOutput{}}
	clusterFlowSource := &source.Kind{Type: &loggingv1alpha2.ClusterFlow{}}
	outputSource := &source.Kind{Type: &loggingv1alpha2.Output{}}
	flowSource := &source.Kind{Type: &loggingv1alpha2.Flow{}}
	secretSource := &source.Kind{Type: &corev1.Secret{}}

	requestMapper := &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(mapObject handler.MapObject) []reconcile.Request {
			object, err := meta.Accessor(mapObject.Object)
			if err != nil {
				return nil
			}
			// get all the logging resources from the cache
			loggingList := &loggingv1alpha2.LoggingList{}
			err = mgr.GetCache().List(context.TODO(), loggingList)
			if err != nil {
				logger.Error(err, "failed to list logging resources")
				return nil
			}
			if o, ok := object.(*corev1.Secret); ok {
				requestList := []reconcile.Request{}
				for key := range o.Annotations {
					r := regexp.MustCompile("logging.banzaicloud.io/(.*)")
					result := r.FindStringSubmatch(key)
					if len(result) > 1 {
						loggingRef := result[1]
						// When loggingRef is default we trigger "empty" and default loggingRef as well, because we can't use an empty string in the annotation, thus we refer to `default` in case the loggingRef is empty
						if loggingRef == "default" {
							requestList = append(requestList, reconcileRequestsForLoggingRef(loggingList, loggingRef)...)
							loggingRef = ""
						}
						requestList = append(requestList, reconcileRequestsForLoggingRef(loggingList, loggingRef)...)
					}
				}
				return requestList
			}
			if o, ok := object.(*loggingv1alpha2.ClusterOutput); ok {
				return reconcileRequestsForLoggingRef(loggingList, o.Spec.LoggingRef)
			}
			if o, ok := object.(*loggingv1alpha2.Output); ok {
				return reconcileRequestsForLoggingRef(loggingList, o.Spec.LoggingRef)
			}
			if o, ok := object.(*loggingv1alpha2.Flow); ok {
				return reconcileRequestsForLoggingRef(loggingList, o.Spec.LoggingRef)
			}
			if o, ok := object.(*loggingv1alpha2.ClusterFlow); ok {
				return reconcileRequestsForLoggingRef(loggingList, o.Spec.LoggingRef)
			}
			return nil
		}),
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1alpha2.Logging{}).
		Owns(&corev1.Pod{}).
		Watches(clusterOutputSource, requestMapper).
		Watches(clusterFlowSource, requestMapper).
		Watches(outputSource, requestMapper).
		Watches(flowSource, requestMapper).
		Watches(secretSource, requestMapper)

	FluentdWatches(builder)
	FluentbitWatches(builder)

	return builder
}

func reconcileRequestsForLoggingRef(loggingList *loggingv1alpha2.LoggingList, loggingRef string) []reconcile.Request {
	filtered := make([]reconcile.Request, 0)
	for _, l := range loggingList.Items {
		if l.Spec.LoggingRef == loggingRef {
			filtered = append(filtered, reconcile.Request{
				NamespacedName: types.NamespacedName{
					// this happens to be empty as long as Logging is cluster scoped
					Namespace: l.Namespace,
					Name:      l.Name,
				},
			})
		}
	}
	return filtered
}

// FluentdWatches for fluentd statefulset
func FluentdWatches(builder *ctrl.Builder) *ctrl.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{})
}

// FluentbitWatches for fluent-bit daemonset
func FluentbitWatches(builder *ctrl.Builder) *ctrl.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{})
}

// GetResources collect all resources referenced by logging resource
func (r *LoggingReconciler) GetResources(logging *loggingv1alpha2.Logging) (*model.LoggingResources, error) {
	loggingResources := model.NewLoggingResources(logging, r.Client, r.Log)
	var err error

	clusterFlows := &loggingv1alpha2.ClusterFlowList{}
	err = r.List(context.TODO(), clusterFlows, client.InNamespace(logging.Spec.ControlNamespace))
	if err != nil {
		return nil, err
	}
	if len(clusterFlows.Items) > 0 {
		for _, i := range clusterFlows.Items {
			if i.Spec.LoggingRef == logging.Spec.LoggingRef {
				loggingResources.ClusterFlows = append(loggingResources.ClusterFlows, i)
			}
		}
	}

	clusterOutputs := &loggingv1alpha2.ClusterOutputList{}
	err = r.List(context.TODO(), clusterOutputs, client.InNamespace(logging.Spec.ControlNamespace))
	if err != nil {
		return nil, err
	}
	if len(clusterOutputs.Items) > 0 {
		for _, i := range clusterOutputs.Items {
			if i.Spec.LoggingRef == logging.Spec.LoggingRef {
				loggingResources.ClusterOutputs = append(loggingResources.ClusterOutputs, i)
			}
		}
	}

	watchNamespaces := logging.Spec.WatchNamespaces

	if len(watchNamespaces) == 0 {
		nsList := &corev1.NamespaceList{}
		err = r.List(context.TODO(), nsList)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to list all namespaces")
		}
		for _, ns := range nsList.Items {
			watchNamespaces = append(watchNamespaces, ns.Name)
		}
	}

	for _, ns := range watchNamespaces {
		flows := &loggingv1alpha2.FlowList{}
		err = r.List(context.TODO(), flows, client.InNamespace(ns))
		if err != nil {
			return nil, err
		}
		if len(flows.Items) > 0 {
			for _, i := range flows.Items {
				if i.Spec.LoggingRef == logging.Spec.LoggingRef {
					loggingResources.Flows = append(loggingResources.Flows, i)
				}
			}
		}
		outputs := &loggingv1alpha2.OutputList{}
		err = r.List(context.TODO(), outputs, client.InNamespace(ns))
		if err != nil {
			return nil, err
		}
		if len(outputs.Items) > 0 {
			for _, i := range outputs.Items {
				if i.Spec.LoggingRef == logging.Spec.LoggingRef {
					loggingResources.Outputs = append(loggingResources.Outputs, i)
				}
			}
		}
	}

	return loggingResources, nil
}
