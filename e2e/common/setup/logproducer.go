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

package setup

import (
	"context"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cisco-open/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/e2e/common"
)

func LogProducer(t *testing.T, c client.Client, opts ...LogProducerOption) {
	options := LogProducerOptions{
		Name:      "log-producer",
		Namespace: "default",
		Replicas:  1,
	}
	for _, opt := range opts {
		opt.ApplyToLogProducerOptions(&options)
	}

	lbls := map[string]string{
		"app.kubernetes.io/name": options.Name,
	}
	lbls = utils.MergeLabels(lbls, options.Labels)
	logProducerDeployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: options.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.IntPointer(int32(options.Replicas)),
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
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "log-generator-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	logProducerConfig := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "log-generator-config",
			Namespace: options.Namespace,
		},
		Data: map[string]string{
			"config.ini": heredoc.Doc(`
				[message]
				count = -1
				
				[golang]
				enabled = true
			`),
		},
	}
	common.RequireNoError(t, c.Create(context.Background(), &logProducerConfig))
	common.RequireNoError(t, c.Create(context.Background(), &logProducerDeployment))
}

type LogProducerOption interface {
	ApplyToLogProducerOptions(*LogProducerOptions)
}

type LogProducerOptionFunc func(*LogProducerOptions)

func (fn LogProducerOptionFunc) ApplyToLogProducerOptions(options *LogProducerOptions) {
	fn(options)
}

type LogProducerOptions struct {
	Labels    map[string]string
	Name      string
	Namespace string
	Replicas  int
}
