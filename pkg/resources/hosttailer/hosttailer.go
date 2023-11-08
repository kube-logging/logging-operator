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
	"path/filepath"
	"strings"

	"emperror.dev/errors"
	"github.com/kube-logging/logging-operator/pkg/resources/kubetool"
	"github.com/kube-logging/logging-operator/pkg/resources/volumepath"
	"github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/tailer"
	v1alpha1 "github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/v1alpha1"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// HostTailer .
type HostTailer struct {
	log logr.Logger
	*reconciler.GenericResourceReconciler
	customResource       v1alpha1.HostTailer
	CommonSelectorLabels map[string]string `json:"selectorLabels,omitempty"`
}

// SetupWithBuilder .
func SetupWithBuilder(builder *builder.Builder) {
	builder.Owns(&appsv1.DaemonSet{})
}

// New constructor
func New(client client.Client, log logr.Logger, opts reconciler.ReconcilerOpts, customResource v1alpha1.HostTailer) *HostTailer {
	return &HostTailer{
		log:                       log,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, log, opts),
		customResource:            customResource,
	}
}

func (h *HostTailer) Reconcile(object runtime.Object) (*reconcile.Result, error) {
	o, state, err := h.Run()
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create desired object")
	}
	if o == nil {
		return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", h)
	}
	result, err := h.ReconcileResource(o, state)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to reconcile resource")
	}
	if result != nil {
		return result, nil
	}

	return nil, nil
}

// RegisterWatches completes the implementation of ComponentReconciler
func (h *HostTailer) RegisterWatches(*builder.Builder) {
	// placeholder
}

// Name returns generated name string
func (h *HostTailer) Name(suffix string) string {
	strs := []string{h.customResource.ObjectMeta.Name, config.HostTailer.TailerAffix}
	if suffix != "" {
		strs = append(strs, suffix)
	}
	return strings.Join(strs, "-")
}

// Run is the implementation of type Resource interface, generates Deamonset for fileTailers
func (h *HostTailer) Run() (runtime.Object, reconciler.DesiredState, error) {
	tailers := h.MergeTailers()
	volumePaths := h.GatherPathStrings(tailers)

	pathList := volumepath.Init(volumePaths).TopLevelPathList().Uniq()

	volumes := h.Volumes(pathList.Strings())
	containers := h.Containers(tailers, pathList.Strings())

	if len(containers) > 0 {
		ds := h.DaemonSet(containers, volumes)
		return ds, reconciler.StatePresent, nil
	}
	return &appsv1.DaemonSet{ObjectMeta: h.objectMeta()}, reconciler.StateAbsent, nil
}

// MergeTailers collects all the tailers through the different tailer types and cast them to Tailer interface array
func (h *HostTailer) MergeTailers() []tailer.Tailer {
	tailers := make([]tailer.Tailer, 0)
	for _, fileTailer := range h.customResource.Spec.FileTailers {
		tailers = append(tailers, fileTailer)
	}
	for _, fileTailer := range h.customResource.Spec.SystemdTailers {
		tailers = append(tailers, fileTailer)
	}
	return tailers
}

// Containers assembles the required containers
func (h *HostTailer) Containers(tailers []tailer.Tailer, volumePaths []string) []corev1.Container {
	containers := []corev1.Container{}

	vpl := volumepath.Init(volumePaths)

	for _, t := range tailers {
		generalDescriptor := t.GeneralDescriptor()
		if generalDescriptor.Disabled {
			continue
		}

		path := vpl.Apply(volumepath.ApplyFn(
			func(strs []string, idx int) *string {
				if strings.HasPrefix(filepath.Dir(generalDescriptor.Path), strs[idx]) {
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
			WithName(volumepath.ConvertFilePath(*path)).
			VolumeMount
		command := t.Command(h.Name(generalDescriptor.Name))
		containers = append(containers, h.Container(generalDescriptor.Name, volumeMount, command, generalDescriptor.ContainerBase, generalDescriptor.Image))
	}
	return containers
}

// DaemonSet produces k8s daemonset struct from the given parameters
func (h *HostTailer) DaemonSet(containers []corev1.Container, volumes []corev1.Volume) *appsv1.DaemonSet {
	ds := &appsv1.DaemonSet{
		ObjectMeta: h.objectMeta(),
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: h.selectorLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: h.customResource.Spec.WorkloadMetaBase.Merge(v1.ObjectMeta{
					Labels: h.selectorLabels(),
				}),
				Spec: h.customResource.Spec.WorkloadBase.Override(corev1.PodSpec{
					Volumes:    volumes,
					Containers: containers,
				}),
			},
		},
	}
	return ds
}

// GatherPathStrings returns the path strings
func (h *HostTailer) GatherPathStrings(tailers []tailer.Tailer) []string {
	result := []string{}
	for _, t := range tailers {
		generalDescriptor := t.GeneralDescriptor()
		if generalDescriptor.Disabled {
			continue
		}
		result = append(result, filepath.Dir(generalDescriptor.Path))
	}
	return result
}
