// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package eventtailer

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ClusterRoleBinding resource for reconciler
func (e *EventTailer) ClusterRoleBinding() (runtime.Object, reconciler.DesiredState, error) {
	clusterRoleBinding := v1.ClusterRoleBinding{
		ObjectMeta: e.clusterObjectMeta(),
		Subjects: []v1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      e.Name(),
				Namespace: e.customResource.Spec.ControlNamespace,
			},
		},
		RoleRef: v1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     e.Name(),
		},
	}
	return &clusterRoleBinding, reconciler.StatePresent, nil
}
