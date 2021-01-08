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

package nodeagent

import (
	"fmt"

	util "github.com/banzaicloud/operator-tools/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeAgentObjectMeta creates an objectMeta for resource fluentbit
func (n *nodeAgentInstance) NodeAgentObjectMeta(name string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", n.logging.QualifiedName(name), n.nodeAgent.Name),
		Namespace: n.logging.Spec.ControlNamespace,
		Labels:    n.getFluentBitLabels(),
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: n.logging.APIVersion,
				Kind:       n.logging.Kind,
				Name:       n.logging.Name,
				UID:        n.logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}

// NodeAgentObjectMetaClusterScope creates an cluster scoped objectMeta for resource fluentbit
func (n *nodeAgentInstance) NodeAgentObjectMetaClusterScope(name string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   fmt.Sprintf("%s-%s", n.logging.QualifiedName(name), n.nodeAgent.Name),
		Labels: n.getFluentBitLabels(),
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: n.logging.APIVersion,
				Kind:       n.logging.Kind,
				Name:       n.logging.Name,
				UID:        n.logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}
