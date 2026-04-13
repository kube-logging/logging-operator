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

	"github.com/go-logr/logr"
	"github.com/siliconbrain/go-seqs/seqs"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/kube-logging/logging-operator/pkg/resources/annotation"
	"github.com/kube-logging/logging-operator/pkg/resources/kubetool"
	"github.com/kube-logging/logging-operator/pkg/resources/volumepath"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
)

// PodHandler .
type PodHandler struct {
	Decoder admission.Decoder
	Log     logr.Logger
}

var _ admission.Handler = &PodHandler{}

// NewPodHandler constructor
func NewPodHandler(log logr.Logger) *PodHandler {
	return &PodHandler{Log: log}
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
	log := p.Log.WithValues("namespace", req.Namespace, "name", req.Name)

	log.Info("webhook handler called")

	pod := &corev1.Pod{}

	err := p.Decoder.Decode(req, pod)
	if err != nil {
		log.Error(err, "unable to decode pod")
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

	// Build a snapshot of original container names and their indices before
	// the mutation loop, so the indices recorded here remain stable.
	type containerRef struct {
		name string
		idx  int
	}
	originalContainers := make([]containerRef, len(pod.Spec.Containers))
	for i, c := range pod.Spec.Containers {
		originalContainers[i] = containerRef{name: c.Name, idx: i}
	}

	// Iterate over the snapshot — not over pod.Spec.Containers directly —
	// so that appended sidecars are never visited.
	for _, ref := range originalContainers {
		filePaths := annotationHandler.FilePathsForContainer(ref.name)
		if len(filePaths) == 0 {
			continue
		}

		sideCars, volumes, volumeMounts := p.sideCarsForContainer(ref.name, filePaths)

		if resp := p.podHandlerHelper(pod, ref.idx, sideCars, volumes, volumeMounts); resp != nil {
			return *resp
		}
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		log.Error(err, "pod marshaling failed")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (p *PodHandler) podHandlerHelper(podToModify *corev1.Pod, targetContainerIdx int, sideCars []corev1.Container, volumes []corev1.Volume, volumeMounts []corev1.VolumeMount) *admission.Response {
	duplicateValuesMsg := "webhook mutation would result in duplicate values, returning"

	// Bounds check: targetContainerIdx must point to a valid original container.
	if targetContainerIdx < 0 || targetContainerIdx >= len(podToModify.Spec.Containers) {
		msg := fmt.Sprintf("targetContainerIdx %d is out of range (containers: %d)", targetContainerIdx, len(podToModify.Spec.Containers))
		p.Log.Info(msg)
		rv := admission.Denied(msg)
		return &rv
	}

	for _, sideCar := range sideCars {
		if seqs.Any(seqs.FromSlice(podToModify.Spec.Containers), func(c corev1.Container) bool {
			return c.Name == sideCar.Name
		}) {
			p.Log.Info(duplicateValuesMsg)
			rv := admission.Denied(duplicateValuesMsg)
			return &rv
		}
		podToModify.Spec.Containers = append(podToModify.Spec.Containers, sideCar)
	}

	// Add shared volumeMounts only to the target container (the one whose
	// files are being tailed), not to all original containers. This prevents
	// masking filesystem paths in unrelated containers.
	for _, vm := range volumeMounts {
		existing, found := findVolumeMount(podToModify.Spec.Containers[targetContainerIdx].VolumeMounts, vm.MountPath)
		if found {
			if existing.Name != vm.Name {
				// Same mountPath but different volume name — the sidecar would
				// mount the webhook's emptyDir while the app keeps its own volume,
				// so the tailer would never see the log files. Deny the mutation
				// to surface the misconfiguration.
				msg := fmt.Sprintf(
					"container %q already has mountPath %q with volume %q, but webhook needs volume %q; "+
						"rename the existing volume to match or remove the conflicting mount",
					podToModify.Spec.Containers[targetContainerIdx].Name, vm.MountPath, existing.Name, vm.Name,
				)
				p.Log.Info(msg)
				rv := admission.Denied(msg)
				return &rv
			}
			// Same mountPath and same volume name — already correctly mounted, skip.
			continue
		}
		podToModify.Spec.Containers[targetContainerIdx].VolumeMounts = append(
			podToModify.Spec.Containers[targetContainerIdx].VolumeMounts, vm)
	}

	// Append volumes. In multi-container pods, two containers tailing files
	// in the same directory produce the same volume name; the second call
	// must skip rather than deny, because the volume was already added by
	// the first call.
	for _, volume := range volumes {
		if existing, found := findVolume(podToModify.Spec.Volumes, volume.Name); found {
			// A volume with this name already exists. Only skip if it is
			// compatible (also an EmptyDir). If the existing volume has a
			// different source type, the sidecar would mount an unexpected
			// volume and tailing would silently fail.
			if existing.EmptyDir == nil {
				msg := fmt.Sprintf(
					"volume %q already exists with a non-EmptyDir source; "+
						"the tailer webhook requires an EmptyDir volume at this name",
					volume.Name,
				)
				p.Log.Info(msg)
				rv := admission.Denied(msg)
				return &rv
			}
			p.Log.V(1).Info("compatible volume already exists, skipping", "volume", volume.Name)
			continue
		}
		podToModify.Spec.Volumes = append(podToModify.Spec.Volumes, volume)
	}

	return nil
}

// findVolumeMount checks if the given slice contains a mount at mountPath.
// If found, it returns the existing mount and true; otherwise zero value and false.
func findVolumeMount(mounts []corev1.VolumeMount, mountPath string) (corev1.VolumeMount, bool) {
	for _, m := range mounts {
		if m.MountPath == mountPath {
			return m, true
		}
	}
	return corev1.VolumeMount{}, false
}

// findVolume checks if the given slice contains a volume with the given name.
// If found, it returns the existing volume and true; otherwise zero value and false.
func findVolume(volumes []corev1.Volume, name string) (corev1.Volume, bool) {
	for _, v := range volumes {
		if v.Name == name {
			return v, true
		}
	}
	return corev1.Volume{}, false
}
