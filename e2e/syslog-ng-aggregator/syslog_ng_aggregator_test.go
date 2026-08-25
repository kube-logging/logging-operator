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

package syslong_ng_aggregator

import (
	"testing"

	"github.com/cisco-open/operator-tools/pkg/typeoverride"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/resources/syslogng"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	syslogngoutput "github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
)

const (
	ns      = "test"
	release = "e2e"
	testTag = "test.tag"
)

var producerLabels = map[string]string{"my-unique-label": "log-producer"}

func TestSyslogNGIsRunningAndForwardingLogs(t *testing.T) {
	env := harness.New(t).
		WithCluster("syslog-ng-forwarding").
		WithRelease(release).
		WithControlNamespace(ns).
		Start()

	logging := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "syslog-ng-aggregator-test",
			Namespace: ns,
		},
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
			SyslogNGSpec: &v1beta1.SyslogNGSpec{
				ConfigReloadImage:        image.SyslogNGReloader().Basic(),
				BufferVolumeMetricsImage: image.NodeExporter().Basic(),
				StatefulSetOverrides: &typeoverride.StatefulSet{
					Spec: typeoverride.StatefulSetSpec{
						Template: typeoverride.PodTemplateSpec{
							Spec: typeoverride.PodSpec{
								Containers: []corev1.Container{
									{
										Name: syslogng.ContainerName,
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("100m"),
												corev1.ResourceMemory: resource.MustParse("100M"),
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("25m"),
												corev1.ResourceMemory: resource.MustParse("10M"),
											},
										},
										VolumeMounts: []corev1.VolumeMount{
											{
												Name:      "buffers",
												MountPath: "/buffers",
											},
										},
									},
								},
								Volumes: []corev1.Volume{
									{
										Name: "buffers",
										VolumeSource: corev1.VolumeSource{
											EmptyDir: &corev1.EmptyDirVolumeSource{},
										},
									},
								},
							},
						},
					},
				},
				BufferVolumeMetrics: &v1beta1.BufferMetrics{
					Metrics: v1beta1.Metrics{
						Interval: "1s",
					},
					MountName: "buffers",
				},
			},
		},
	}

	output := &v1beta1.SyslogNGOutput{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-output",
			Namespace: ns,
		},
		Spec: v1beta1.SyslogNGOutputSpec{
			HTTP: &syslogngoutput.HTTPOutput{
				URL: env.Receiver.URL(testTag),
				Headers: []string{
					"Content-type: application/json",
				},
				Method: "POST",
				DiskBuffer: &syslogngoutput.DiskBuffer{
					DiskBufSize: 100 * 1024 * 1024,
					Reliable:    true,
					Dir:         syslogng.BufferPath,
				},
			},
		},
	}

	flow := &v1beta1.SyslogNGFlow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-flow",
			Namespace: ns,
		},
		Spec: v1beta1.SyslogNGFlowSpec{
			Match: &v1beta1.SyslogNGMatch{
				Regexp: &filter.RegexpMatchExpr{
					Pattern: "log-producer",
					Value:   "json.kubernetes.labels.my-unique-label",
					Type:    "string",
				},
			},
			LocalOutputRefs: []string{output.Name},
		},
	}

	env.Create(logging, output, flow)
	env.StartLogProducer(ns, producerLabels)

	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.SyslogNGAggregator(ns),
	)
	env.Receiver.MustReceive(testTag)
}
