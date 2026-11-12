// Copyright © 2026 Kube logging authors
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

package fluentbit_multiline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"

	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

func TestMultilineExceptions(t *testing.T) {
	const (
		namespace = "multiline"
		release   = "multiline"
		tag       = "multiline-test"
	)
	env := harness.New(t).
		WithCluster("fluentbit-multiline").
		WithRelease(release).
		WithControlNamespace(namespace).
		Start()

	var values struct {
		Logging struct {
			Fluentbit v1beta1.FluentbitSpec `json:"fluentbit"`
		} `json:"logging"`
	}
	example, err := os.ReadFile(filepath.Join("..", "..", "..", "charts", "logging-operator", "values-multiline-example.yaml"))
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(example, &values))
	fluentbit := values.Logging.Fluentbit
	fluentbit.BufferVolumeImage = image.NodeExporter().Spec()
	fluentbit.ConfigHotReload = &v1beta1.HotReload{Image: image.ConfigReloader().Spec()}

	env.Create(&v1beta1.Logging{
		Name: "multiline",
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: namespace,
			FluentbitSpec:    &fluentbit,
			FluentdSpec: &v1beta1.FluentdSpec{
				Image:               image.Fluentd().Spec(),
				ConfigReloaderImage: image.ConfigReloader().Spec(),
				BufferVolumeImage:   image.NodeExporter().Spec(),
				DisablePvc:          true,
				Workers:             2,
				Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("50m"),
					corev1.ResourceMemory: resource.MustParse("50M"),
				}},
			},
		},
	})
	tags := "time"
	env.Create(&v1beta1.Output{
		Name: "http", Namespace: namespace,
		Spec: v1beta1.OutputSpec{
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    env.Receiver.URL(tag),
				ContentType: "application/json",
				Buffer:      &output.Buffer{Tags: &tags, Timekey: "1s", TimekeyWait: "0s"},
			},
		},
	}, &v1beta1.Flow{
		Name: "exceptions", Namespace: namespace,
		Spec: v1beta1.FlowSpec{
			Match:           []v1beta1.Match{{Select: &v1beta1.Select{Labels: map[string]string{"app": "multiline-producer"}}}},
			LocalOutputRefs: []string{"http"},
		},
	})
	env.WaitFor(wait.Operator(release), wait.Fluentbit(namespace), wait.FluentdAggregator(namespace))

	env.Create(&corev1.Pod{
		Name: "multiline-producer", Namespace: namespace,
		Labels: map[string]string{"app": "multiline-producer"},
		Spec: corev1.PodSpec{Containers: []corev1.Container{{
			Name: "producer", Image: image.Fluentd().Ref(),
			Command: []string{"sh", "-c", `while true; do
printf '%s\n' 'java.lang.RuntimeException: stdout' '    at example.Main.run(Main.java:12)'
printf '%s\n' 'Traceback (most recent call last):' '  File "app.py", line 1, in <module>' 'ValueError: stdout'
printf '%s\n' 'java.lang.RuntimeException: stderr' '    at example.Main.run(Main.java:12)' >&2
printf '%s\n' 'Traceback (most recent call last):' '  File "app.py", line 1, in <module>' 'ValueError: stderr' >&2
sleep 10
done`},
		}}},
	})
	env.WaitFor(wait.Pod(namespace, "multiline-producer"))
	env.Receiver.MustReceive(
		"java.lang.RuntimeException: stdout\n    at example.Main.run(Main.java:12)",
		"java.lang.RuntimeException: stderr\n    at example.Main.run(Main.java:12)",
		"Traceback (most recent call last):\n  File \"app.py\", line 1, in <module>\nValueError: stdout",
		"Traceback (most recent call last):\n  File \"app.py\", line 1, in <module>\nValueError: stderr",
	)
}
