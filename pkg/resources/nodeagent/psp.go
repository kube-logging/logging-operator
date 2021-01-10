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
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

func (n *nodeAgentInstance) clusterPodSecurityPolicy() (runtime.Object, reconciler.DesiredState, error) {
	if n.nodeAgent.FluentbitSpec.Security.PodSecurityPolicyCreate {
		allowedHostPaths := []policyv1beta1.AllowedHostPath{{
			PathPrefix: n.nodeAgent.FluentbitSpec.MountPath,
			ReadOnly:   true,
		}, {
			PathPrefix: "/var/log",
			ReadOnly:   true,
		}}

		for _, vMnt := range n.nodeAgent.FluentbitSpec.ExtraVolumeMounts {
			allowedHostPaths = append(allowedHostPaths, policyv1beta1.AllowedHostPath{
				PathPrefix: vMnt.Source,
				ReadOnly:   vMnt.ReadOnly,
			})
		}

		if n.nodeAgent.FluentbitSpec.PositionDB.HostPath != nil {
			n.nodeAgent.FluentbitSpec.PositionDB.WithDefaultHostPath(
				fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, TailPositionVolume))

			allowedHostPaths = append(allowedHostPaths, policyv1beta1.AllowedHostPath{
				PathPrefix: n.nodeAgent.FluentbitSpec.PositionDB.HostPath.Path,
				ReadOnly:   false,
			})
		}

		return &policyv1beta1.PodSecurityPolicy{
			ObjectMeta: n.NodeAgentObjectMetaClusterScope(fluentbitPodSecurityPolicyName),
			Spec: policyv1beta1.PodSecurityPolicySpec{
				Volumes: []policyv1beta1.FSType{
					"configMap",
					"emptyDir",
					"secret",
					"hostPath"},
				SELinux: policyv1beta1.SELinuxStrategyOptions{
					Rule: policyv1beta1.SELinuxStrategyRunAsAny,
				},
				RunAsUser: policyv1beta1.RunAsUserStrategyOptions{
					Rule: policyv1beta1.RunAsUserStrategyRunAsAny,
				},
				SupplementalGroups: policyv1beta1.SupplementalGroupsStrategyOptions{
					Rule: policyv1beta1.SupplementalGroupsStrategyRunAsAny,
				},
				FSGroup: policyv1beta1.FSGroupStrategyOptions{
					Rule: policyv1beta1.FSGroupStrategyRunAsAny,
				},
				ReadOnlyRootFilesystem:   true,
				AllowPrivilegeEscalation: util.BoolPointer(false),
				AllowedHostPaths:         allowedHostPaths,
			},
		}, reconciler.StatePresent, nil
	}
	return &policyv1beta1.PodSecurityPolicy{
		ObjectMeta: n.NodeAgentObjectMeta(fluentbitPodSecurityPolicyName),
		Spec:       policyv1beta1.PodSecurityPolicySpec{},
	}, reconciler.StateAbsent, nil
}

func (n *nodeAgentInstance) pspClusterRole() (runtime.Object, reconciler.DesiredState, error) {
	if *n.nodeAgent.FluentbitSpec.Security.RoleBasedAccessControlCreate && n.nodeAgent.FluentbitSpec.Security.PodSecurityPolicyCreate {
		return &rbacv1.ClusterRole{
			ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleName + "-psp"),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups:     []string{"policy"},
					Resources:     []string{"podsecuritypolicies"},
					ResourceNames: []string{n.logging.QualifiedName(fluentbitPodSecurityPolicyName)},
					Verbs:         []string{"use"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRole{
		ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleName + "-psp"),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (n *nodeAgentInstance) pspClusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *n.nodeAgent.FluentbitSpec.Security.RoleBasedAccessControlCreate && n.nodeAgent.FluentbitSpec.Security.PodSecurityPolicyCreate {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleBindingName + "-psp"),
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     n.logging.QualifiedName(clusterRoleName + "-psp"),
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
		ObjectMeta: n.NodeAgentObjectMetaClusterScope(clusterRoleBindingName + "-psp"),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}
