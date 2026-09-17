// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package harness

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/cisco-open/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	logProducerName   = "log-producer"
	logProducerConfig = "log-generator-config"
)

// logProducer returns the ConfigMap first: the Deployment mounts it.
func logProducer(namespace string, extraLabels map[string]string) []client.Object {
	lbls := utils.MergeLabels(map[string]string{types.NameLabel: logProducerName}, extraLabels)
	return []client.Object{
		&corev1.ConfigMap{
			Name:      logProducerConfig,
			Namespace: namespace,
			Data: map[string]string{
				"config.ini": heredoc.Doc(`
					[message]
					count = -1

					[golang]
					enabled = true
				`),
			},
		},
		&appsv1.Deployment{
			Name:      logProducerName,
			Namespace: namespace,
			Spec: appsv1.DeploymentSpec{
				Replicas: new(int32(1)),
				Selector: metav1.SetAsLabelSelector(labels.Set(lbls)),
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: lbls,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "producer",
								Image: "ghcr.io/kube-logging/log-generator:latest",
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "config",
										MountPath: "/conf",
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "config",
								ConfigMap: &corev1.ConfigMapVolumeSource{
									Name: logProducerConfig,
								},
							},
						},
					},
				},
			},
		},
	}
}
