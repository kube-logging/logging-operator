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
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) role() (runtime.Object, k8sutil.DesiredState, error) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.Role{
			ObjectMeta: templates.FluentdObjectMeta(r.Logging.QualifiedName(roleName), r.Logging.Labels, r.Logging),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"configmaps", "secrets"},
					Verbs:     []string{"*"},
				},
			},
		}, k8sutil.StatePresent, nil
	}
	return &rbacv1.Role{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(roleName),
			r.Logging.Labels, r.Logging,
		),
		Rules: []rbacv1.PolicyRule{}}, k8sutil.StateAbsent, nil
}

func (r *Reconciler) roleBinding() (runtime.Object, k8sutil.DesiredState, error) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.RoleBinding{
			ObjectMeta: templates.FluentdObjectMeta(r.Logging.QualifiedName(roleBindingName), r.Logging.Labels, r.Logging),
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.Logging.QualifiedName(roleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      r.Logging.QualifiedName(defaultServiceAccountName),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, k8sutil.StatePresent, nil
	}
	return &rbacv1.RoleBinding{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(roleBindingName),
			r.Logging.Labels, r.Logging,
		),
		RoleRef: rbacv1.RoleRef{}}, k8sutil.StateAbsent, nil
}

func (r *Reconciler) serviceAccount() (runtime.Object, k8sutil.DesiredState, error) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.FluentdSpec.Security.ServiceAccount == "" {
		return &corev1.ServiceAccount{
			ObjectMeta: templates.FluentdObjectMeta(
				r.Logging.QualifiedName(defaultServiceAccountName),
				r.Logging.Labels,
				r.Logging,
			),
		}, k8sutil.StatePresent, nil
	}
	return &corev1.ServiceAccount{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(defaultServiceAccountName),
			r.Logging.Labels,
			r.Logging,
		),
	}, k8sutil.StateAbsent, nil
}
