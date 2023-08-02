// Copyright © 2019 Banzai Cloud
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

package v1beta1

import (
	"github.com/cisco-open/operator-tools/pkg/typeoverride"
	"github.com/cisco-open/operator-tools/pkg/volume"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/input"
	corev1 "k8s.io/api/core/v1"
)

// +name:"FluentdSpec"
// +weight:"200"
type _hugoFluentdSpec interface{} //nolint:deadcode,unused

// +name:"FluentdSpec"
// +version:"v1beta1"
// +description:"FluentdSpec defines the desired state of Fluentd"
type _metaFluentdSpec interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true

// FluentdSpec defines the desired state of Fluentd
type FluentdSpec struct {
	Affinity    *corev1.Affinity  `json:"affinity,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	// BufferStorageVolume is by default configured as PVC using FluentdPvcSpec
	// +docLink:"volume.KubernetesVolume,https://github.com/cisco-open/operator-tools/tree/master/docs/types"
	BufferStorageVolume     volume.KubernetesVolume     `json:"bufferStorageVolume,omitempty"`
	BufferVolumeMetrics     *Metrics                    `json:"bufferVolumeMetrics,omitempty"`
	BufferVolumeImage       ImageSpec                   `json:"bufferVolumeImage,omitempty"`
	BufferVolumeArgs        []string                    `json:"bufferVolumeArgs,omitempty"`
	ConfigCheckAnnotations  map[string]string           `json:"configCheckAnnotations,omitempty"`
	ConfigReloaderImage     ImageSpec                   `json:"configReloaderImage,omitempty"`
	ConfigCheckResources    corev1.ResourceRequirements `json:"configCheckResources,omitempty"`
	ConfigReloaderResources corev1.ResourceRequirements `json:"configReloaderResources,omitempty"`
	CompressConfigFile      bool                        `json:"compressConfigFile,omitempty"`
	DisablePvc              bool                        `json:"disablePvc,omitempty"`
	DNSPolicy               corev1.DNSPolicy            `json:"dnsPolicy,omitempty"`
	DNSConfig               *corev1.PodDNSConfig        `json:"dnsConfig,omitempty"`
	// Allows Time object in buffer's MessagePack serde
	// +docLink:"more info, https://docs.fluentd.org/deployment/system-config#enable_msgpack_time_support"
	EnableMsgpackTimeSupport bool            `json:"enableMsgpackTimeSupport,omitempty"`
	EnvVars                  []corev1.EnvVar `json:"envVars,omitempty"`
	ExtraArgs                []string        `json:"extraArgs,omitempty"`
	ExtraVolumes             []ExtraVolume   `json:"extraVolumes,omitempty"`
	// Deprecated, use bufferStorageVolume
	FluentdPvcSpec *volume.KubernetesVolume `json:"fluentdPvcSpec,omitempty"`
	// +kubebuilder:validation:enum=stdout,null
	FluentLogDestination string `json:"fluentLogDestination,omitempty"`
	// FluentOutLogrotate sends fluent's stdout to file and rotates it
	FluentOutLogrotate *FluentOutLogrotate       `json:"fluentOutLogrotate,omitempty"`
	ForwardInputConfig *input.ForwardInputConfig `json:"forwardInputConfig,omitempty"`
	// Ignore same log lines
	// +docLink:"more info, https://docs.fluentd.org/deployment/logging#ignore_same_log_interval"
	IgnoreSameLogInterval string `json:"ignoreSameLogInterval,omitempty"`
	// Ignore repeated log lines
	// +docLink:"more info, https://docs.fluentd.org/deployment/logging#ignore_repeated_log_interval"
	IgnoreRepeatedLogInterval string            `json:"ignoreRepeatedLogInterval,omitempty"`
	Image                     ImageSpec         `json:"image,omitempty"`
	Labels                    map[string]string `json:"labels,omitempty"`
	LivenessProbe             *corev1.Probe     `json:"livenessProbe,omitempty"`
	LivenessDefaultCheck      bool              `json:"livenessDefaultCheck,omitempty"`
	// +kubebuilder:validation:enum=fatal,error,warn,info,debug,trace
	LogLevel                  string                            `json:"logLevel,omitempty"`
	Metrics                   *Metrics                          `json:"metrics,omitempty"`
	NodeSelector              map[string]string                 `json:"nodeSelector,omitempty"`
	PodPriorityClassName      string                            `json:"podPriorityClassName,omitempty"`
	Port                      int32                             `json:"port,omitempty"`
	ReadinessProbe            *corev1.Probe                     `json:"readinessProbe,omitempty"`
	ReadinessDefaultCheck     ReadinessDefaultCheck             `json:"readinessDefaultCheck,omitempty"`
	Resources                 corev1.ResourceRequirements       `json:"resources,omitempty"`
	RootDir                   string                            `json:"rootDir,omitempty"`
	Scaling                   *FluentdScaling                   `json:"scaling,omitempty"`
	Security                  *Security                         `json:"security,omitempty"`
	ServiceAccountOverrides   *typeoverride.ServiceAccount      `json:"serviceAccount,omitempty"`
	StatefulSetAnnotations    map[string]string                 `json:"statefulsetAnnotations,omitempty"`
	TLS                       FluentdTLS                        `json:"tls,omitempty"`
	Tolerations               []corev1.Toleration               `json:"tolerations,omitempty"`
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	VolumeMountChmod          bool                              `json:"volumeMountChmod,omitempty"`
	VolumeModImage            ImageSpec                         `json:"volumeModImage,omitempty"`
	Workers                   int32                             `json:"workers,omitempty"`
}

// +kubebuilder:object:generate=true

type FluentOutLogrotate struct {
	Age     string `json:"age,omitempty"`
	Enabled bool   `json:"enabled"`
	Path    string `json:"path,omitempty"`
	Size    string `json:"size,omitempty"`
}

// +kubebuilder:object:generate=true

// ExtraVolume defines the fluentd extra volumes
type ExtraVolume struct {
	ContainerName string                   `json:"containerName,omitempty"`
	Path          string                   `json:"path,omitempty"`
	Volume        *volume.KubernetesVolume `json:"volume,omitempty"`
	VolumeName    string                   `json:"volumeName,omitempty"`
}

func (e *ExtraVolume) GetVolume() (corev1.Volume, error) {
	return e.Volume.GetVolume(e.VolumeName)
}

func (e *ExtraVolume) ApplyVolumeForPodSpec(spec *corev1.PodSpec) error {
	return e.Volume.ApplyVolumeForPodSpec(e.VolumeName, e.ContainerName, e.Path, spec)
}

// +kubebuilder:object:generate=true

// FluentdScaling enables configuring the scaling behaviour of the fluentd statefulset
type FluentdScaling struct {
	Drain               FluentdDrainConfig `json:"drain,omitempty"`
	PodManagementPolicy string             `json:"podManagementPolicy,omitempty"`
	Replicas            int                `json:"replicas,omitempty"`
}

// +kubebuilder:object:generate=true

// FluentdTLS defines the TLS configs
type FluentdTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName,omitempty"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

// +kubebuilder:object:generate=true

// FluentdDrainConfig enables configuring the drain behavior when scaling down the fluentd statefulset
type FluentdDrainConfig struct {
	// Container image to use for the drain watch sidecar
	Annotations map[string]string `json:"annotations,omitempty"`
	// Should buffers on persistent volumes left after scaling down the statefulset be drained
	Enabled bool `json:"enabled,omitempty"`
	// Should persistent volume claims be deleted after draining is done
	DeleteVolume bool      `json:"deleteVolume,omitempty"`
	Image        ImageSpec `json:"image,omitempty"`
	// Container image to use for the fluentd placeholder pod
	PauseImage ImageSpec `json:"pauseImage,omitempty"`
}
