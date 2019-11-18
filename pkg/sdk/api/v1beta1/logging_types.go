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
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LoggingSpec defines the desired state of Logging
type LoggingSpec struct {
	LoggingRef              string         `json:"loggingRef,omitempty"`
	FlowConfigCheckDisabled bool           `json:"flowConfigCheckDisabled,omitempty"`
	FlowConfigOverride      string         `json:"flowConfigOverride,omitempty"`
	FluentbitSpec           *FluentbitSpec `json:"fluentbit,omitempty"`
	FluentdSpec             *FluentdSpec   `json:"fluentd,omitempty"`
	WatchNamespaces         []string       `json:"watchNamespaces,omitempty"`
	ControlNamespace        string         `json:"controlNamespace"`
}

// LoggingStatus defines the observed state of Logging
type LoggingStatus struct {
	ConfigCheckResults map[string]bool `json:"configCheckResults,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=loggings,scope=Cluster

// Logging is the Schema for the loggings API
type Logging struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoggingSpec   `json:"spec,omitempty"`
	Status LoggingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoggingList contains a list of Logging
type LoggingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Logging `json:"items"`
}

// SetDefaults fill empty attributes
func (l *Logging) SetDefaults() (*Logging, error) {
	copy := l.DeepCopy()
	if !copy.Spec.FlowConfigCheckDisabled && copy.Status.ConfigCheckResults == nil {
		copy.Status.ConfigCheckResults = make(map[string]bool)
	}
	if copy.Spec.WatchNamespaces == nil {
		copy.Spec.WatchNamespaces = []string{}
	}
	if copy.Spec.FluentdSpec != nil {
		if copy.Spec.FluentdSpec.Image.Repository == "" {
			copy.Spec.FluentdSpec.Image.Repository = "banzaicloud/fluentd"
		}
		if copy.Spec.FluentdSpec.Image.Tag == "" {
			copy.Spec.FluentdSpec.Image.Tag = "v1.7.4-alpine-8"
		}
		if copy.Spec.FluentdSpec.Image.PullPolicy == "" {
			copy.Spec.FluentdSpec.Image.PullPolicy = "IfNotPresent"
		}
		if copy.Spec.FluentdSpec.Annotations == nil {
			copy.Spec.FluentdSpec.Annotations = make(map[string]string)
		}
		if copy.Spec.FluentdSpec.Security == nil {
			copy.Spec.FluentdSpec.Security = &Security{}
		}
		if copy.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate == nil {
			copy.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate = util.BoolPointer(true)
		}
		if copy.Spec.FluentdSpec.Security.SecurityContext == nil {
			copy.Spec.FluentdSpec.Security.SecurityContext = &v1.SecurityContext{}
		}
		if copy.Spec.FluentdSpec.Security.PodSecurityContext == nil {
			copy.Spec.FluentdSpec.Security.PodSecurityContext = &v1.PodSecurityContext{}
		}
		if copy.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup == nil {
			copy.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup = util.IntPointer64(101)
		}
		if copy.Spec.FluentdSpec.Metrics != nil {
			if copy.Spec.FluentdSpec.Metrics.Path == "" {
				copy.Spec.FluentdSpec.Metrics.Path = "/metrics"
			}
			if copy.Spec.FluentdSpec.Metrics.Port == 0 {
				copy.Spec.FluentdSpec.Metrics.Port = 24231
			}
			if copy.Spec.FluentdSpec.Metrics.Timeout == "" {
				copy.Spec.FluentdSpec.Metrics.Timeout = "5s"
			}
			if copy.Spec.FluentdSpec.Metrics.Interval == "" {
				copy.Spec.FluentdSpec.Metrics.Interval = "15s"
			}

			if copy.Spec.FluentdSpec.Metrics.PrometheusAnnotations {
				copy.Spec.FluentdSpec.Annotations["prometheus.io/scrape"] = "true"

				copy.Spec.FluentdSpec.Annotations["prometheus.io/path"] = copy.Spec.FluentdSpec.Metrics.Path
				copy.Spec.FluentdSpec.Annotations["prometheus.io/port"] = string(copy.Spec.FluentdSpec.Metrics.Port)
			}
		}
		if copy.Spec.FluentdSpec.FluentdPvcSpec.AccessModes == nil {
			copy.Spec.FluentdSpec.FluentdPvcSpec.AccessModes = []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			}
		}
		if copy.Spec.FluentdSpec.FluentdPvcSpec.VolumeMode == nil {
			copy.Spec.FluentdSpec.FluentdPvcSpec.VolumeMode = persistentVolumeModePointer(v1.PersistentVolumeFilesystem)
		}
		if copy.Spec.FluentdSpec.FluentdPvcSpec.Resources.Requests == nil {
			copy.Spec.FluentdSpec.FluentdPvcSpec.Resources.Requests = map[v1.ResourceName]resource.Quantity{
				"storage": resource.MustParse("20Gi"),
			}
		}
		// Temporarily copy the FluentdPvcSpec for backward compatibility
		// if BufferStorageVolume.PersistentVolumeClaim is not set.
		// DisablePvc will stay for a while. The alternative would be to set a hostPath or emptyDir explicitly
		if copy.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim == nil {
			copy.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim = &PersistentVolumeClaim{
				PersistentVolumeClaimSpec: copy.Spec.FluentdSpec.FluentdPvcSpec,
				PersistentVolumeSource: v1.PersistentVolumeClaimVolumeSource{
					ClaimName: l.QualifiedName("fluentd-buffer"),
					ReadOnly:  false,
				},
			}
		}
		if copy.Spec.FluentdSpec.VolumeModImage.Repository == "" {
			copy.Spec.FluentdSpec.VolumeModImage.Repository = "busybox"
		}
		if copy.Spec.FluentdSpec.VolumeModImage.Tag == "" {
			copy.Spec.FluentdSpec.VolumeModImage.Tag = "latest"
		}
		if copy.Spec.FluentdSpec.VolumeModImage.PullPolicy == "" {
			copy.Spec.FluentdSpec.VolumeModImage.PullPolicy = "IfNotPresent"
		}
		if copy.Spec.FluentdSpec.ConfigReloaderImage.Repository == "" {
			copy.Spec.FluentdSpec.ConfigReloaderImage.Repository = "jimmidyson/configmap-reload"
		}
		if copy.Spec.FluentdSpec.ConfigReloaderImage.Tag == "" {
			copy.Spec.FluentdSpec.ConfigReloaderImage.Tag = "v0.2.2"
		}
		if copy.Spec.FluentdSpec.ConfigReloaderImage.PullPolicy == "" {
			copy.Spec.FluentdSpec.ConfigReloaderImage.PullPolicy = "IfNotPresent"
		}
		if copy.Spec.FluentdSpec.Resources.Limits == nil {
			copy.Spec.FluentdSpec.Resources.Limits = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("200M"),
				v1.ResourceCPU:    resource.MustParse("1000m"),
			}
		}
		if copy.Spec.FluentdSpec.Resources.Requests == nil {
			copy.Spec.FluentdSpec.Resources.Requests = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("100M"),
				v1.ResourceCPU:    resource.MustParse("500m"),
			}
		}
		if copy.Spec.FluentdSpec.Port == 0 {
			copy.Spec.FluentdSpec.Port = 24240
		}
		if copy.Spec.FluentdSpec.Scaling == nil {
			copy.Spec.FluentdSpec.Scaling = &FluentdScaling{
				Replicas: 1,
			}
		}
	}
	if copy.Spec.FluentbitSpec != nil {
		if copy.Spec.FluentbitSpec.Image.Repository == "" {
			copy.Spec.FluentbitSpec.Image.Repository = "fluent/fluent-bit"
		}
		if copy.Spec.FluentbitSpec.Image.Tag == "" {
			copy.Spec.FluentbitSpec.Image.Tag = "1.3.2"
		}
		if copy.Spec.FluentbitSpec.Image.PullPolicy == "" {
			copy.Spec.FluentbitSpec.Image.PullPolicy = "IfNotPresent"
		}
		if copy.Spec.FluentbitSpec.Resources.Limits == nil {
			copy.Spec.FluentbitSpec.Resources.Limits = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("100M"),
				v1.ResourceCPU:    resource.MustParse("200m"),
			}
		}
		if copy.Spec.FluentbitSpec.Resources.Requests == nil {
			copy.Spec.FluentbitSpec.Resources.Requests = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("50M"),
				v1.ResourceCPU:    resource.MustParse("100m"),
			}
		}
		// For backward compatibility
		if copy.Spec.FluentbitSpec.Parser != "" {
			if copy.Spec.FluentbitSpec.InputTail.Parser == "" {
				copy.Spec.FluentbitSpec.InputTail.Parser = copy.Spec.FluentbitSpec.Parser
			}
		}
		if copy.Spec.FluentbitSpec.InputTail.Path == "" {
			copy.Spec.FluentbitSpec.InputTail.Path = "/var/log/containers/*.log"
		}
		if copy.Spec.FluentbitSpec.InputTail.RefreshInterval == "" {
			copy.Spec.FluentbitSpec.InputTail.RefreshInterval = "5"
		}
		if copy.Spec.FluentbitSpec.InputTail.SkipLongLines == "" {
			copy.Spec.FluentbitSpec.InputTail.SkipLongLines = "On"
		}
		if copy.Spec.FluentbitSpec.InputTail.DB == nil {
			copy.Spec.FluentbitSpec.InputTail.DB = util.StringPointer("/tail-db/tail-containers-state.db")
		}
		if copy.Spec.FluentbitSpec.InputTail.MemBufLimit == "" {
			copy.Spec.FluentbitSpec.InputTail.MemBufLimit = "5MB"
		}
		if copy.Spec.FluentbitSpec.InputTail.Tag == "" {
			copy.Spec.FluentbitSpec.InputTail.Tag = "kubernetes.*"
		}
		if copy.Spec.FluentbitSpec.PositionDBLegacy != nil {
			copy.Spec.FluentbitSpec.PositionDB = *copy.Spec.FluentbitSpec.PositionDBLegacy.DeepCopy()
		}
		if copy.Spec.FluentbitSpec.Annotations == nil {
			copy.Spec.FluentbitSpec.Annotations = make(map[string]string)
		}
		if copy.Spec.FluentbitSpec.Security == nil {
			copy.Spec.FluentbitSpec.Security = &Security{}
		}
		if copy.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate == nil {
			copy.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate = util.BoolPointer(true)
		}
		if copy.Spec.FluentbitSpec.Security.SecurityContext == nil {
			copy.Spec.FluentbitSpec.Security.SecurityContext = &v1.SecurityContext{}
		}
		if copy.Spec.FluentbitSpec.Security.PodSecurityContext == nil {
			copy.Spec.FluentbitSpec.Security.PodSecurityContext = &v1.PodSecurityContext{}
		}
		if copy.Spec.FluentbitSpec.Metrics != nil {
			if copy.Spec.FluentbitSpec.Metrics.Path == "" {
				copy.Spec.FluentbitSpec.Metrics.Path = "/api/v1/metrics/prometheus"
			}
			if copy.Spec.FluentbitSpec.Metrics.Port == 0 {
				copy.Spec.FluentbitSpec.Metrics.Port = 2020
			}
			if copy.Spec.FluentbitSpec.Metrics.Timeout == "" {
				copy.Spec.FluentbitSpec.Metrics.Timeout = "5s"
			}
			if copy.Spec.FluentbitSpec.Metrics.Interval == "" {
				copy.Spec.FluentbitSpec.Metrics.Interval = "15s"
			}
			if copy.Spec.FluentbitSpec.Metrics.PrometheusAnnotations {
				copy.Spec.FluentbitSpec.Annotations["prometheus.io/scrape"] = "true"
				copy.Spec.FluentbitSpec.Annotations["prometheus.io/path"] = copy.Spec.FluentbitSpec.Metrics.Path
				copy.Spec.FluentbitSpec.Annotations["prometheus.io/port"] = string(copy.Spec.FluentbitSpec.Metrics.Port)
			}
		}
		if copy.Spec.FluentbitSpec.MountPath == "" {
			copy.Spec.FluentbitSpec.MountPath = "/var/lib/docker/containers"
		}
		if copy.Spec.FluentbitSpec.BufferStorage.StoragePath == "" {
			copy.Spec.FluentbitSpec.BufferStorage.StoragePath = "/buffers"
		}
	}
	return copy, nil
}

// QualifiedName is the "logging-resource" name combined
func (l *Logging) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", l.Name, name)
}

// QualifiedNamespacedName is the "namespace-logging-resource" name combined
func (l *Logging) QualifiedNamespacedName(name string) string {
	return fmt.Sprintf("%s-%s-%s", l.Spec.ControlNamespace, l.Name, name)
}

func init() {
	SchemeBuilder.Register(&Logging{}, &LoggingList{})
}

func persistentVolumeModePointer(mode v1.PersistentVolumeMode) *v1.PersistentVolumeMode {
	return &mode
}
