// Copyright © 2025 Kube logging authors
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

package fluentbit_agent_namespace

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
	release    = "fluentbit-agents-namespace"
	nsControl  = "logging"
	nsAgents   = "logging-fluentbit-agents"
	nsProducer = "default"
	tag        = "tag_fluentbit_agents"
)

var producerLabels = map[string]string{"my-unique-label": "log-producer"}

func realTimeBuffer() *output.Buffer {
	tags := "time"
	return &output.Buffer{Tags: &tags, Timekey: "1s", TimekeyWait: "0s"}
}

// Verifies that Fluentbit agents are deployed to a dedicated namespace when
// logging.spec.fluentBitAgentNamespace is set, while the aggregator stays in
// the control namespace.
func TestFluentbitAgentDedicatedNamespace(t *testing.T) {
	env := harness.New(t).
		WithCluster(release).
		WithRelease(release).
		WithControlNamespace(nsControl).
		WithNamespaces(nsAgents).
		Start()

	out := &v1beta1.ClusterOutput{
		ObjectMeta: metav1.ObjectMeta{Name: "http", Namespace: nsControl},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				LoggingRef: "infra",
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    env.Receiver.URL(tag),
					ContentType: "application/json",
					Buffer:      realTimeBuffer(),
				},
			},
		},
	}
	env.Create(out)

	env.Create(&v1beta1.ClusterFlow{
		ObjectMeta: metav1.ObjectMeta{Name: "flow", Namespace: nsControl},
		Spec: v1beta1.ClusterFlowSpec{
			LoggingRef: "infra",
			Match: []v1beta1.ClusterMatch{
				{ClusterSelect: &v1beta1.ClusterSelect{Labels: producerLabels}},
			},
			GlobalOutputRefs: []string{out.Name},
		},
	})

	env.Create(&v1beta1.FluentbitAgent{
		ObjectMeta: metav1.ObjectMeta{Name: "infra"},
		Spec: v1beta1.FluentbitSpec{
			LoggingRef: "infra",
			ConfigHotReload: &v1beta1.HotReload{
				Image: image.ConfigReloader().Spec(),
			},
			BufferVolumeImage: image.NodeExporter().Spec(),
		},
	})

	env.Create(&v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: "infra"},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:              "infra",
			ControlNamespace:        nsControl,
			FluentbitAgentNamespace: nsAgents,
			FluentdSpec: &v1beta1.FluentdSpec{
				Image:               image.Fluentd().Spec(),
				ConfigReloaderImage: image.ConfigReloader().Spec(),
				BufferVolumeImage:   image.NodeExporter().Spec(),
				DisablePvc:          true,
				Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("50m"),
					corev1.ResourceMemory: resource.MustParse("50M"),
				}},
			},
		},
	})

	env.StartLogProducer(nsProducer, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(nsControl),
		wait.Fluentbit(nsAgents),
	)

	env.Receiver.MustReceive(tag)
}
