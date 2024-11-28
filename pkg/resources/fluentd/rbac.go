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
	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) role() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.Role{
			ObjectMeta: r.FluentdObjectMeta(roleName, ComponentFluentd),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"configmaps", "secrets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.Role{
		ObjectMeta: r.FluentdObjectMeta(roleName, ComponentFluentd),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) roleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.RoleBinding{
			ObjectMeta: r.FluentdObjectMeta(roleBindingName, ComponentFluentd),
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				APIGroup: rbacv1.GroupName,
				Name:     r.Logging.QualifiedName(roleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      r.getServiceAccount(),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.RoleBinding{
		ObjectMeta: r.FluentdObjectMeta(roleBindingName, ComponentFluentd),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) sccRole() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.CreateOpenShiftSCC {
		return &rbacv1.Role{
			ObjectMeta: r.FluentdObjectMeta(sccRoleName, ComponentFluentd),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups:     []string{"security.openshift.io"},
					ResourceNames: []string{"anyuid"},
					Resources:     []string{"securitycontextconstraints"},
					Verbs:         []string{"use"},
				},
			},
		}, reconciler.StatePresent, nil
	}

	return &rbacv1.Role{
		ObjectMeta: r.FluentdObjectMeta(sccRoleName, ComponentFluentd),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) sccRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.CreateOpenShiftSCC {
		return &rbacv1.RoleBinding{
			ObjectMeta: r.FluentdObjectMeta(sccRoleName, ComponentFluentd),
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				APIGroup: rbacv1.GroupName,
				Name:     r.Logging.QualifiedName(sccRoleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      r.getServiceAccount(),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.RoleBinding{
		ObjectMeta: r.FluentdObjectMeta(sccRoleName, ComponentFluentd),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) clusterRole() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRole{
			ObjectMeta: r.FluentdObjectMetaClusterScope(clusterRoleName, ComponentFluentd),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{
						"configmaps",
						"events",
						"nodes",
						"endpoints",
						"services",
						"pods",
					},
					Verbs: []string{"get", "list", "watch"},
				}, {
					APIGroups: []string{"apps"},
					Resources: []string{
						"daemonsets",
						"deployments",
						"replicasets",
						"statefulsets",
					},
					Verbs: []string{"get", "list", "watch"},
				},
				{
					APIGroups: []string{"events.k8s.io"},
					Resources: []string{
						"events",
					},
					Verbs: []string{"get", "list", "watch"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRole{
		ObjectMeta: r.FluentdObjectMetaClusterScope(clusterRoleName, ComponentFluentd),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) clusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: r.FluentdObjectMetaClusterScope(clusterRoleBindingName, ComponentFluentd),
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				APIGroup: rbacv1.GroupName,
				Name:     r.Logging.QualifiedName(roleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      r.getServiceAccount(),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: r.FluentdObjectMetaClusterScope(clusterRoleBindingName, ComponentFluentd),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceAccount() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentdSpec.Security.RoleBasedAccessControlCreate && r.fluentdSpec.Security.ServiceAccount == "" {
		desired := &corev1.ServiceAccount{
			ObjectMeta: r.FluentdObjectMeta(defaultServiceAccountName, ComponentFluentd),
		}
		err := merge.Merge(desired, r.fluentdSpec.ServiceAccountOverrides)
		if err != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
		}

		return desired, reconciler.StatePresent, nil
	} else {
		desired := &corev1.ServiceAccount{
			ObjectMeta: r.FluentdObjectMeta(defaultServiceAccountName, ComponentFluentd),
		}
		return desired, reconciler.StateAbsent, nil
	}
}
