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
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/kube-logging/logging-operator/pkg/resources/eventtailer"
	loggingextensionsv1alpha1 "github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/v1alpha1"
)

// EventTailerReconciler reconciles a EventTailer object
type EventTailerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging-extensions.banzaicloud.io,resources=eventtailers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging-extensions.banzaicloud.io,resources=eventtailers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extensions;apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="events.k8s.io",resources=events,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=*

// Reconcile .
func (r *EventTailerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("eventtailer", req.NamespacedName)

	// your logic here

	eventTailer := loggingextensionsv1alpha1.EventTailer{}

	if err := r.Client.Get(ctx, req.NamespacedName, &eventTailer); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	reconcilerOpts := reconciler.ReconcilerOpts{
		EnableRecreateWorkloadOnImmutableFieldChange: true,
	}

	reconcilers := make([]reconciler.ComponentReconciler, 0)

	reconcilers = append(reconcilers, eventtailer.New(r.Client, log, reconcilerOpts, eventTailer))

	for _, rec := range reconcilers {
		result, err := rec.Reconcile(nil)
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

// SetupWithManager .
func (r *EventTailerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingextensionsv1alpha1.EventTailer{}).
		Complete(r)
}
