// Copyright © 2021 Banzai Cloud
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

package nodeagent

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
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
	containerName                  = "fluent-bit"
)

var NodeAgentFluentbitWindowsDefaults = &v1beta1.NodeAgentFluentbit{
	Flush: 1,
}
var NodeAgentFluentbitLinuxDefaults = &v1beta1.NodeAgentFluentbit{
	Flush: 2,
}

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func (n *nodeAgentInstance) getFluentBitLabels() map[string]string {
	return util.MergeLabels(n.logging.Spec.FluentbitSpec.Labels, map[string]string{
		"app.kubernetes.io/name": "fluentbit"}, generateLoggingRefLabels(n.logging.ObjectMeta.GetName()))
}

//func (r *Reconciler) getServiceAccount() string {
//	if r.Logging.Spec.FluentbitSpec.Security.ServiceAccount != "" {
//		return r.Logging.Spec.FluentbitSpec.Security.ServiceAccount
//	}
//	return r.Logging.QualifiedName(defaultServiceAccountName)
//}
//
//type DesiredObject struct {
//	Object runtime.Object
//	State  reconciler.DesiredState
//}
//
// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging *v1beta1.Logging
	*reconciler.GenericResourceReconciler
	configs map[string][]byte
}

// NewReconciler creates a new NodeAgent reconciler
func New(client client.Client, logger logr.Logger, logging *v1beta1.Logging, opts reconciler.ReconcilerOpts) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, logger, opts),
	}
}

type nodeAgentInstance struct {
	nodeAgent  *v1beta1.NodeAgent
	reconciler *reconciler.GenericResourceReconciler
	logging    *v1beta1.Logging
}

// Reconcile reconciles the NodeAgent resource
func (r *Reconciler) Reconcile() (*reconcile.Result, error) {
	for _, a := range r.Logging.Spec.NodeAgents {
		var instance nodeAgentInstance
		switch a.Type {
		case "windows":
			err := mergo.Merge(a, NodeAgentFluentbitWindowsDefaults)
			if err != nil {
				return nil, err
			}
			instance = nodeAgentInstance{
				nodeAgent:  a,
				reconciler: r.GenericResourceReconciler,
				logging:    r.Logging,
			}
		default:
			err := mergo.Merge(a, NodeAgentFluentbitLinuxDefaults)
			if err != nil {
				return nil, err
			}
			instance = nodeAgentInstance{
				nodeAgent:  a,
				reconciler: r.GenericResourceReconciler,
				logging:    r.Logging,
			}

		}

		result, err := instance.Reconcile()
		if err != nil {
			return nil, errors.WrapWithDetails(err,
				"failed to reconcile instances", "NodeName", a.Name)
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
}

// Reconcile reconciles the nodeAgentInstance resource
func (n *nodeAgentInstance) Reconcile() (*reconcile.Result, error) {
	for _, factory := range []resources.Resource{
		n.serviceAccount,
		//n.clusterRole,
		//n.clusterRoleBinding,
		//n.clusterPodSecurityPolicy,
		//n.pspClusterRole,
		//n.pspClusterRoleBinding,
		//n.configSecret,
		//n.daemonSet,
		//n.serviceMetrics,
		//n.monitorServiceMetrics,
	} {
		o, state, err := factory()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", factory)
		}
		result, err := n.reconciler.ReconcileResource(o, state)
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

//
//func RegisterWatches(builder *builder.Builder) *builder.Builder {
//	return builder.
//		Owns(&corev1.ConfigMap{}).
//		Owns(&appsv1.DaemonSet{}).
//		Owns(&rbacv1.ClusterRole{}).
//		Owns(&rbacv1.ClusterRoleBinding{}).
//		Owns(&corev1.ServiceAccount{})
//}
