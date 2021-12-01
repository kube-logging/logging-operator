// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package v1alpha1

import (
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +name:"Promtail"
// +weight:"200"
type _hugoPromtail = interface{}

// +name:"Promtail"
// +version:"v1alpha1"
// +description:"Promtail's main goal is to tail custom files and transmit their changes to stdout. This way the logging-operator is able to process them."
type _metaPromtail = interface{}

// PromtailSpec defines the desired state of Promtail
type PromtailSpec struct {
	//+kubebuilder:validation:Required
	// The resources of Promtail will be placed into this namespace
	Namespace string `json:"namespace"`
	// EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the
	// daemonset (and possibly other resource in the future) in case there is a change in an immutable field
	// that otherwise couldn't be managed with a simple update.
	EnableRecreateWorkloadOnImmutableFieldChange bool `json:"enableRecreateWorkloadOnImmutableFieldChange,omitempty"`
	//+kubebuilder:validation:Required
	// Override metadata of the created resources
	WorkloadMetaBase *types.MetaBase `json:"workloadMetaOverrides,omitempty"`
	// Override podSpec fields for the given daemonset
	WorkloadBase *types.PodSpecBase `json:"workloadOverrides,omitempty"`
	// Override container fields for the given statefulset
	ContainerBase *types.ContainerBase `json:"containerOverrides,omitempty"`
	// Container Runtime  (docker, containerd)
	ContainerRuntime string `json:"containerRuntime,omitempty"`
	// PipelineStages  (docker, cri)
	PipelineStages []string `json:"pipelineStages,omitempty"`
	// Loki URL http://loki:3100/loki/api/v1/push
	LokiUrl string `json:"lokiUrl,omitempty"`
	// Security defines promtail deployment security properties
	Security             *Security                   `json:"security,omitempty"`
	Tolerations          []corev1.Toleration         `json:"tolerations,omitempty"`
	NodeSelector         map[string]string           `json:"nodeSelector,omitempty"`
	Affinity             *corev1.Affinity            `json:"affinity,omitempty"`
	PodPriorityClassName string                      `json:"podPriorityClassName,omitempty"`
	Resources            corev1.ResourceRequirements `json:"resources,omitempty"`
	Image                ImageSpec                   `json:"image,omitempty"`
}

// PromtailStatus defines the observed state of Promtail
type PromtailStatus struct {
}

// +kubebuilder:object:root=true

// Promtail is the Schema for the promtails API
type Promtail struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PromtailSpec   `json:"spec,omitempty"`
	Status PromtailStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PromtailList contains a list of Promtail
type PromtailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Promtail `json:"items"`
}

// Security defines promtail deployment security properties
type Security struct {
	ServiceAccount               string                     `json:"serviceAccount,omitempty"`
	RoleBasedAccessControlCreate *bool                      `json:"roleBasedAccessControlCreate,omitempty"`
	PodSecurityPolicyCreate      bool                       `json:"podSecurityPolicyCreate,omitempty"`
	SecurityContext              *corev1.SecurityContext    `json:"securityContext,omitempty"`
	PodSecurityContext           *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
}

// ImageSpec struct hold information about image specification
type ImageSpec struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
	PullPolicy string `json:"pullPolicy,omitempty"`
}

// SetDefaults fill empty attributes
func (p *Promtail) SetDefaults() {
	if p.Spec.Namespace == "" {
		p.Spec.Namespace = "default"
	}
	if p.Spec.Image.Repository == "" {
		p.Spec.Image.Repository = "grafana/promtail"
	}
	if p.Spec.Image.Tag == "" {
		p.Spec.Image.Tag = "1.6.0"
	}
	if p.Spec.Image.PullPolicy == "" {
		p.Spec.Image.PullPolicy = "IfNotPresent"
	}

	if p.Spec.Security == nil {
		p.Spec.Security = &Security{}
	}
	if p.Spec.Security.RoleBasedAccessControlCreate == nil {
		p.Spec.Security.RoleBasedAccessControlCreate = utils.BoolPointer(true)
	}
	if p.Spec.Security.SecurityContext == nil {
		p.Spec.Security.SecurityContext = &corev1.SecurityContext{}
	}
	if p.Spec.Security.PodSecurityContext == nil {
		p.Spec.Security.PodSecurityContext = &corev1.PodSecurityContext{}
	}
	if p.Spec.Security.PodSecurityContext.FSGroup == nil {
		p.Spec.Security.PodSecurityContext.FSGroup = utils.IntPointer64(101)
	}
	if p.Spec.LokiUrl == "" {
		p.Spec.LokiUrl = "http://loki:3100/loki/api/v1/push"
	}
	if p.Spec.Resources.Limits == nil {
		p.Spec.Resources.Limits = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("50Mi"),
			corev1.ResourceCPU:    resource.MustParse("250m"),
		}
	}
	if p.Spec.Resources.Requests == nil {
		p.Spec.Resources.Requests = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("25Mi"),
			corev1.ResourceCPU:    resource.MustParse("10m"),
		}
	}
}

func init() {
	SchemeBuilder.Register(&Promtail{}, &PromtailList{})
}
