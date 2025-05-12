// Copyright © 2025 Kube logging authors
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
	"fmt"
	"reflect"
	"runtime"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/resources"
	axosyslogresources "github.com/kube-logging/logging-operator/pkg/resources/axosyslog"
	v1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// NewAxoSyslogReconciler creates a new AxoSyslogReconciler instance
func NewAxoSyslogReconciler(client client.Client, log logr.Logger, opts reconciler.ReconcilerOpts) *AxoSyslogReconciler {
	return &AxoSyslogReconciler{
		Client:                    client,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, log, opts),
		Log:                       log,
	}
}

type AxoSyslogReconciler struct {
	client.Client
	*reconciler.GenericResourceReconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=axosyslogs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=axosyslogs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;persistentvolumeclaims;serviceaccounts;pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconcile implements the reconciliation logic for AxoSyslog resources
func (r *AxoSyslogReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.V(1).Info("Reconciling AxoSyslog")

	var axoSyslog v1beta1.AxoSyslog
	if err := r.Get(ctx, req.NamespacedName, &axoSyslog); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log := r.Log.WithValues("axosyslog", fmt.Sprintf("%s/%s", axoSyslog.Namespace, axoSyslog.Name))

	if err := axoSyslog.SetDefaults(); err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "failed to set defaults")
	}

	// TODO: config-check ?

	if result, err := r.reconcileWorkloadResources(log, &axoSyslog); err != nil {
		return ctrl.Result{}, err
	} else if result != nil {
		return *result, nil
	}

	return ctrl.Result{}, nil
}

// reconcileWorkloadResources handles resources related to AxoSyslog and requires it's spec
func (r *AxoSyslogReconciler) reconcileWorkloadResources(log logr.Logger, axoSyslog *v1beta1.AxoSyslog) (*ctrl.Result, error) {
	resourceBuilders := []resources.ResourceWithSpec{
		axosyslogresources.CreateAxoSyslogConfig,
		axosyslogresources.StatefulSet,
		axosyslogresources.Service,
		axosyslogresources.HeadlessService,
		// TODO: service-metrics & buffer-metrics ?
		// axosyslogresources.ServiceMetrics,
		// axosyslogresources.ServiceBufferMetrics,
	}

	for _, buildObject := range resourceBuilders {
		builderName := getFunctionName(buildObject)
		log.V(2).Info("Processing resource", "builder", builderName)

		o, state, err := buildObject(axoSyslog)
		if err != nil {
			return nil, errors.WrapIff(err, "failed to build object with %s", builderName)
		}
		if o == nil {
			return nil, errors.Errorf("reconcile error: %s returned nil object", builderName)
		}

		metaObj, ok := o.(metav1.Object)
		if !ok {
			return nil, errors.Errorf("reconcile error: %s returned non-metav1.Object", builderName)
		}

		if metaObj.GetNamespace() == "" {
			return nil, errors.Errorf("reconcile error: %s returned resource without namespace set", builderName)
		}

		if err := ctrl.SetControllerReference(axoSyslog, metaObj, r.Scheme()); err != nil {
			return nil, errors.WrapIff(err, "failed to set controller reference for %s", metaObj.GetName())
		}

		result, err := r.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIff(err, "failed to reconcile resource %s/%s", metaObj.GetNamespace(), metaObj.GetName())
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func SetupAxoSyslogWithManager(mgr ctrl.Manager, logger logr.Logger) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.AxoSyslog{}).
		Named("axosyslog").
		Complete(NewAxoSyslogReconciler(mgr.GetClient(), logger, reconciler.ReconcilerOpts{}))
}

func getFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
