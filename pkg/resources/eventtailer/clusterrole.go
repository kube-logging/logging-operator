// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package eventtailer

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ClusterRole resource for reconciler
func (e *EventTailer) ClusterRole() (runtime.Object, reconciler.DesiredState, error) {
	clusterRole := v1.ClusterRole{
		ObjectMeta: e.clusterObjectMeta(),
		Rules: []v1.PolicyRule{
			{
				Verbs:     []string{"get", "watch", "list"},
				APIGroups: []string{"", "events.k8s.io"},
				Resources: []string{"events"},
			},
		},
	}
	return &clusterRole, reconciler.StatePresent, nil
}
