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
	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/extensionsconfig"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	spec := appsv1.StatefulSetSpec{
		Replicas: utils.IntPointer(1),
		Selector: &v1.LabelSelector{
			MatchLabels: e.selectorLabels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: e.customResource.Spec.WorkloadMetaBase.Merge(v1.ObjectMeta{
				Labels: e.selectorLabels(),
			}),
			Spec: e.customResource.Spec.WorkloadBase.Override(
				corev1.PodSpec{
					Containers: []corev1.Container{
						e.customResource.Spec.ContainerBase.Override(corev1.Container{
							Name:            config.EventTailer.TailerAffix,
							Image:           "banzaicloud/eventrouter:v0.1.0",
							ImagePullPolicy: corev1.PullIfNotPresent,
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
				}),
		},
	}
	return &spec
}
