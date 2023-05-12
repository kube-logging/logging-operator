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
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) clusterRole() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentbitSpec.Security.RoleBasedAccessControlCreate {
		clusterRoleResources := []string{"pods", "namespaces"}
		if r.fluentbitSpec.FilterKubernetes.UseKubelet == "On" {
			clusterRoleResources = append(clusterRoleResources, "nodes", "nodes/proxy")
		}
		return &rbacv1.ClusterRole{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleName),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: clusterRoleResources,
					Verbs:     []string{"get", "list", "watch"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.ClusterRole{
		ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleName),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) clusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentbitSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleBindingName),
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.nameProvider.ComponentName(clusterRoleName),
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
		ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleBindingName),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceAccount() (runtime.Object, reconciler.DesiredState, error) {
	if *r.fluentbitSpec.Security.RoleBasedAccessControlCreate && r.fluentbitSpec.Security.ServiceAccount == "" {
		desired := &corev1.ServiceAccount{
			ObjectMeta: r.FluentbitObjectMeta(defaultServiceAccountName),
		}
		err := merge.Merge(desired, r.fluentbitSpec.ServiceAccountOverrides)
		if err != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
		}

		return desired, reconciler.StatePresent, nil
	} else {
		desired := &corev1.ServiceAccount{
			ObjectMeta: r.FluentbitObjectMeta(defaultServiceAccountName),
		}
		return desired, reconciler.StateAbsent, nil
	}
}
