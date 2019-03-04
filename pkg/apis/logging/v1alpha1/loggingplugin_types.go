/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PluginSpec defines the desired state of Plugin
// +k8s:openapi-gen=true
type PluginSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Input  Input     `json:"input,omitempty"`
	Filter []FPlugin `json:"filter,omitempty"`
	Output []FPlugin `json:"output,omitempty"`
}

// Input this determines the log origin
type Input struct {
	Label map[string]string `json:"label"`
}

// FPlugin struct for fluentd plugins
type FPlugin struct {
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

// Parameter generic parameter type to handle values from different sources
type Parameter struct {
	Name      string     `json:"name"`
	ValueFrom *ValueFrom `json:"valueFrom,omitempty"`
	Value     string     `json:"value"`
}

// GetValue for a Parameter
func (p Parameter) GetValue(namespace string, client client.Client) (string, string) {
	if p.ValueFrom != nil {
		value, error := p.ValueFrom.GetValue(namespace, client)
		if error != nil {
			return "", ""
		}
		return p.Name, value
	}
	return p.Name, p.Value
}

// ValueFrom generic type to determine value origin
type ValueFrom struct {
	SecretKeyRef KubernetesSecret `json:"secretKeyRef"`
}

// GetValue handles the different origin of ValueFrom
func (vf *ValueFrom) GetValue(namespace string, client client.Client) (string, error) {
	return vf.SecretKeyRef.GetValue(namespace, client)
}

// KubernetesSecret is a ValueFrom type
type KubernetesSecret struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
}

// GetValue implement GetValue interface
func (ks KubernetesSecret) GetValue(namespace string, client client.Client) (string, error) {
	secret := &corev1.Secret{}
	nSpace := namespace
	if ks.Namespace != "" {
		nSpace = ks.Namespace
	}
	err := client.Get(context.TODO(), types.NamespacedName{Name: ks.Name, Namespace: nSpace}, secret)
	if err != nil {
		return "", err
	}
	value, ok := secret.Data[ks.Key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret %q in namespace %q", ks.Key, secret.ObjectMeta.Name, secret.ObjectMeta.Namespace)
	}
	return string(value), nil
}

// PluginStatus defines the observed state of Plugin
// +k8s:openapi-gen=true
type PluginStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Plugin is the Schema for the Plugin API
// +k8s:openapi-gen=true
type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PluginSpec   `json:"spec,omitempty"`
	Status PluginStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PluginList contains a list of Plugin
type PluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Plugin{}, &PluginList{})
}
