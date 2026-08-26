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

package volumedrain

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	release     = "volumedrain"
	loggingName = "drainer-test"
	testTag     = "test.volumedrain"

	// The replica the Logging drops to one, so the drainer runs against it.
	drainedReplica = loggingName + "-fluentd-1"
	drainerJob     = drainedReplica + "-drainer"
	drainedPVC     = loggingName + "-fluentd-buffer-" + drainedReplica
)

var producerLabels = map[string]string{"my-unique-label": "log-producer"}

// deleteVolume is what separates the two tests: whether the drained claim is
// kept and labeled, or removed with the replica.
func drainingLogging(ns string, deleteVolume bool) *v1beta1.Logging {
	return &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: loggingName, Namespace: ns},
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
				BufferVolumeMetrics: &v1beta1.Metrics{Enabled: new(true)},
				Scaling: &v1beta1.FluentdScaling{
					Replicas: 2,
					Drain: v1beta1.FluentdDrainConfig{
						Enabled:      true,
						DeleteVolume: deleteVolume,
						Image:        image.DrainWatch().Spec(),
					},
				},
			},
		},
	}
}

func httpOutput(env *harness.Env, ns, timekey string) *v1beta1.Output {
	tags := "time"
	return &v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{Name: "test-output", Namespace: ns},
		Spec: v1beta1.OutputSpec{
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    env.Receiver.URL(testTag),
				ContentType: "application/json",
				Buffer: &output.Buffer{
					Type:        "file",
					Tags:        &tags,
					Timekey:     timekey,
					TimekeyWait: "0s",
				},
			},
		},
	}
}

func flowToOutput(ns, outputName string) *v1beta1.Flow {
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

// bufferFiles counts what the aggregator has queued. This is still a shell-out:
// reading a path inside a container is the kubectl exec seam tracked in #2325,
// which the receiver helpers do not cover.
func bufferFiles(env *harness.Env, ns string) (int, error) {
	out, err := common.CmdEnv(exec.Command("kubectl",
		"-n", ns, "exec", drainedReplica, "-c", "fluentd", "--", "ls", "-1", "/buffers"), env.Cluster).Output()
	if err != nil {
		return 0, err
	}
	return strings.Count(string(out), "\n"), nil
}

// How long each step of the arc is allowed. They are the suite's own numbers:
// settling within thirty seconds and settling within the shared five minutes
// are different claims, and the second would not catch a drain-speed
// regression.
type drainBudget struct {
	buffered, bufferedTick time.Duration
	started, startedTick   time.Duration
	podGoneTick            time.Duration
}

// drainOneReplica is the shared arc: fill the buffers with the receiver away,
// drop to one replica, and let the drainer finish once it is back.
func drainOneReplica(t *testing.T, env *harness.Env, ns string, logging *v1beta1.Logging, b drainBudget) {
	t.Helper()

	env.Receiver.Scale(0)

	require.Eventually(t, func() bool {
		n, err := bufferFiles(env, ns)
		if err != nil {
			t.Logf("listing the buffer directory: %v", err)
			return false
		}
		return n > 2
	}, b.buffered, b.bufferedTick, "the aggregator never buffered anything")

	patch := client.MergeFrom(logging.DeepCopy())
	logging.Spec.FluentdSpec.Scaling.Replicas = 1
	require.NoError(t, env.Client.Patch(env.Ctx, logging, patch))

	env.WaitWithin(b.started, b.startedTick, wait.JobStarted(ns, drainerJob))
	env.WaitWithin(30*time.Second, time.Second/2, wait.Pod(ns, drainedReplica))

	env.Receiver.Scale(1)

	env.WaitWithin(3*time.Minute, 3*time.Second, wait.JobGone(ns, drainerJob))
	env.WaitWithin(30*time.Second, b.podGoneTick, wait.PodGone(ns, drainedReplica))
}

func TestVolumeDrain_Downscale(t *testing.T) {
	const ns = "testing-1"

	env := harness.New(t).
		WithCluster("drain").
		WithRelease(release).
		WithControlNamespace(ns).
		Start()

	logging := drainingLogging(ns, false)
	env.Create(logging)

	out := httpOutput(env, ns, "1s")
	env.Create(out, flowToOutput(ns, out.Name))
	env.StartLogProducer(ns, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(ns),
	)
	env.Receiver.MustReceive(testTag)

	drainOneReplica(t, env, ns, logging, drainBudget{
		buffered: time.Minute, bufferedTick: 3 * time.Second,
		started: time.Minute, startedTick: 3 * time.Second,
		podGoneTick: time.Second,
	})

	// The claim outlives the replica, carrying the operator's verdict.
	var pvc corev1.PersistentVolumeClaim
	require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKey{Namespace: ns, Name: drainedPVC}, &pvc))
	assert.Equal(t, "drained", pvc.GetLabels()["logging.banzaicloud.io/drain-status"])
}

func TestVolumeDrain_Downscale_DeleteVolume(t *testing.T) {
	const ns = "testing-2"

	env := harness.New(t).
		WithCluster("drain-2").
		WithRelease(release).
		WithControlNamespace(ns).
		Start()

	logging := drainingLogging(ns, true)
	env.Create(logging)

	out := httpOutput(env, ns, "10s")
	env.Create(out, flowToOutput(ns, out.Name))
	env.StartLogProducer(ns, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(ns),
	)
	env.Receiver.MustReceive(testTag)

	drainOneReplica(t, env, ns, logging, drainBudget{
		buffered: 3 * time.Minute, bufferedTick: 10 * time.Second,
		started: 2 * time.Minute, startedTick: time.Second,
		podGoneTick: time.Second / 2,
	})

	// deleteVolume, so the claim goes with the replica rather than being kept.
	env.WaitWithin(30*time.Second, time.Second/2, wait.PVCGone(ns, drainedPVC))
}
