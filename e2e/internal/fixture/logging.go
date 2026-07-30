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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// Defaults are the value most call sites already use, counted across the
// suites rather than chosen. Anything that diverges passes an option; a
// default must never absorb a divergence.
const (
	// DefaultWorkers is 2, used by four call sites against three using 1.
	DefaultWorkers = 2
)

// LoggingOption modifies a Logging under construction.
type LoggingOption func(*v1beta1.Logging)

// Logging builds the Logging every suite starts from: recreate-on-immutable
// enabled, control namespace set, and no agent or aggregator until asked.
// Logging is cluster-scoped, so controlNamespace sets Spec.ControlNamespace and
// no ObjectMeta.Namespace.
func Logging(controlNamespace, name string, opts ...LoggingOption) *v1beta1.Logging {
	l := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: v1beta1.LoggingSpec{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			ControlNamespace: controlNamespace,
		},
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

// WithFluentbit adds the agent spec the suites share: keepalive off, hot
// reload and buffer volume on the local images.
func WithFluentbit(opts ...func(*v1beta1.FluentbitSpec)) LoggingOption {
	return func(l *v1beta1.Logging) {
		s := &v1beta1.FluentbitSpec{
			Network:           &v1beta1.FluentbitNetwork{Keepalive: new(false)},
			ConfigHotReload:   &v1beta1.HotReload{Image: image(ConfigReloaderRepo)},
			BufferVolumeImage: image(NodeExporterRepo),
		}
		for _, o := range opts {
			o(s)
		}
		l.Spec.FluentbitSpec = s
	}
}

// WithFluentd adds the aggregator spec the suites share.
func WithFluentd(opts ...FluentdOption) LoggingOption {
	return func(l *v1beta1.Logging) {
		l.Spec.FluentdSpec = FluentdSpec(opts...)
	}
}

// FluentdOption modifies a FluentdSpec under construction. It is exported so
// the same options serve Logging.FluentdSpec and a detached FluentdConfig.
type FluentdOption func(*v1beta1.FluentdSpec)

// FluentdSpec builds the aggregator spec shared by the Logging-embedded and
// detached forms.
func FluentdSpec(opts ...FluentdOption) *v1beta1.FluentdSpec {
	s := &v1beta1.FluentdSpec{
		Image:               image(FluentdImageRepo),
		ConfigReloaderImage: image(ConfigReloaderRepo),
		BufferVolumeImage:   image(NodeExporterRepo),
		Resources:           defaultResources(),
		BufferVolumeMetrics: &v1beta1.Metrics{},
		Workers:             DefaultWorkers,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Workers overrides the worker count. Three call sites use 1.
func Workers(n int32) FluentdOption {
	return func(s *v1beta1.FluentdSpec) { s.Workers = n }
}

// Drain enables buffer draining on downscale, with one replica.
func Drain() FluentdOption {
	return func(s *v1beta1.FluentdSpec) {
		s.Scaling = &v1beta1.FluentdScaling{
			Replicas: 1,
			Drain: v1beta1.FluentdDrainConfig{
				Enabled: true,
				Image:   image(FluentdDrainWatchRepo),
			},
		}
	}
}

// defaultResources is the limits/requests pair used by twelve of the CR call
// sites. LoggingInfra and the syslog-ng overrides use different values and set
// them explicitly.
func defaultResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("200M"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("250m"),
			corev1.ResourceMemory: resource.MustParse("50M"),
		},
	}
}
