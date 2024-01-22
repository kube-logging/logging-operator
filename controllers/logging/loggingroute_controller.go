// Copyright © 2023 Kube logging authors
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

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	apitypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentbit"
	loggingv1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func NewLoggingRouteReconciler(client client.Client, log logr.Logger) *LoggingRouteReconciler {
	return &LoggingRouteReconciler{
		Client: client,
		Log:    log,
	}
}

// LoggingRouteReconciler reconciles a LoggingRoute object
type LoggingRouteReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggingroutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggingroutes/status,verbs=get;update;patch

// Reconcile routes between logging domains
func (r *LoggingRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var loggingRoute loggingv1beta1.LoggingRoute
	if err := r.Client.Get(ctx, req.NamespacedName, &loggingRoute); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	tenants, err := fluentbit.FindTenants(ctx, loggingRoute.Spec.Targets, r.Client)
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "listing tenants")
	}

	var problems []string
	var notices []string
	loggingRoute.Status.Tenants = make([]loggingv1beta1.Tenant, 0)

	for _, t := range tenants {
		valid := true
		if t.AllNamespace {
			notices = append(notices, fmt.Sprintf("tenant %s receives logs from ALL namespaces", t.Name))
		} else if len(t.Namespaces) == 0 {
			problems = append(problems, fmt.Sprintf("tenant %s will be skipped as it does not provide valid target namespaces", t.Name))
			valid = false
		}
		tenantStatus := loggingv1beta1.Tenant{
			Name:       t.Name,
			Namespaces: t.Namespaces,
		}
		if valid {
			loggingRoute.Status.Tenants = append(loggingRoute.Status.Tenants, tenantStatus)
		}
	}

	loggingRoute.Status.Problems = problems
	loggingRoute.Status.ProblemsCount = len(problems)
	loggingRoute.Status.Notices = notices
	loggingRoute.Status.NoticesCount = len(notices)

	err = r.Status().Update(ctx, &loggingRoute)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func SetupLoggingRouteWithManager(mgr ctrl.Manager, logger logr.Logger) error {
	// In case we receive an update about a logging resource
	// we better notify all the logging routes to check if their target list has changed
	// rather than complicate the watch logic here.
	// The number and processing time of logging routes is not expected to cause issues.
	loggingRequestMapper := handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		var requests []reconcile.Request
		if _, ok := obj.(*loggingv1beta1.Logging); ok {
			var lrList loggingv1beta1.LoggingRouteList
			if err := mgr.GetClient().List(ctx, &lrList); err != nil {
				logger.Error(err, "failed to list logging route resources")
				return nil
			}
			for _, lr := range lrList.Items {
				requests = append(requests, reconcile.Request{NamespacedName: apitypes.NamespacedName{Name: lr.Name}})
			}
		}
		return requests
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.LoggingRoute{}).
		Watches(&loggingv1beta1.Logging{}, loggingRequestMapper).
		Complete(NewLoggingRouteReconciler(mgr.GetClient(), logger))
}
