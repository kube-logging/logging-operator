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

package eventtailer

import (
	"github.com/cisco-open/operator-tools/pkg/reconciler"
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
