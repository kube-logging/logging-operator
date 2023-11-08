// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (n *nodeAgentInstance) clusterRole() (runtime.Object, reconciler.DesiredState, error) {
	if *n.nodeAgent.FluentbitSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRole{
			ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleName),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"pods", "namespaces"},
					Verbs:     []string{"get", "list", "watch"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRole{
		ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleName),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (n *nodeAgentInstance) clusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *n.nodeAgent.FluentbitSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleBindingName),
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     n.QualifiedName(clusterRoleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      n.getServiceAccount(),
					Namespace: n.logging.Spec.ControlNamespace,
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleBindingName),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}

func (n *nodeAgentInstance) serviceAccount() (runtime.Object, reconciler.DesiredState, error) {
	if *n.nodeAgent.FluentbitSpec.Security.RoleBasedAccessControlCreate && n.nodeAgent.FluentbitSpec.Security.ServiceAccount == "" {
		desired := &corev1.ServiceAccount{
			ObjectMeta: n.nodeAgent.Metadata.Merge(n.NodeAgentObjectMeta(defaultServiceAccountName)),
		}
		err := merge.Merge(desired, n.nodeAgent.FluentbitSpec.ServiceAccountOverrides)
		if err != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
		}

		return desired, reconciler.StatePresent, nil
	} else {
		desired := &corev1.ServiceAccount{
			ObjectMeta: n.NodeAgentObjectMeta(defaultServiceAccountName),
		}

		err := merge.Merge(desired, n.nodeAgent.FluentbitSpec.ServiceAccountOverrides)
		if err != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
		}
		return desired, reconciler.StateAbsent, nil
	}
}
