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

package v1alpha1

import (
	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/tailer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +name:"HostTailer"
// +weight:"200"
type _hugoHostTailer = interface{} //nolint:deadcode,unused

// +name:"HostTailer"
// +version:"v1alpha1"
// +description:"HostTailer's main goal is to tail custom files and transmit their changes to stdout. This way the logging-operator is able to process them."
type _metaHostTailer = interface{} //nolint:deadcode,unused

// HostTailerSpec defines the desired state of HostTailer
type HostTailerSpec struct {
	// List of [file tailers](#filetailer).
	FileTailers []FileTailer `json:"fileTailers,omitempty"`
	// List of [systemd tailers](#systemdtailer).
	SystemdTailers []SystemdTailer `json:"systemdTailers,omitempty"`
	// EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the
	// daemonset (and possibly other resource in the future) in case there is a change in an immutable field
	// that otherwise couldn't be managed with a simple update.
	EnableRecreateWorkloadOnImmutableFieldChange bool `json:"enableRecreateWorkloadOnImmutableFieldChange,omitempty"`
	//+kubebuilder:validation:Required
	// Override metadata of the created resources
	WorkloadMetaBase *types.MetaBase `json:"workloadMetaOverrides,omitempty"`
	// Override podSpec fields for the given daemonset
	WorkloadBase *types.PodSpecBase `json:"workloadOverrides,omitempty"`
	Image        tailer.ImageSpec   `json:"image,omitempty"`
}

// HostTailerStatus defines the observed state of [HostTailer](#hosttailer).
type HostTailerStatus struct {
}

// +kubebuilder:object:root=true

// HostTailer is the Schema for the hosttailers API
type HostTailer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HostTailerSpec   `json:"spec,omitempty"`
	Status HostTailerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HostTailerList contains a list of [HostTailers](#hosttailer).
type HostTailerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HostTailer `json:"items"`
}

// FileTailer configuration options
type FileTailer struct {
	// Name for the tailer
	Name string `json:"name"`
	// Path to the loggable file
	Path string `json:"path,omitempty"`
	// Disable tailing the file
	Disabled bool `json:"disabled,omitempty"`
	// Set the limit of the buffer size per active filetailer
	BufferMaxSize string `json:"buffer_max_size,omitempty"`
	// Set the buffer chunk size per active filetailer
	BufferChunkSize string `json:"buffer_chunk_size,omitempty"`
	// Skip long line when exceeding Buffer_Max_Size
	SkipLongLines string `json:"skip_long_lines,omitempty"`
	// Start reading from the head of new log files
	ReadFromHead bool `json:"read_from_head,omitempty"`
	// Override container fields for the given tailer
	ContainerBase *types.ContainerBase `json:"containerOverrides,omitempty"`
	// Override image field for the given trailer
	Image *tailer.ImageSpec `json:"image,omitempty"`
}

// SystemdTailer configuration options
type SystemdTailer struct {
	// Name for the tailer
	Name string `json:"name"`
	// Override systemd log path
	Path string `json:"path,omitempty"`
	// Disable component
	Disabled bool `json:"disabled,omitempty"`
	// Filter to select systemd unit example: kubelet.service
	SystemdFilter string `json:"systemdFilter,omitempty"`
	// Maximum entries to read when starting to tail logs to avoid high pressure
	MaxEntries int `json:"maxEntries,omitempty"`
	// Override container fields for the given tailer
	ContainerBase *types.ContainerBase `json:"containerOverrides,omitempty"`
	// Override image field for the given trailer
	Image *tailer.ImageSpec `json:"image,omitempty"`
}

func init() {
	SchemeBuilder.Register(&HostTailer{}, &HostTailerList{})
}
