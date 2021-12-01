// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package eventtailer

import (
	"fmt"

	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensionsconfig"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Resource redeclaration of function with return type kubernetes Object
type Resource func() (runtime.Object, reconciler.DesiredState, error)

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
