// Copyright © 2022 Cisco Systems, Inc. and/or its affiliates
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
	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) role() (runtime.Object, reconciler.DesiredState, error) {
	role := &rbacv1.Role{
		ObjectMeta: r.SyslogNGObjectMeta(roleName, ComponentSyslogNG),
	}
	if r.syslogNGSpec == nil || r.syslogNGSpec.SkipRBACCreate {
		return role, reconciler.StateAbsent, nil
	}
	role.Rules = []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps", "secrets"},
			Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
		},
	}
	return role, reconciler.StatePresent, nil
}

func (r *Reconciler) roleBinding() (runtime.Object, reconciler.DesiredState, error) {
	binding := &rbacv1.RoleBinding{
		ObjectMeta: r.SyslogNGObjectMeta(roleBindingName, ComponentSyslogNG),
	}
	if r.syslogNGSpec == nil || r.syslogNGSpec.SkipRBACCreate {
		return binding, reconciler.StateAbsent, nil
	}
	binding.RoleRef = rbacv1.RoleRef{
		Kind:     "Role",
		APIGroup: "rbac.authorization.k8s.io",
		Name:     r.Logging.QualifiedName(roleName),
	}
	binding.Subjects = []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      r.getServiceAccountName(),
			Namespace: r.Logging.Spec.ControlNamespace,
		},
	}
	return binding, reconciler.StatePresent, nil
}

func (r *Reconciler) isEnhanceK8sFilter() bool {
	for _, f := range r.Logging.Spec.GlobalFilters {
		if f.EnhanceK8s != nil {
			return true
		}
	}
	return false
}

func (r *Reconciler) clusterRole() (runtime.Object, reconciler.DesiredState, error) {
	role := &rbacv1.ClusterRole{
		ObjectMeta: r.SyslogNGObjectMetaClusterScope(clusterRoleName, ComponentSyslogNG),
	}
	if r.syslogNGSpec == nil || r.syslogNGSpec.SkipRBACCreate || !r.isEnhanceK8sFilter() {
		return role, reconciler.StateAbsent, nil
	}
	role.Rules = []rbacv1.PolicyRule{
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
		},
		{
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
	}
	return role, reconciler.StatePresent, nil
}

func (r *Reconciler) clusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	binding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: r.SyslogNGObjectMetaClusterScope(clusterRoleBindingName, ComponentSyslogNG),
	}
	if r.syslogNGSpec == nil || r.syslogNGSpec.SkipRBACCreate || !r.isEnhanceK8sFilter() {
		return binding, reconciler.StateAbsent, nil
	}
	binding.RoleRef = rbacv1.RoleRef{
		Kind:     "ClusterRole",
		APIGroup: "rbac.authorization.k8s.io",
		Name:     r.Logging.QualifiedName(roleName),
	}
	binding.Subjects = []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      r.getServiceAccountName(),
			Namespace: r.Logging.Spec.ControlNamespace,
		},
	}
	return binding, reconciler.StatePresent, nil
}

func (r *Reconciler) serviceAccount() (runtime.Object, reconciler.DesiredState, error) {
	account := &corev1.ServiceAccount{
		ObjectMeta: r.SyslogNGObjectMeta(serviceAccountName, ComponentSyslogNG),
	}
	err := merge.Merge(account, r.syslogNGSpec.ServiceAccountOverrides)
	err = errors.WrapIf(err, "unable to merge overrides to base object")

	if r.syslogNGSpec == nil || r.syslogNGSpec.SkipRBACCreate {
		return account, reconciler.StateAbsent, err
	}

	return account, reconciler.StatePresent, err
}
