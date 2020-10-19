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
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentbit"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/resources/model"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/render"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/secret"
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

	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
)

// NewLoggingReconciler returns a new LoggingReconciler instance
func NewLoggingReconciler(client client.Client, log logr.Logger) *LoggingReconciler {
	return &LoggingReconciler{
		Client: client,
		Log:    log,
	}
}

// LoggingReconciler reconciles a Logging object
type LoggingReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggings;flows;clusterflows;outputs;clusteroutputs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;policy,resources=podsecuritypolicies,verbs=get;list;watch;create;update;patch;delete;use
// +kubebuilder:rbac:groups=apps,resources=statefulsets;daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;persistentvolumeclaims;serviceaccounts;pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=nodes;namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete

// Reconcile logging resources
func (r *LoggingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	log := r.Log.WithValues("logging", req.NamespacedName)

	logging := new(loggingv1beta1.Logging)
	if err := r.Client.Get(ctx, req.NamespacedName, logging); err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logging, err := logging.SetDefaults()
	if err != nil {
		return reconcile.Result{}, err
	}

	fluentdConfig, secretList, err := r.clusterConfiguration(logging)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.V(1).Info("flow configuration", "config", fluentdConfig)

	var reconcilers []resources.ComponentReconciler

	reconcilerOpts := reconciler.ReconcilerOpts{
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

func (r *LoggingReconciler) clusterConfiguration(logging *loggingv1beta1.Logging) (string, *secret.MountSecrets, error) {
	if logging.Spec.FlowConfigOverride != "" {
		return logging.Spec.FlowConfigOverride, nil, nil
	}
	loggingResources, err := model.NewLoggingResourceRepository(r.Client).LoggingResourcesFor(context.TODO(), *logging)
	if err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to get logging resources", "logging", logging)
	}
	slf := secretLoaderFactory{
		Client: r.Client,
	}
	fluentConfig, err := model.CreateSystem(loggingResources, &slf, r.Log)
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
	return output.String(), &slf.Secrets, nil
}

type secretLoaderFactory struct {
	Client  client.Client
	Secrets secret.MountSecrets
}

func (f *secretLoaderFactory) OutputSecretLoaderForNamespace(namespace string) secret.SecretLoader {
	return secret.NewSecretLoader(f.Client, namespace, fluentd.OutputSecretPath, &f.Secrets)
}

// SetupLoggingWithManager setup logging manager
func SetupLoggingWithManager(mgr ctrl.Manager, logger logr.Logger) *ctrl.Builder {
	clusterOutputSource := &source.Kind{Type: &loggingv1beta1.ClusterOutput{}}
	clusterFlowSource := &source.Kind{Type: &loggingv1beta1.ClusterFlow{}}
	outputSource := &source.Kind{Type: &loggingv1beta1.Output{}}
	flowSource := &source.Kind{Type: &loggingv1beta1.Flow{}}
	secretSource := &source.Kind{Type: &corev1.Secret{}}

	requestMapper := &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(mapObject handler.MapObject) []reconcile.Request {
			object, err := meta.Accessor(mapObject.Object)
			if err != nil {
				return nil
			}
			// get all the logging resources from the cache
			var loggingList loggingv1beta1.LoggingList
			if err := mgr.GetCache().List(context.TODO(), &loggingList); err != nil {
				logger.Error(err, "failed to list logging resources")
				return nil
			}
			if o, ok := object.(*corev1.Secret); ok {
				var requestList []reconcile.Request
				r := regexp.MustCompile("logging.banzaicloud.io/(.*)")
				for key := range o.Annotations {
					result := r.FindStringSubmatch(key)
					if len(result) > 1 {
						loggingRef := result[1]
						// When loggingRef is "default" we also trigger for the empty ("") loggingRef as well, because the empty string cannot be used in the annotation, thus "default" refers to the empty case.
						if loggingRef == "default" {
							requestList = append(requestList, reconcileRequestsForLoggingRef(loggingList.Items, "")...)
						}
						requestList = append(requestList, reconcileRequestsForLoggingRef(loggingList.Items, loggingRef)...)
					}
				}
				return requestList
			}
			switch o := object.(type) {
			case *loggingv1beta1.ClusterOutput:
				return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
			case *loggingv1beta1.Output:
				return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
			case *loggingv1beta1.Flow:
				return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
			case *loggingv1beta1.ClusterFlow:
				return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
			}
			return nil
		}),
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Logging{}).
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

func reconcileRequestsForLoggingRef(loggings []loggingv1beta1.Logging, loggingRef string) (reqs []reconcile.Request) {
	for _, l := range loggings {
		if l.Spec.LoggingRef == loggingRef {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: l.Namespace, // this happens to be empty as long as Logging is cluster scoped
					Name:      l.Name,
				},
			})
		}
	}
	return
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
