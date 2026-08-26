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
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const release = "e2e"

var (
	producerLabels    = map[string]string{"my-unique-label": "log-producer"}
	configCheckLabels = map[string]string{"my-unique-label": "configcheck"}

	// The checksum in the message is the config that failed, so the problem can
	// only be matched by shape.
	configCheckFailure = regexp.MustCompile(`^Configuration with checksum (.+) has failed. .*`)
)

func fluentbitSpec() *v1beta1.FluentbitSpec {
	return &v1beta1.FluentbitSpec{
		Network: &v1beta1.FluentbitNetwork{
			Keepalive: new(false),
		},
		ConfigHotReload: &v1beta1.HotReload{
			Image: image.ConfigReloader().Spec(),
		},
		BufferVolumeImage: image.NodeExporter().Spec(),
	}
}

// Workers is the field the multiworker test exists to exercise, so it stays a
// parameter rather than a default.
func fluentdSpec(workers int32) *v1beta1.FluentdSpec {
	return &v1beta1.FluentdSpec{
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
		BufferVolumeMetrics: &v1beta1.Metrics{},
		Scaling: &v1beta1.FluentdScaling{
			Replicas: 1,
			Drain: v1beta1.FluentdDrainConfig{
				Enabled: true,
				Image:   image.DrainWatch().Spec(),
			},
		},
		Workers: workers,
	}
}

// readOnlyRootFilesystem is what the two config-check tests are about: the
// check has nowhere to write, so the strategy has to cope.
func readOnlyRootFilesystem() *v1beta1.Security {
	return &v1beta1.Security{
		SecurityContext: &corev1.SecurityContext{
			ReadOnlyRootFilesystem: new(true),
		},
	}
}

func logging(name, ns string, fluentd *v1beta1.FluentdSpec) *v1beta1.Logging {
	return &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: v1beta1.LoggingSpec{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			ControlNamespace: ns,
			FluentbitSpec:    fluentbitSpec(),
			FluentdSpec:      fluentd,
		},
	}
}

// The three config-check tests route to a file rather than the receiver: what
// they assert is the Logging status, not delivery.
func fileOutput(ns string) *v1beta1.Output {
	return &v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{Name: "test-output", Namespace: ns},
		Spec: v1beta1.OutputSpec{
			FileOutput: &output.FileOutputConfig{
				Path:   "/tmp/logs/${tag}/%Y/%m/%d.%H.%M",
				Append: true,
				Buffer: &output.Buffer{
					Type:        "file",
					Timekey:     "1m",
					TimekeyWait: "10s",
				},
			},
		},
	}
}

func flowTo(ns, outputName string) *v1beta1.Flow {
	return &v1beta1.Flow{
		ObjectMeta: metav1.ObjectMeta{Name: "test-flow", Namespace: ns},
		Spec: v1beta1.FlowSpec{
			Match: []v1beta1.Match{
				{Select: &v1beta1.Select{Labels: producerLabels}},
			},
			LocalOutputRefs: []string{outputName},
		},
	}
}

func TestFluentdAggregator_MultiWorker(t *testing.T) {
	const ns = "testing-1"
	testTag := "test.fluentd_aggregator_multiworker"

	env := harness.New(t).
		WithCluster("fluentd-multiworker").
		WithRelease(release).
		WithControlNamespace(ns).
		Start()

	env.Create(logging("fluentd-aggregator-multiworker-test", ns, fluentdSpec(2)))

	tags := "time"
	out := &v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{Name: "test-output", Namespace: ns},
		Spec: v1beta1.OutputSpec{
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
	}
	env.Create(out, flowTo(ns, out.Name))
	env.StartLogProducer(ns, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(ns),
	)

	env.Receiver.MustReceive(testTag)
}

func TestFluentdAggregator_ConfigChecks(t *testing.T) {
	const ns = "testing-2"
	const loggingName = "fluentd-aggregator-configchecks-test"

	env := harness.New(t).
		WithCluster("fluentd-configcheck").
		WithRelease(release).
		WithControlNamespace(ns).
		Start()

	env.Create(logging(loggingName, ns, fluentdSpec(1)))

	out := fileOutput(ns)
	env.Create(out, flowTo(ns, out.Name))
	env.StartLogProducer(ns, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(ns),
		wait.LoggingHealthy(loggingName),
	)

	t.Log("breaking the file output with an invalid config")
	patch := client.MergeFrom(out.DeepCopy())
	out.Spec.FileOutput.Path = "/tmp/zzz"
	require.NoError(t, env.Client.Patch(env.Ctx, out, patch))
	env.WaitFor(wait.LoggingProblem(loggingName, configCheckFailure.MatchString))

	t.Log("fixing it again")
	patch = client.MergeFrom(out.DeepCopy())
	out.Spec.FileOutput.Path = "/tmp/logs/${tag}/%Y/%m/%d.%H.%M"
	require.NoError(t, env.Client.Patch(env.Ctx, out, patch))
	env.WaitFor(wait.LoggingProblemCleared(loggingName, configCheckFailure.MatchString))
}

// The two differ only in the strategy: both want the check to finish and leave
// the Logging without problems, on a filesystem it cannot write to.
func TestFluentdAggregator_ConfigChecks_ReadOnlyRootFilesystem(t *testing.T) {
	for _, c := range []struct {
		name, cluster string
		check         v1beta1.ConfigCheck
	}{
		{
			"dry run",
			"fluentd-configcheck-dry-run-readonly-rootfs",
			v1beta1.ConfigCheck{Strategy: v1beta1.ConfigCheckStrategyDryRun, Labels: configCheckLabels},
		},
		{
			"start timeout",
			"fluentd-configcheck-start-timeout-readonly-rootfs",
			v1beta1.ConfigCheck{Strategy: v1beta1.ConfigCheckStrategyTimeout, TimeoutSeconds: 30, Labels: configCheckLabels},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			const ns = "testing-3"
			loggingName := "fluentd-aggregator-configchecks-ro-rootfs-test"

			spec := fluentdSpec(1)
			spec.Security = readOnlyRootFilesystem()
			spec.ConfigCheck = &c.check

			env := harness.New(t).
				WithCluster(c.cluster).
				WithRelease(release).
				WithControlNamespace(ns).
				Start()

			env.Create(logging(loggingName, ns, spec))

			out := fileOutput(ns)
			env.Create(out, flowTo(ns, out.Name))
			env.StartLogProducer(ns, producerLabels)

			env.WaitFor(
				wait.Operator(release),
				wait.Producer(producerLabels),
				wait.ConfigCheck(configCheckLabels),
				wait.FluentdAggregator(ns),
			)

			var deployed v1beta1.Logging
			require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKey{Name: loggingName}, &deployed))
			require.Equal(t, 0, deployed.Status.ProblemsCount)
		})
	}
}
