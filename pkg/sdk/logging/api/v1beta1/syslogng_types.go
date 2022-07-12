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
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/input"
	"github.com/banzaicloud/operator-tools/pkg/typeoverride"
	"github.com/banzaicloud/operator-tools/pkg/volume"
	corev1 "k8s.io/api/core/v1"
)

// +name:"SyslogNGSpec"
// +weight:"200"
type _hugoSyslogNGSpec interface{} //nolint:deadcode,unused

// +name:"SyslogNGSpec"
// +version:"v1beta1"
// +description:"SyslogNGSpec defines the desired state of SyslogNG"
type _metaSyslogNGSpec interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true

// SyslogNGSpec defines the desired state of SyslogNG
type SyslogNGSpec struct {
	StatefulSetAnnotations map[string]string `json:"statefulsetAnnotations,omitempty"`
	Annotations            map[string]string `json:"annotations,omitempty"`
	ConfigCheckAnnotations map[string]string `json:"configCheckAnnotations,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty"`
	EnvVars                []corev1.EnvVar   `json:"envVars,omitempty"`
	TLS                    SyslogNGTLS       `json:"tls,omitempty"`
	Image                  ImageSpec         `json:"image,omitempty"`
	DisablePvc             bool              `json:"disablePvc,omitempty"`
	// BufferStorageVolume is by default configured as PVC using SyslogNGPvcSpec
	// +docLink:"volume.KubernetesVolume,https://github.com/banzaicloud/operator-tools/tree/master/docs/types"
	BufferStorageVolume volume.KubernetesVolume `json:"bufferStorageVolume,omitempty"`
	ExtraVolumes        []SyslogNGExtraVolume   `json:"extraVolumes,omitempty"`
	// Deprecated, use bufferStorageVolume
	SyslogNGPvcSpec             *volume.KubernetesVolume          `json:"SyslogNGPvcSpec,omitempty"`
	VolumeMountChmod            bool                              `json:"volumeMountChmod,omitempty"`
	VolumeModImage              ImageSpec                         `json:"volumeModImage,omitempty"`
	ConfigReloaderImage         ImageSpec                         `json:"configReloaderImage,omitempty"`
	PrometheusExporterImage     ImageSpec                         `json:"prometheusExporterImage,omitempty"`
	PrometheusExporterResources corev1.ResourceRequirements       `json:"prometheusExporterResources,omitempty"`
	Resources                   corev1.ResourceRequirements       `json:"resources,omitempty"`
	ConfigCheckResources        corev1.ResourceRequirements       `json:"configCheckResources,omitempty"`
	LivenessProbe               *corev1.Probe                     `json:"livenessProbe,omitempty"`
	LivenessDefaultCheck        bool                              `json:"livenessDefaultCheck,omitempty"`
	ReadinessProbe              *corev1.Probe                     `json:"readinessProbe,omitempty"`
	ReadinessDefaultCheck       ReadinessDefaultCheck             `json:"readinessDefaultCheck,omitempty"`
	PortTCP                     *int32                            `json:"portTCP,omitempty"`
	PortUDP                     *int32                            `json:"portUDP,omitempty"`
	Tolerations                 []corev1.Toleration               `json:"tolerations,omitempty"`
	NodeSelector                map[string]string                 `json:"nodeSelector,omitempty"`
	Affinity                    *corev1.Affinity                  `json:"affinity,omitempty"`
	TopologySpreadConstraints   []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	Metrics                     *Metrics                          `json:"metrics,omitempty"`
	BufferVolumeMetrics         *Metrics                          `json:"bufferVolumeMetrics,omitempty"`
	BufferVolumeImage           ImageSpec                         `json:"bufferVolumeImage,omitempty"`
	BufferVolumeArgs            []string                          `json:"bufferVolumeArgs,omitempty"`
	Security                    *Security                         `json:"security,omitempty"`
	Scaling                     *SyslogNGScaling                  `json:"scaling,omitempty"`
	Workers                     int32                             `json:"workers,omitempty"`
	RootDir                     string                            `json:"rootDir,omitempty"`
	// +kubebuilder:validation:enum=fatal,error,warn,info,debug,trace
	LogLevel string `json:"logLevel,omitempty"`
	// Ignore same log lines
	// +docLink:"more info, https://docs.SyslogNG.org/deployment/logging#ignore_same_log_interval"
	IgnoreSameLogInterval string `json:"ignoreSameLogInterval,omitempty"`
	// Ignore repeated log lines
	// +docLink:"more info, https://docs.SyslogNG.org/deployment/logging#ignore_repeated_log_interval"
	IgnoreRepeatedLogInterval string `json:"ignoreRepeatedLogInterval,omitempty"`
	PodPriorityClassName      string `json:"podPriorityClassName,omitempty"`
	// +kubebuilder:validation:enum=stdout,null
	SyslogNGLogDestination string `json:"syslogNGLogDestination,omitempty"`
	// SyslogNGOutLogrotate sends syslog-ng's stdout to file and rotates it
	SyslogNGOutLogrotate    *SyslogNGOutLogrotate        `json:"syslogNGOutLogrotate,omitempty"`
	ForwardInputConfig      *input.ForwardInputConfig    `json:"forwardInputConfig,omitempty"`
	ServiceAccountOverrides *typeoverride.ServiceAccount `json:"serviceAccount,omitempty"`
	DNSPolicy               corev1.DNSPolicy             `json:"dnsPolicy,omitempty"`
	DNSConfig               *corev1.PodDNSConfig         `json:"dnsConfig,omitempty"`
}

// +kubebuilder:object:generate=true

type SyslogNGOutLogrotate struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path,omitempty"`
	Age     string `json:"age,omitempty"`
	Size    string `json:"size,omitempty"`
}

// +kubebuilder:object:generate=true

// ExtraVolume defines the SyslogNG extra volumes
type SyslogNGExtraVolume struct {
	VolumeName    string                   `json:"volumeName,omitempty"`
	Path          string                   `json:"path,omitempty"`
	ContainerName string                   `json:"containerName,omitempty"`
	Volume        *volume.KubernetesVolume `json:"volume,omitempty"`
}

func (e *SyslogNGExtraVolume) GetVolume() (corev1.Volume, error) {
	return e.Volume.GetVolume(e.VolumeName)
}

func (e *SyslogNGExtraVolume) ApplyVolumeForPodSpec(spec *corev1.PodSpec) error {
	return e.Volume.ApplyVolumeForPodSpec(e.VolumeName, e.ContainerName, e.Path, spec)
}

// +kubebuilder:object:generate=true

// SyslogNGScaling enables configuring the scaling behaviour of the SyslogNG statefulset
type SyslogNGScaling struct {
	Replicas            int                 `json:"replicas,omitempty"`
	PodManagementPolicy string              `json:"podManagementPolicy,omitempty"`
	Drain               SyslogNGDrainConfig `json:"drain,omitempty"`
}

// +kubebuilder:object:generate=true

// SyslogNGTLS defines the TLS configs
type SyslogNGTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName,omitempty"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

// +kubebuilder:object:generate=true

// SyslogNGDrainConfig enables configuring the drain behavior when scaling down the SyslogNG statefulset
type SyslogNGDrainConfig struct {
	// Should buffers on persistent volumes left after scaling down the statefulset be drained
	Enabled bool `json:"enabled,omitempty"`
	// Container image to use for the drain watch sidecar
	Annotations map[string]string `json:"annotations,omitempty"`
	Image       ImageSpec         `json:"image,omitempty"`
	// Container image to use for the SyslogNG placeholder pod
	PauseImage ImageSpec `json:"pauseImage,omitempty"`
}
