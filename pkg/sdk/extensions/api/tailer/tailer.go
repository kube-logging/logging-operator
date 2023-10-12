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

package tailer

import (
	"github.com/cisco-open/operator-tools/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// Tailer .
type Tailer interface {
	Command(Name string) []string
	GeneralDescriptor() General
}

// +kubebuilder:object:generate=true

// General descriptor for hosttailers
type General struct {
	Name          string               `json:"name"`
	Path          string               `json:"path,omitempty"`
	Disabled      bool                 `json:"disabled,omitempty"`
	ContainerBase *types.ContainerBase `json:"containerOverrides,omitempty"`
	Image         *ImageSpec           `json:"image,omitempty"`
}

// +kubebuilder:object:generate=true
type ImageSpec struct {
	Repository       string                        `json:"repository,omitempty"`
	Tag              string                        `json:"tag,omitempty"`
	PullPolicy       string                        `json:"pullPolicy,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

func (s ImageSpec) RepositoryWithTag() string {
	return RepositoryWithTag(s.Repository, s.Tag)
}

func RepositoryWithTag(repository, tag string) string {
	res := repository
	if tag != "" {
		res += ":" + tag
	}
	return res
}
