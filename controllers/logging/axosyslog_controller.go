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

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// NewAxoSyslogReconciler creates a new AxoSyslogReconciler instance
func NewAxoSyslogReconciler(client client.Client, log logr.Logger) *AxoSyslogReconciler {
	return &AxoSyslogReconciler{
		Client: client,
		Log:    log,
	}
}

type AxoSyslogReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=axosyslogs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=axosyslogs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;persistentvolumeclaims;serviceaccounts;pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconciles axosyslog resource
func (r *AxoSyslogReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("axosyslog", req.NamespacedName)

	log.V(1).Info("Reconciling AxoSyslog")

	var axosyslog v1beta1.AxoSyslog
	if err := r.Get(ctx, req.NamespacedName, &axosyslog); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := axosyslog.SetDefaults(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func SetupAxoSyslogWithManager(mgr ctrl.Manager, logger logr.Logger) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.AxoSyslog{}).
		Named("AxoSyslogController").
		Complete(NewAxoSyslogReconciler(mgr.GetClient(), logger))
}
