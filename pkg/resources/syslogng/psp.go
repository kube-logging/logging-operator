// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) clusterPodSecurityPolicy() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.SyslogNGSpec.Security.PodSecurityPolicyCreate {
		return &policyv1beta1.PodSecurityPolicy{
			ObjectMeta: r.SyslogNGObjectMetaClusterScope(PodSecurityPolicyName, ComponentSyslogNG),
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
					Rule:   policyv1beta1.RunAsUserStrategyMustRunAs,
					Ranges: []policyv1beta1.IDRange{{Min: 100, Max: 100}}},
				SupplementalGroups: policyv1beta1.SupplementalGroupsStrategyOptions{
					Rule:   policyv1beta1.SupplementalGroupsStrategyMustRunAs,
					Ranges: []policyv1beta1.IDRange{{Min: 101, Max: 101}},
				},
				FSGroup: policyv1beta1.FSGroupStrategyOptions{
					Rule:   policyv1beta1.FSGroupStrategyMustRunAs,
					Ranges: []policyv1beta1.IDRange{{Min: 101, Max: 101}},
				},
				AllowPrivilegeEscalation: util.BoolPointer(false),
			},
		}, reconciler.StatePresent, nil
	}
	return &policyv1beta1.PodSecurityPolicy{
		ObjectMeta: r.SyslogNGObjectMeta(PodSecurityPolicyName, ComponentSyslogNG),
		Spec:       policyv1beta1.PodSecurityPolicySpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) pspRole() (runtime.Object, reconciler.DesiredState, error) {
	if *r.Logging.Spec.SyslogNGSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.SyslogNGSpec.Security.PodSecurityPolicyCreate {
		return &rbacv1.Role{
			ObjectMeta: r.SyslogNGObjectMeta(roleName+"-psp", ComponentSyslogNG),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups:     []string{"policy"},
					Resources:     []string{"podsecuritypolicies"},
					ResourceNames: []string{r.Logging.QualifiedName(PodSecurityPolicyName)},
					Verbs:         []string{"use"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.Role{
		ObjectMeta: r.SyslogNGObjectMeta(roleName+"-psp", ComponentSyslogNG),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) pspRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.Logging.Spec.SyslogNGSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.SyslogNGSpec.Security.PodSecurityPolicyCreate {
		return &rbacv1.RoleBinding{
			ObjectMeta: r.SyslogNGObjectMeta(roleBindingName+"-psp", ComponentSyslogNG),
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.Logging.QualifiedName(roleName + "-psp"),
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
	return &rbacv1.RoleBinding{
		ObjectMeta: r.SyslogNGObjectMeta(roleBindingName+"-psp", ComponentSyslogNG),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}
