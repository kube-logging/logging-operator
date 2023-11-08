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
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) clusterPodSecurityPolicy() (runtime.Object, reconciler.DesiredState, error) {
	if r.fluentbitSpec.Security.PodSecurityPolicyCreate {
		allowedHostPaths := []policyv1beta1.AllowedHostPath{{
			PathPrefix: r.fluentbitSpec.MountPath,
			ReadOnly:   true,
		}, {
			PathPrefix: "/var/log",
			ReadOnly:   true,
		}}

		if r.fluentbitSpec.BufferStorageVolume.HostPath != nil {
			allowedHostPaths = append(allowedHostPaths, policyv1beta1.AllowedHostPath{
				PathPrefix: r.fluentbitSpec.BufferStorageVolume.HostPath.Path,
				ReadOnly:   false,
			})
		}

		for _, vMnt := range r.fluentbitSpec.ExtraVolumeMounts {
			allowedHostPaths = append(allowedHostPaths, policyv1beta1.AllowedHostPath{
				PathPrefix: vMnt.Source,
				ReadOnly:   *vMnt.ReadOnly,
			})
		}

		if r.fluentbitSpec.PositionDB.HostPath != nil {
			r.fluentbitSpec.PositionDB.WithDefaultHostPath(
				fmt.Sprintf(v1beta1.HostPath, r.nameProvider.Name(), TailPositionVolume))

			allowedHostPaths = append(allowedHostPaths, policyv1beta1.AllowedHostPath{
				PathPrefix: r.fluentbitSpec.PositionDB.HostPath.Path,
				ReadOnly:   false,
			})
		}

		hostPorts := []policyv1beta1.HostPortRange{}
		if r.fluentbitSpec.HostNetwork {
			if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.Metrics.Port != 0 {
				hostPorts = append(hostPorts, policyv1beta1.HostPortRange{
					Min: r.fluentbitSpec.Metrics.Port,
					Max: r.fluentbitSpec.Metrics.Port,
				})
			}

			if r.fluentbitSpec.BufferVolumeMetrics != nil && r.fluentbitSpec.BufferVolumeMetrics.Port != 0 {
				hostPorts = append(hostPorts, policyv1beta1.HostPortRange{
					Min: r.fluentbitSpec.BufferVolumeMetrics.Port,
					Max: r.fluentbitSpec.BufferVolumeMetrics.Port,
				})
			}
		}

		return &policyv1beta1.PodSecurityPolicy{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(fluentbitPodSecurityPolicyName),
			Spec: policyv1beta1.PodSecurityPolicySpec{
				Volumes: []policyv1beta1.FSType{
					"configMap",
					"emptyDir",
					"secret",
					"hostPath",
					"persistentVolumeClaim"},
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
				HostNetwork:              r.fluentbitSpec.HostNetwork,
				HostPorts:                hostPorts,
			},
		}, reconciler.StatePresent, nil
	}
	return &policyv1beta1.PodSecurityPolicy{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitPodSecurityPolicyName),
		Spec:       policyv1beta1.PodSecurityPolicySpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) pspClusterRole() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentbitSpec.Security.RoleBasedAccessControlCreate && r.fluentbitSpec.Security.PodSecurityPolicyCreate {
		return &rbacv1.ClusterRole{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleName + "-psp"),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups:     []string{"policy"},
					Resources:     []string{"podsecuritypolicies"},
					ResourceNames: []string{r.nameProvider.ComponentName(fluentbitPodSecurityPolicyName)},
					Verbs:         []string{"use"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRole{
		ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleName + "-psp"),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) pspClusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentbitSpec.Security.RoleBasedAccessControlCreate && r.fluentbitSpec.Security.PodSecurityPolicyCreate {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleBindingName + "-psp"),
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.nameProvider.ComponentName(clusterRoleName + "-psp"),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      r.getServiceAccount(),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleBindingName + "-psp"),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}
