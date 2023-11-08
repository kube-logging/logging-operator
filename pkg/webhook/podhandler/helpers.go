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

package podhandler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kube-logging/logging-operator/pkg/resources/kubetool"
	"github.com/kube-logging/logging-operator/pkg/resources/volumepath"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
	corev1 "k8s.io/api/core/v1"
)

// Command returns the desired command for the current filetailer
func (p *PodHandler) Command(filePath string) []string {
	command := []string{
		"/fluent-bit/bin/fluent-bit", "-i", "tail",
		"-p", fmt.Sprintf("path=%s", filePath),
		"-o", "file",
		"-p", "format=template",
		"-p", "template={log}",
	}
	command = append(command, config.TailerWebhook.VersionedFluentBitPathArgs("/dev/stdout")...)

	return command
}

// Container returns the assembled container for the current tailer
func (p *PodHandler) Container(name string, volumeMount corev1.VolumeMount, command []string) corev1.Container {
	container := corev1.Container{
		Name:            name,
		Image:           config.TailerWebhook.FluentBitImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command:         command,
		VolumeMounts: []corev1.VolumeMount{
			volumeMount,
		},
	}

	return container
}

// Containers assembles the required containers
func (p *PodHandler) Containers(filePaths []string, volumePaths []string, containerName string) []corev1.Container {
	containers := []corev1.Container{}

	vpl := volumepath.Init(volumePaths)

	for _, filePath := range filePaths {
		filePath := filePath
		path := vpl.Apply(volumepath.ApplyFn(
			func(strs []string, idx int) *string {
				if strings.HasPrefix(filepath.Dir(filePath), strs[idx]) {
					return &strs[idx]
				}
				return nil
			},
		)).First()

		if path == nil {
			return containers
		}

		volumeMount := kubetool.NewVolumeMountBuilder().
			WithMountPath(*path).
			WithName(p.ContainerizedVolumeName(containerName, *path)).
			VolumeMount
		command := p.Command(filePath)
		containerName := p.ContainerizedVolumeName(containerName, filePath)
		containers = append(containers, p.Container(containerName, volumeMount, command))
	}
	return containers
}

// ContainerNames .
func (p *PodHandler) ContainerNames(containers []corev1.Container) []string {
	result := []string{}

	for _, container := range containers {
		result = append(result, container.Name)
	}

	return result
}

// ContainerizedVolumeName .
func (p *PodHandler) ContainerizedVolumeName(containerName string, volumePath string) string {
	containerizedVolumePath := fmt.Sprintf("/%s%s", containerName, volumePath)
	return volumepath.ConvertFilePath(containerizedVolumePath)
}

// Volumes makes corev1.Volumes from the given paths
func (p *PodHandler) Volumes(volumePaths []string) []corev1.Volume {
	emptyDir := corev1.EmptyDirVolumeSource{
		Medium:    corev1.StorageMediumDefault,
		SizeLimit: nil, // default: unlimited
	}
	volumes := []corev1.Volume{}
	for _, v := range volumePaths {
		volumes = append(volumes,
			kubetool.NewVolumeBuilder().
				WithEmptyDir(emptyDir).
				WithName(volumepath.ConvertFilePath(v)).
				Volume)
	}

	return volumes
}
