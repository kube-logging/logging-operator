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

package hosttailer

import (
	"github.com/banzaicloud/logging-operator/pkg/resources/kubetool"
	"github.com/banzaicloud/logging-operator/pkg/resources/volumepath"
	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/extensionsconfig"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *HostTailer) ownerReferences() []v1.OwnerReference {
	ownerReferences := []v1.OwnerReference{
		{
			APIVersion: h.customResource.TypeMeta.APIVersion,
			Kind:       h.customResource.TypeMeta.Kind,
			Name:       h.customResource.ObjectMeta.Name,
			UID:        h.customResource.ObjectMeta.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	return ownerReferences
}

func (h *HostTailer) selectorLabels() map[string]string {
	base := map[string]string{
		types.NameLabel:     config.HostTailer.TailerAffix,
		types.InstanceLabel: h.Name(""),
	}
	if len(h.CommonSelectorLabels) > 0 {
		for key, val := range h.CommonSelectorLabels {
			base[key] = val
		}
	}
	return base
}

func (h *HostTailer) objectMeta() v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:            h.Name(""),
		Namespace:       h.customResource.ObjectMeta.Namespace,
		Labels:          h.selectorLabels(),
		OwnerReferences: h.ownerReferences(),
	}
	return meta
}

// Container returns the assembled container for the current tailer
func (h *HostTailer) Container(name string, volumeMount corev1.VolumeMount, command []string, overrides *types.ContainerBase) corev1.Container {
	container := corev1.Container{
		Name:            name,
		Image:           config.HostTailer.FluentBitImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command:         command,
		VolumeMounts: []corev1.VolumeMount{
			kubetool.NewVolumeMountBuilder().
				WithMountPath(config.Global.FluentBitPosFilePath).
				WithName(config.Global.FluentBitPosVolumeName).
				VolumeMount,
			volumeMount,
		},
	}
	if overrides != nil {
		container = overrides.Override(container)
	}

	return container
}

// Volumes makes corev1.Volumes from the given paths
func (h *HostTailer) Volumes(volumePaths []string) []corev1.Volume {
	volumes := []corev1.Volume{
		kubetool.NewVolumeBuilder().
			WithHostPathFromPath(config.Global.FluentBitPosFilePath).
			WithName(config.Global.FluentBitPosVolumeName).
			Volume,
	}

	for _, v := range volumePaths {
		volumes = append(volumes,
			kubetool.NewVolumeBuilder().
				WithHostPathFromPath(v).
				WithName(volumepath.ConvertFilePath(v)).
				Volume)
	}

	return volumes
}
