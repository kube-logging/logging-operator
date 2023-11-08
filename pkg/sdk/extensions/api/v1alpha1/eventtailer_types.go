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
	"github.com/cisco-open/operator-tools/pkg/volume"
	"github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/tailer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +name:"EventTailer"
// +weight:"200"
type _hugoEventTailer = interface{} //nolint:deadcode,unused

// +name:"EventTailer"
// +version:"v1alpha1"
// +description:"Eventtailer's main goal is to listen kubernetes events and transmit their changes to stdout. This way the logging-operator is able to process them."
type _metaEventTailer = interface{} //nolint:deadcode,unused

// EventTailerSpec defines the desired state of EventTailer
type EventTailerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:Required
	// The resources of EventTailer will be placed into this namespace
	ControlNamespace string `json:"controlNamespace"`
	// Volume definition for tracking fluentbit file positions (optional)
	PositionVolume volume.KubernetesVolume `json:"positionVolume,omitempty"`
	// Override metadata of the created resources
	WorkloadMetaBase *types.MetaBase `json:"workloadMetaOverrides,omitempty"`
	// Override podSpec fields for the given statefulset
	WorkloadBase *types.PodSpecBase `json:"workloadOverrides,omitempty"`
	// Override container fields for the given statefulset
	ContainerBase *types.ContainerBase `json:"containerOverrides,omitempty"`
	// Override image related fields for the given statefulset, highest precedence
	Image *tailer.ImageSpec `json:"image,omitempty"`
}

// EventTailerStatus defines the observed state of EventTailer
type EventTailerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=eventtailers,scope=Cluster

// EventTailer is the Schema for the eventtailers API
type EventTailer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventTailerSpec   `json:"spec,omitempty"`
	Status EventTailerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EventTailerList contains a list of EventTailer
type EventTailerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventTailer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventTailer{}, &EventTailerList{})
}
