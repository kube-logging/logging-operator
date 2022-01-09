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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/banzaicloud/logging-operator/pkg/resources/annotation"
	"github.com/banzaicloud/logging-operator/pkg/resources/kubetool"
	"github.com/banzaicloud/logging-operator/pkg/resources/volumepath"
	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/extensionsconfig"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// PodHandler .
type PodHandler struct {
	Client  client.Client
	decoder *admission.Decoder
}

var _ admission.Handler = &PodHandler{}

// NewPodHandler constructor
func NewPodHandler(client client.Client) *PodHandler {
	return &PodHandler{Client: client}
}

func (p *PodHandler) sideCarsForContainer(containerName string, filesToTail []string) (sideCars []corev1.Container, volumes []corev1.Volume, volumeMounts []corev1.VolumeMount) {
	fileList := volumepath.Init(filesToTail).Uniq().RemoveInvalidPath(nil)

	// get list of dirs from fileList
	dirList := fileList.Apply(volumepath.ApplyFn(
		func(paths []string, idx int) *string {
			fileDir := filepath.Dir(paths[idx])
			return &fileDir
		},
	)).TopLevelPathList().Uniq()

	// generate containerized path list /{containername}/path
	volumePathList := dirList.Apply(volumepath.ApplyFn(
		func(paths []string, idx int) *string {
			fileDir := fmt.Sprintf("/%s%s", containerName, paths[idx])
			return &fileDir
		},
	))

	// create volumes needed by the container
	volumes = p.Volumes(volumePathList.Strings())

	// create sidecar containers
	sideCars = p.Containers(fileList.Strings(), dirList.Strings(), containerName)

	// generate volumemounts to mount them to then container
	for _, dir := range *dirList {
		volumeMounts = append(
			volumeMounts, kubetool.NewVolumeMountBuilder().
				WithMountPath(dir).
				WithName(p.ContainerizedVolumeName(containerName, dir)).
				VolumeMount)
	}

	return
}

// Handle .
func (p *PodHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := p.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// check annotations
	tailAnnotation, ok := pod.Annotations[config.TailerWebhook.AnnotationKey]
	if !ok {
		return admission.Allowed("changing podspec is not required")
	}

	// collect existing containers to process annotations properly
	containerNames := p.ContainerNames(pod.Spec.Containers)

	// handle the tail annotation string of the pod
	annotationHandler := annotation.NewHandler(containerNames)
	annotationHandler.AddTailerAnnotation(tailAnnotation)

	for idx, container := range pod.Spec.Containers {
		filePaths := annotationHandler.FilePathsForContainer(container.Name)

		sideCars, volumes, volumeMounts := p.sideCarsForContainer(container.Name, filePaths)

		// Append the new data to the podspec
		pod.Spec.Containers = append(pod.Spec.Containers, sideCars...)
		pod.Spec.Volumes = append(pod.Spec.Volumes, volumes...)
		pod.Spec.Containers[idx].VolumeMounts = append(pod.Spec.Containers[idx].VolumeMounts, volumeMounts...)
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// InjectDecoder injects the decoder.
func (p *PodHandler) InjectDecoder(d *admission.Decoder) error {
	p.decoder = d
	return nil
}
