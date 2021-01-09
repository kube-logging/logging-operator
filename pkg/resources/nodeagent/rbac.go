// Copyright © 2021 Banzai Cloud
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

package nodeagent

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (n *nodeAgentInstance) serviceAccount() (runtime.Object, reconciler.DesiredState, error) {
	if *n.nodeAgent.FluentbitSpec.Security.RoleBasedAccessControlCreate && n.nodeAgent.FluentbitSpec.Security.ServiceAccount == "" {
		return &corev1.ServiceAccount{
			ObjectMeta: n.nodeAgent.FluentbitSpec.MetaOverride.Merge(n.NodeAgentObjectMeta(defaultServiceAccountName)),
		}, reconciler.StatePresent, nil
	}
	return &corev1.ServiceAccount{
		ObjectMeta: n.NodeAgentObjectMeta(defaultServiceAccountName),
	}, reconciler.StateAbsent, nil
}
