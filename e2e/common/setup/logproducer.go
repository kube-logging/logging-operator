// Copyright © 2021 Banzai Cloud
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
	"time"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/banzaicloud/logging-operator/e2e/common/cond"
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

	logProducerSource := `setInterval(() => console.log('this is a line of log', new Date()), 500);`
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
							Image: "node:latest",
							Command: []string{
								"node",
							},
							Args: []string{
								"-e",
								logProducerSource,
							},
						},
					},
				},
			},
		},
	}
	require.NoError(t, c.Create(context.Background(), &logProducerDeployment))
	require.Eventually(t, cond.AnyPodShouldBeRunning(t, c, client.MatchingLabels(lbls)), 2*time.Minute, 5*time.Second)
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
