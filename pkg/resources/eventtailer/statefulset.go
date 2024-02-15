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
	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/cisco-open/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
)

// StatefulSet resource for reconciler
func (e *EventTailer) StatefulSet() (runtime.Object, reconciler.DesiredState, error) {
	var err error
	spec := e.statefulSetSpec()

	if e.customResource.Spec.PositionVolume.PersistentVolumeClaim != nil {
		err = e.customResource.Spec.PositionVolume.ApplyPVCForStatefulSet(config.EventTailer.TailerAffix, config.Global.FluentBitPosFilePath, spec, func(name string) metav1.ObjectMeta {
			return e.objectMeta()
		})
	} else {
		err = e.customResource.Spec.PositionVolume.ApplyVolumeForPodSpec(config.EventTailer.PositionVolumeName, config.EventTailer.TailerAffix, config.Global.FluentBitPosFilePath, &spec.Template.Spec)
	}

	statefulSet := appsv1.StatefulSet{
		ObjectMeta: e.objectMeta(),
		Spec:       *spec,
	}
	return &statefulSet, reconciler.StatePresent, err
}

func (e *EventTailer) statefulSetSpec() *appsv1.StatefulSetSpec {

	if e.customResource.Spec.Image != nil {
		if repositoryWithTag := e.customResource.Spec.Image.RepositoryWithTag(); repositoryWithTag != "" {
			if e.customResource.Spec.ContainerBase == nil {
				e.customResource.Spec.ContainerBase = &types.ContainerBase{}
			}
			e.customResource.Spec.ContainerBase.Image = repositoryWithTag
		}
	}
	if e.customResource.Spec.Image != nil && e.customResource.Spec.Image.PullPolicy != "" {
		if e.customResource.Spec.ContainerBase == nil {
			e.customResource.Spec.ContainerBase = &types.ContainerBase{}
		}
		e.customResource.Spec.ContainerBase.PullPolicy = corev1.PullPolicy(e.customResource.Spec.ContainerBase.PullPolicy)
	}

	var imagePullSecrets []corev1.LocalObjectReference
	if e.customResource.Spec.Image != nil {
		imagePullSecrets = e.customResource.Spec.Image.ImagePullSecrets
	}

	spec := appsv1.StatefulSetSpec{
		Replicas: utils.IntPointer(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: e.selectorLabels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: e.customResource.Spec.WorkloadMetaBase.Merge(metav1.ObjectMeta{
				Labels: e.selectorLabels(),
			}),
			Spec: e.customResource.Spec.WorkloadBase.Override(
				corev1.PodSpec{
					Containers: []corev1.Container{
						e.customResource.Spec.ContainerBase.Override(corev1.Container{
							Name:            config.EventTailer.TailerAffix,
							Image:           config.EventTailer.ImageWithTag,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									ReadOnly:  true,
									MountPath: "/etc/eventrouter",
								},
							},
						}),
					},
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:      utils.IntPointer64(2000),
						RunAsNonRoot: utils.BoolPointer(true),
						RunAsUser:    utils.IntPointer64(1000),
					},
					ServiceAccountName: e.Name(),
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: e.Name(),
									},
								},
							},
						},
					},
					ImagePullSecrets: imagePullSecrets,
				}),
		},
	}
	return &spec
}
