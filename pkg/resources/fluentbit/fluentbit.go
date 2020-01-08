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

package fluentbit

import (
	"emperror.dev/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	defaultServiceAccountName      = "fluentbit"
	clusterRoleBindingName         = "fluentbit"
	clusterRoleName                = "fluentbit"
	fluentBitSecretConfigName      = "fluentbit"
	fluentbitDaemonSetName         = "fluentbit"
	fluentbitPodSecurityPolicyName = "fluentbit"
	fluentbitServiceName           = "fluentbit"
)

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func (r *Reconciler) getFluentBitLabels() map[string]string {
	return util.MergeLabels(r.Logging.Spec.FluentbitSpec.Labels, map[string]string{
		"app.kubernetes.io/name": "fluentbit"}, generateLoggingRefLabels(r.Logging.ObjectMeta.GetName()))
}

func (r *Reconciler) getServiceAccount() string {
	if r.Logging.Spec.FluentbitSpec.Security.ServiceAccount != "" {
		return r.Logging.Spec.FluentbitSpec.Security.ServiceAccount
	}
	return r.Logging.QualifiedName(defaultServiceAccountName)
}

type DesiredObject struct {
	Object runtime.Object
	State  k8sutil.DesiredState
}

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging *v1beta1.Logging
	*k8sutil.GenericResourceReconciler
	desiredConfig string
}

// NewReconciler creates a new Fluentbit reconciler
func New(client client.Client, logger logr.Logger, logging *v1beta1.Logging, opts k8sutil.ReconcilerOpts) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		GenericResourceReconciler: k8sutil.NewReconciler(client, logger, opts),
	}
}

// Reconcile reconciles the fluentBit resource
func (r *Reconciler) Reconcile() (*reconcile.Result, error) {
	for _, factory := range []resources.Resource{
		r.serviceAccount,
		r.clusterRole,
		r.clusterRoleBinding,
		r.clusterPodSecurityPolicy,
		r.pspClusterRole,
		r.pspClusterRoleBinding,
		r.configSecret,
		r.daemonSet,
		r.serviceMetrics,
		r.monitorServiceMetrics,
	} {
		o, state, err := factory()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", factory)
		}
		result, err := r.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapWithDetails(err,
				"failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}
