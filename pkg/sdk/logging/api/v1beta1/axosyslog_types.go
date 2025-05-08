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

package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=axosyslogs,scope=Namespaced,categories=logging-all
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Problems",type="integer",JSONPath=".status.problemsCount",description="Number of problems"

// AxoSyslog is the Schema for the AxoSyslogs API
type AxoSyslog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AxoSyslogSpec   `json:"spec,omitempty"`
	Status AxoSyslogStatus `json:"status,omitempty"`
}

// AxoSyslogSpec defines the desired state of AxoSyslog
type AxoSyslogSpec struct {
	// LogPaths is a list of log paths to be rendered in the AxoSyslog configuration
	LogPaths []LogPath `json:"logPaths,omitempty"`

	// Image is the image specification for AxoSyslog
	Image *BasicImageSpec `json:"image,omitempty"`

	// ConfigReloadImage is the image specification for the config reload
	ConfigReloadImage *BasicImageSpec `json:"configReloadImage,omitempty"`

	// Destinations is a list of destinations to be rendered in the AxoSyslog configuration
	Destinations []Destination `json:"destinations,omitempty"`
}

// LogPath defines a single log path that will be rendered in the AxoSyslog configuration
type LogPath struct {
	// filterx block to be rendered within the log path
	Filterx string `json:"filterx,omitempty"`
	// name of a destination to be used in the log path
	Destination string `json:"destination,omitempty"`
}

// Destination defines a single destination that will be rendered in the AxoSyslog configuration
type Destination struct {
	// Name of the destination
	Name string `json:"name,omitempty"`
	// Config is the configuration for the destination
	Config string `json:"config,omitempty"`
}

// AxoSyslogStatus defines the observed state of AxoSyslog
type AxoSyslogStatus struct {
	// Sources configured for AxoSyslog
	Sources []Source `json:"sources,omitempty"`
	// Problems with the AxoSyslog resource
	Problems []string `json:"problems,omitempty"`
	// Count of problems with the AxoSyslog resource
	ProblemsCount int `json:"problemsCount,omitempty"`
}

// Source represents the source of logs for AxoSyslog
type Source struct {
	// OTLP specific configuration
	OTLP *OTLPSource `json:"otlp,omitempty"`
}

// OTLPSource contains configuration for OpenTelemetry Protocol sources
type OTLPSource struct {
	// Endpoint for the OTLP source
	Endpoint string `json:"endpoint,omitempty"`
}

// +kubebuilder:object:root=true

// AxoSyslogList contains a list of AxoSyslog
type AxoSyslogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AxoSyslog `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AxoSyslog{}, &AxoSyslogList{})
}

func (a *AxoSyslog) SetDefaults() error {
	if a != nil {
		if a.Spec.LogPaths == nil {
			a.Spec.LogPaths = []LogPath{}
		}
		if a.Spec.Destinations == nil {
			a.Spec.Destinations = []Destination{}
		}

		if a.Spec.Image == nil {
			a.Spec.Image = &BasicImageSpec{
				Repository: defaultAxoSyslogImageRepository,
				Tag:        defaultAxoSyslogImageTag,
			}
		}
		if a.Spec.ConfigReloadImage == nil {
			a.Spec.ConfigReloadImage = &BasicImageSpec{}
		}
		if a.Spec.ConfigReloadImage.Repository == "" {
			a.Spec.ConfigReloadImage.Repository = defaultConfigReloaderImageRepository
		}
		if a.Spec.ConfigReloadImage.Tag == "" {
			if Version == "" {
				a.Spec.ConfigReloadImage.Tag = defaultConfigReloaderImageTag
			} else {
				a.Spec.ConfigReloadImage.Tag = Version
			}
		}
	}

	return nil
}
