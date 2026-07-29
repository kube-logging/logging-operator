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

package fixture

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// The package deliberately does not import common, so the constants are
// duplicated while both exist. This is the guard against them drifting:
// setup.defaultImages and common disagreed on the syslog-ng reloader name for
// as long as the Makefile happened to mask it.
func TestImageNamesMatchCommon(t *testing.T) {
	require.Equal(t, common.FluentdImageRepo, FluentdImageRepo)
	require.Equal(t, common.ConfigReloaderRepo, ConfigReloaderRepo)
	require.Equal(t, common.SyslogNGReloaderRepo, SyslogNGReloaderRepo)
	require.Equal(t, common.FluentdDrainWatchRepo, FluentdDrainWatchRepo)
	require.Equal(t, common.NodeExporterRepo, NodeExporterRepo)

	for _, tag := range []string{
		common.FluentdImageTag, common.ConfigReloaderTag, common.SyslogNGReloaderTag,
		common.FluentdDrainWatchTag, common.NodeExporterTag,
	} {
		require.Equal(t, localTag, tag)
	}
}

func TestLoggingDefaults(t *testing.T) {
	l := Logging("ns", "lg")

	require.Equal(t, "lg", l.Name)
	require.Equal(t, "ns", l.Namespace)
	require.Equal(t, "ns", l.Spec.ControlNamespace)
	require.True(t, l.Spec.EnableRecreateWorkloadOnImmutableFieldChange)

	// Neither half is present until asked for: suites that want only an agent
	// must not silently get an aggregator.
	require.Nil(t, l.Spec.FluentbitSpec)
	require.Nil(t, l.Spec.FluentdSpec)
}

func TestWithFluentbit(t *testing.T) {
	l := Logging("ns", "lg", WithFluentbit())
	fb := l.Spec.FluentbitSpec

	require.NotNil(t, fb)
	require.False(t, *fb.Network.Keepalive)
	require.Equal(t, "config-reloader", fb.ConfigHotReload.Image.Repository)
	require.Equal(t, "node-exporter", fb.BufferVolumeImage.Repository)
	require.Equal(t, "local", fb.BufferVolumeImage.Tag)
}

func TestWithFluentdDefaults(t *testing.T) {
	l := Logging("ns", "lg", WithFluentd())
	fd := l.Spec.FluentdSpec

	require.NotNil(t, fd)
	require.EqualValues(t, 2, fd.Workers, "the counted majority is 2, not 1")
	require.Equal(t, FluentdImageRepo, fd.Image.Repository)
	require.Equal(t, ConfigReloaderRepo, fd.ConfigReloaderImage.Repository)
	require.Equal(t, NodeExporterRepo, fd.BufferVolumeImage.Repository)
	require.NotNil(t, fd.BufferVolumeMetrics)

	require.Equal(t, resource.MustParse("500m"), fd.Resources.Limits[corev1.ResourceCPU])
	require.Equal(t, resource.MustParse("200M"), fd.Resources.Limits[corev1.ResourceMemory])
	require.Equal(t, resource.MustParse("250m"), fd.Resources.Requests[corev1.ResourceCPU])
	require.Equal(t, resource.MustParse("50M"), fd.Resources.Requests[corev1.ResourceMemory])

	// Draining is opt-in: it is absent unless Drain() is passed.
	require.Nil(t, fd.Scaling)
}

func TestFluentdOptions(t *testing.T) {
	fd := FluentdSpec(Workers(1), Drain())

	require.EqualValues(t, 1, fd.Workers)
	require.NotNil(t, fd.Scaling)
	require.EqualValues(t, 1, fd.Scaling.Replicas)
	require.True(t, fd.Scaling.Drain.Enabled)
	require.Equal(t, FluentdDrainWatchRepo, fd.Scaling.Drain.Image.Repository)
}

func TestBufferDefaultsAndOverrides(t *testing.T) {
	b := Buffer("time")
	require.Equal(t, "file", b.Type)
	require.Equal(t, "time", *b.Tags)
	// Literals, not the constants: asserting a constant against itself proves
	// nothing, and these two are not pinned by the literal test below.
	require.Equal(t, "1s", b.Timekey, "counted majority: 1s at eight sites")
	require.Equal(t, "0s", b.TimekeyWait, "counted majority: 0s at nine sites")

	o := Buffer("time", Timekey("1m"), TimekeyWait("10s"))
	require.Equal(t, "1m", o.Timekey)
	require.Equal(t, "10s", o.TimekeyWait)
}

func TestReceiverURL(t *testing.T) {
	require.Equal(t, "http://e2e-test-receiver:8080/test.tag", ReceiverURL("e2e", "test.tag"))
}

func TestHTTPOutputAndFlow(t *testing.T) {
	out := HTTPOutput("ns", "test-output", ReceiverURL("e2e", "tag"), Buffer("time"))
	require.Equal(t, "ns", out.Namespace)
	require.Equal(t, "application/json", out.Spec.HTTPOutput.ContentType)
	require.Equal(t, "http://e2e-test-receiver:8080/tag", out.Spec.HTTPOutput.Endpoint)
	require.Equal(t, "1s", out.Spec.HTTPOutput.Buffer.Timekey)

	labels := map[string]string{"my-unique-label": "log-producer"}
	f := Flow("ns", "test-flow", labels, out.Name)
	require.Equal(t, labels, f.Spec.Match[0].Select.Labels)
	require.Equal(t, []string{"test-output"}, f.Spec.LocalOutputRefs)
}

// The §6 safety net: build the Logging that fluentd-aggregator's MultiWorker
// test writes by hand and assert the builder produces exactly it. A default
// that quietly absorbed a divergence would show up here rather than after a
// twenty minute cluster run.
func TestBuilderReproducesTheMultiWorkerLiteral(t *testing.T) {
	const ns = "testing-1"

	want := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-aggregator-multiworker-test",
			Namespace: ns,
		},
		Spec: v1beta1.LoggingSpec{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			ControlNamespace: ns,
			FluentbitSpec: &v1beta1.FluentbitSpec{
				Network: &v1beta1.FluentbitNetwork{Keepalive: new(false)},
				ConfigHotReload: &v1beta1.HotReload{
					Image: v1beta1.ImageSpec{Repository: common.ConfigReloaderRepo, Tag: common.ConfigReloaderTag},
				},
				BufferVolumeImage: v1beta1.ImageSpec{Repository: common.NodeExporterRepo, Tag: common.NodeExporterTag},
			},
			FluentdSpec: &v1beta1.FluentdSpec{
				Image:               v1beta1.ImageSpec{Repository: common.FluentdImageRepo, Tag: common.FluentdImageTag},
				ConfigReloaderImage: v1beta1.ImageSpec{Repository: common.ConfigReloaderRepo, Tag: common.ConfigReloaderTag},
				BufferVolumeImage:   v1beta1.ImageSpec{Repository: common.NodeExporterRepo, Tag: common.NodeExporterTag},
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
						Image:   v1beta1.ImageSpec{Repository: common.FluentdDrainWatchRepo, Tag: common.FluentdDrainWatchTag},
					},
				},
				Workers: 2,
			},
		},
	}

	got := Logging(ns, "fluentd-aggregator-multiworker-test",
		WithFluentbit(),
		WithFluentd(Drain()),
	)

	require.Equal(t, want, got)
}
