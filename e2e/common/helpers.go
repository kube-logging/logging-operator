// Copyright © 2023 Kube logging authors
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

package common

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	FluentdImageRepo      = "fluentd-full"
	FluentdImageTag       = "local"
	ConfigReloaderRepo    = "config-reloader"
	ConfigReloaderTag     = "local"
	SyslogNGReloaderRepo  = "syslog-ng-reloader"
	SyslogNGReloaderTag   = "local"
	FluentdDrainWatchRepo = "fluentd-drain-watch"
	FluentdDrainWatchTag  = "local"
	NodeExporterRepo      = "node-exporter"
	NodeExporterTag       = "local"
)

var sequence uint32

func RequireNoError(t *testing.T, err error) {
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Received unexpected error:\n%#v %+v", err, errors.GetDetails(err)))
		t.FailNow()
	}
}

func Initialize(t *testing.T) {
	localSeq := atomic.AddUint32(&sequence, 1)
	shards := cast.ToUint32(os.Getenv("SHARDS"))
	shard := cast.ToUint32(os.Getenv("SHARD"))
	if shards > 0 {
		if localSeq%shards != shard {
			t.Skipf("skipping %s as sequence %d not in shard %d", t.Name(), localSeq, shard)
		}
	}
	t.Parallel()
}

func LoggingInfra(
	ctx context.Context,
	t *testing.T,
	c client.Client,
	nsInfra string,
	release string,
	tag string,
	buffer *output.Buffer,
	producerLabels map[string]string) {

	output := v1beta1.ClusterOutput{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http",
			Namespace: nsInfra,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				LoggingRef: "infra",
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    fmt.Sprintf("http://%s-test-receiver:8080/%s", release, tag),
					ContentType: "application/json",
					Buffer:      buffer,
				},
			},
		},
	}

	RequireNoError(t, c.Create(ctx, &output))
	flow := v1beta1.ClusterFlow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flow",
			Namespace: nsInfra,
		},
		Spec: v1beta1.ClusterFlowSpec{
			LoggingRef: "infra",
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						Labels: producerLabels,
					},
				},
			},
			GlobalOutputRefs: []string{output.Name},
		},
	}
	RequireNoError(t, c.Create(ctx, &flow))

	agent := v1beta1.FluentbitAgent{
		ObjectMeta: metav1.ObjectMeta{
			Name: "infra",
		},
		Spec: v1beta1.FluentbitSpec{
			LoggingRef: "infra",
			ConfigHotReload: &v1beta1.HotReload{
				Image: v1beta1.ImageSpec{
					Repository: ConfigReloaderRepo,
					Tag:        ConfigReloaderTag,
				},
			},
			BufferVolumeImage: v1beta1.ImageSpec{
				Repository: NodeExporterRepo,
				Tag:        NodeExporterTag,
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &agent))

	logging := v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name: "infra",
			Labels: map[string]string{
				"tenant": "infra",
			},
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       "infra",
			ControlNamespace: nsInfra,
			FluentdSpec: &v1beta1.FluentdSpec{
				Image: v1beta1.ImageSpec{
					Repository: FluentdImageRepo,
					Tag:        FluentdImageTag,
				},
				ConfigReloaderImage: v1beta1.ImageSpec{
					Repository: ConfigReloaderRepo,
					Tag:        ConfigReloaderTag,
				},
				BufferVolumeImage: v1beta1.ImageSpec{
					Repository: NodeExporterRepo,
					Tag:        NodeExporterTag,
				},
				DisablePvc: true,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &logging))
}

func LoggingTenant(
	ctx context.Context,
	t *testing.T,
	c client.Client,
	nsTenant,
	nsInfra,
	release,
	tag string,
	buffer *output.Buffer,
	producerLabels map[string]string) {
	output := v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http",
			Namespace: nsTenant,
		},
		Spec: v1beta1.OutputSpec{
			LoggingRef: "tenant",
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    fmt.Sprintf("http://%s-test-receiver.%s:8080/%s", release, nsInfra, tag),
				ContentType: "application/json",
				Buffer:      buffer,
			},
		},
	}

	RequireNoError(t, c.Create(ctx, &output))
	flow := v1beta1.Flow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flow",
			Namespace: nsTenant,
		},
		Spec: v1beta1.FlowSpec{
			LoggingRef: "tenant",
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{
						Labels: producerLabels,
					},
				},
			},
			LocalOutputRefs: []string{output.Name},
		},
	}
	RequireNoError(t, c.Create(ctx, &flow))

	logging := v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tenant",
			Labels: map[string]string{
				"tenant": "tenant",
			},
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       "tenant",
			ControlNamespace: nsTenant,
			WatchNamespaces:  []string{"tenant"},
			FluentdSpec: &v1beta1.FluentdSpec{
				Image: v1beta1.ImageSpec{
					Repository: FluentdImageRepo,
					Tag:        FluentdImageTag,
				},
				ConfigReloaderImage: v1beta1.ImageSpec{
					Repository: ConfigReloaderRepo,
					Tag:        ConfigReloaderTag,
				},
				BufferVolumeImage: v1beta1.ImageSpec{
					Repository: NodeExporterRepo,
					Tag:        NodeExporterTag,
				},
				DisablePvc: true,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &logging))
}

func LoggingRoute(ctx context.Context, t *testing.T, c client.Client) {
	ap := v1beta1.LoggingRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tenants",
		},
		Spec: v1beta1.LoggingRouteSpec{
			Source: "infra",
			Targets: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "tenant",
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &ap))
}

// WaitForPodReady waits for a pod to be in Running phase and Ready condition
func WaitForPodReady(ctx context.Context, c client.Client, pod *corev1.Pod, pollInterval, pollTimeout time.Duration) error {
	return wait.PollUntilContextTimeout(ctx, pollInterval, pollTimeout, true, wait.ConditionWithContextFunc(func(ctx context.Context) (bool, error) {
		var updatedPod corev1.Pod
		err := c.Get(ctx, client.ObjectKeyFromObject(pod), &updatedPod)
		if client.IgnoreNotFound(err) != nil {
			return false, fmt.Errorf("failed to get pod status: %w", err)
		}

		isReady := updatedPod.Status.Phase == corev1.PodRunning
		for _, cond := range updatedPod.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return isReady, nil
	}))
}

// SetupCurlPod creates a curl pod for testing HTTP endpoints and waits for it to be ready
func SetupCurlPod(ctx context.Context, c client.Client, namespace, name string, pollInterval, pollTimeout time.Duration) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "curl",
					Image:   "curlimages/curl:latest",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}

	if err := c.Create(ctx, pod); err != nil {
		return nil, fmt.Errorf("failed to create curl pod: %w", err)
	}

	if err := WaitForPodReady(ctx, c, pod, pollInterval, pollTimeout); err != nil {
		return nil, fmt.Errorf("failed to wait for curl pod to be ready: %w", err)
	}

	return pod, nil
}
