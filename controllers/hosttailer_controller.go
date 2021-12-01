// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ht "github.com/banzaicloud/logging-operator/pkg/resources/hosttailer"
	loggingextensionsv1alpha1 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1alpha1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
)

const (
	// HostTailerReconcilerOptsHelpString is a generic help string for reconcileOpts
	HostTailerReconcilerOptsHelpString = "Object has to be recreated, but refusing to remove without explicitly being told so. " +
		"Use hosttailer.spec.enableRecreateWorkloadOnImmutableFieldChange to move on but make sure to understand the consequences. " +
		"As of rule, to avoid data loss, make sure to use a persistent volume for buffers, which is the default, unless explicitly disabled or configured differently."
)

// HostTailerReconciler reconciles a HostTailer object
type HostTailerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging-extensions.banzaicloud.io,resources=hosttailers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging-extensions.banzaicloud.io,resources=hosttailers/status,verbs=get;update;patch

// Reconcile .
func (r *HostTailerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("hosttailer", req.NamespacedName)

	// your logic here

	hosttailer := loggingextensionsv1alpha1.HostTailer{}

	if err := r.Client.Get(ctx, req.NamespacedName, &hosttailer); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	reconcilerOpts := reconciler.ReconcilerOpts{
		EnableRecreateWorkloadOnImmutableFieldChange:     hosttailer.Spec.EnableRecreateWorkloadOnImmutableFieldChange,
		EnableRecreateWorkloadOnImmutableFieldChangeHelp: HostTailerReconcilerOptsHelpString,
	}

	reconcilers := make([]reconciler.ComponentReconciler, 0)

	reconcilers = append(reconcilers, ht.New(r.Client, log, reconcilerOpts, hosttailer))

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
func (r *HostTailerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingextensionsv1alpha1.HostTailer{}).
		Complete(r)
}
