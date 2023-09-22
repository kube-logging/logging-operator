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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentbit"
	loggingv1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func NewAggregationPolicyReconciler(client client.Client, log logr.Logger) *AggregationPolicyReconciler {
	return &AggregationPolicyReconciler{
		Client: client,
		Log:    log,
	}
}

// AggregationPolicyReconciler reconciles an AggregationPolicy object
type AggregationPolicyReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=aggregationpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=aggregationpolicies/status,verbs=get;update;patch

// Reconcile aggregation policies
func (r *AggregationPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var agrPol loggingv1beta1.AggregationPolicy
	if err := r.Client.Get(ctx, req.NamespacedName, &agrPol); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	tenants, err := fluentbit.FindTenants(ctx, agrPol.Spec.WatchNamespaceTargets, r.Client)
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "listing tenants")
	}

	var problems []string
	agrPol.Status.Tenants = make([]loggingv1beta1.Tenant, 0)

	for _, t := range tenants {
		valid := true
		if t.AllNamespace {
			problems = append(problems, fmt.Sprintf("tenant %s receives logs from ALL namespaces", t.Name))
		} else {
			if len(t.Namespaces) == 0 {
				problems = append(problems, fmt.Sprintf("tenant %s will be skipped as it does not provide valid target namespaces", t.Name))
				valid = false
			}
		}
		tenantStatus := loggingv1beta1.Tenant{
			Name:       t.Name,
			Namespaces: t.Namespaces,
		}
		if valid {
			agrPol.Status.Tenants = append(agrPol.Status.Tenants, tenantStatus)
		}
	}

	agrPol.Status.Problems = problems

	err = r.Status().Update(ctx, &agrPol)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
