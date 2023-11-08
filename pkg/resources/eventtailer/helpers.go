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
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/cisco-open/operator-tools/pkg/utils"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Name .
func (e *EventTailer) Name() string {
	return fmt.Sprintf("%v-%v", e.customResource.ObjectMeta.Name, config.EventTailer.TailerAffix)
}

func (e *EventTailer) objectMeta() v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:            e.Name(),
		Namespace:       e.customResource.Spec.ControlNamespace,
		Labels:          e.selectorLabels(),
		OwnerReferences: e.ownerReferences(),
	}
	return meta
}

func (e *EventTailer) clusterObjectMeta() v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:            e.Name(),
		Labels:          e.selectorLabels(),
		OwnerReferences: e.ownerReferences(),
	}
	return meta
}

func (e *EventTailer) ownerReferences() []v1.OwnerReference {
	ownerReferences := []v1.OwnerReference{
		{
			APIVersion: e.customResource.TypeMeta.APIVersion,
			Kind:       e.customResource.TypeMeta.Kind,
			Name:       e.customResource.ObjectMeta.Name,
			UID:        e.customResource.ObjectMeta.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	return ownerReferences
}

func (e *EventTailer) selectorLabels() map[string]string {
	base := map[string]string{
		types.NameLabel:     config.EventTailer.TailerAffix,
		types.InstanceLabel: e.Name(),
	}
	if len(e.CommonSelectorLabels) > 0 {
		for key, val := range e.CommonSelectorLabels {
			base[key] = val
		}
	}
	return base
}
