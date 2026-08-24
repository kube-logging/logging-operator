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

package fluentd_aggregator

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	clusterName = "fluentd-aggregator"
	release     = "e2e"
	ns          = "testing-1"
	labeledNs   = "labeled-namespace"
	testTag     = "test.fluentd_aggregator_nslabel"
)

var producerLabels = map[string]string{"my-unique-label": "log-producer"}

func TestFluentdAggregator_NamespaceLabel(t *testing.T) {
	env := harness.New(t).
		WithCluster(clusterName).
		WithRelease(release).
		WithControlNamespace(ns).
		WithNamespaces(labeledNs).
		Start()

	env.Create(&v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: "fluentd-aggregator-nslabel-test"},
		Spec: v1beta1.LoggingSpec{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			ControlNamespace: ns,
			FluentbitSpec: &v1beta1.FluentbitSpec{
				Network: &v1beta1.FluentbitNetwork{
					Keepalive: new(false),
				},
				ConfigHotReload: &v1beta1.HotReload{
					Image: image.ConfigReloader().Spec(),
				},
				BufferVolumeImage: image.NodeExporter().Spec(),
				FilterKubernetes:  v1beta1.FilterKubernetes{
					// Namespace labels enrichment is enabled by default starting with version 4.9
					// NamespaceLabels: "On",
				},
			},
			FluentdSpec: &v1beta1.FluentdSpec{
				Image:               image.Fluentd().Spec(),
				ConfigReloaderImage: image.ConfigReloader().Spec(),
				BufferVolumeImage:   image.NodeExporter().Spec(),
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("200M"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("250m"),
						corev1.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	})

	tags := "time"
	out := &v1beta1.ClusterOutput{
		ObjectMeta: metav1.ObjectMeta{Name: "test-output", Namespace: ns},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    env.Receiver.URL(testTag),
					ContentType: "application/json",
					Buffer: &output.Buffer{
						Type:        "file",
						Tags:        &tags,
						Timekey:     "1s",
						TimekeyWait: "0s",
					},
				},
			},
		},
	}
	env.Create(out)

	env.Create(&v1beta1.ClusterFlow{
		ObjectMeta: metav1.ObjectMeta{Name: "test-flow", Namespace: ns},
		Spec: v1beta1.ClusterFlowSpec{
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						NamespaceLabels: map[string]string{
							"kubernetes.io/metadata.name": labeledNs,
						},
					},
				},
			},
			GlobalOutputRefs: []string{out.Name},
		},
	})

	env.StartLogProducer(labeledNs, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(ns),
	)

	env.Receiver.MustReceive(testTag)
}
