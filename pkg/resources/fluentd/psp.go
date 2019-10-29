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

package fluentd

import (
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) clusterPodSecurityPolicy() (runtime.Object, k8sutil.DesiredState) {
	if r.Logging.Spec.FluentdSpec.Security.PodSecurityPolicyCreate {

		return &policyv1beta1.PodSecurityPolicy{
			ObjectMeta: templates.FluentdObjectMetaClusterScope(
				r.Logging.QualifiedName(PodSecurityPolicyName),
				util.MergeLabels(r.Logging.Labels, r.getFluentdLabels()), r.Logging),
			Spec: policyv1beta1.PodSecurityPolicySpec{
				Volumes: []policyv1beta1.FSType{
					"configMap",
					"emptyDir",
					"secret",
					"persistentVolumeClaim"},
				SELinux: policyv1beta1.SELinuxStrategyOptions{
					Rule: policyv1beta1.SELinuxStrategyRunAsAny,
				},
				RunAsUser: policyv1beta1.RunAsUserStrategyOptions{
					Rule:   policyv1beta1.RunAsUserStrategyMustRunAs,
					Ranges: []policyv1beta1.IDRange{{Min: 1, Max: 65535}}},
				SupplementalGroups: policyv1beta1.SupplementalGroupsStrategyOptions{
					Rule:   policyv1beta1.SupplementalGroupsStrategyMustRunAs,
					Ranges: []policyv1beta1.IDRange{{Min: 1, Max: 65535}},
				},
				FSGroup: policyv1beta1.FSGroupStrategyOptions{
					Rule:   policyv1beta1.FSGroupStrategyMustRunAs,
					Ranges: []policyv1beta1.IDRange{{Min: 1, Max: 65535}},
				},
				AllowPrivilegeEscalation: util.BoolPointer(false),
			},
		}, k8sutil.StatePresent

	}
	return &policyv1beta1.PodSecurityPolicy{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(PodSecurityPolicyName),
			util.MergeLabels(r.Logging.Labels, r.getFluentdLabels()), r.Logging),
		Spec: policyv1beta1.PodSecurityPolicySpec{},
	}, k8sutil.StateAbsent
}

func (r *Reconciler) pspRole() (runtime.Object, k8sutil.DesiredState) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.FluentdSpec.Security.PodSecurityPolicyCreate {

		return &rbacv1.Role{
			ObjectMeta: templates.FluentdObjectMeta(r.Logging.QualifiedName(roleName+"-psp"), r.Logging.Labels, r.Logging),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"policy"},
					Resources: []string{r.Logging.QualifiedName(PodSecurityPolicyName)},
					Verbs:     []string{"use"},
				},
			},
		}, k8sutil.StatePresent
	}
	return &rbacv1.Role{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(roleName+"-psp"),
			r.Logging.Labels, r.Logging,
		),
		Rules: []rbacv1.PolicyRule{}}, k8sutil.StateAbsent
}

func (r *Reconciler) pspRoleBinding() (runtime.Object, k8sutil.DesiredState) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.FluentdSpec.Security.PodSecurityPolicyCreate {

		return &rbacv1.RoleBinding{
			ObjectMeta: templates.FluentdObjectMeta(r.Logging.QualifiedName(roleBindingName+"-psp"), r.Logging.Labels, r.Logging),
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.Logging.QualifiedName(roleName + "-psp"),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      r.Logging.QualifiedName(defaultServiceAccountName),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, k8sutil.StatePresent
	}
	return &rbacv1.RoleBinding{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(roleBindingName+"-psp"),
			r.Logging.Labels, r.Logging,
		),
		RoleRef: rbacv1.RoleRef{}}, k8sutil.StateAbsent
}
