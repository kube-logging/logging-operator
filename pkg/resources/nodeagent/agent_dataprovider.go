// Copyright © 2023 Kube logging authors
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

package nodeagent

import (
	"context"
	"fmt"

	util "github.com/cisco-open/operator-tools/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/resources/loggingdataprovider"
)

type AgentDataProvider interface {
	loggingdataprovider.LoggingDataProvider
	TargetHost() string
	ConfigFileName() string
	QualifiedName(string) string
	Namespace() string
	ResourceLabels() map[string]string
	ResourceAnnotations() map[string]string
	OwnerRefs() []v1.OwnerReference
	GetConstants() Constants
}

type Constants struct {
	LoggingName      string
	ControlNamespace string
	Name             string
	ContainerName    string
	ConfigFileName   string
	Kind             string
	APIVersion       string
	VolumeName       string
	StoragePath      string
	ConfigPath       string
}

type GenericDataProvider struct {
	LoggingDataProvider loggingdataprovider.LoggingDataProvider
	AgentObject         client.Object
	Constants           Constants
}

func (f *GenericDataProvider) TargetHost() string {
	return f.LoggingDataProvider.TargetHost()
}

func (f *GenericDataProvider) GetReplicaCount(ctx context.Context) (*int32, error) {
	return f.LoggingDataProvider.GetReplicaCount(ctx)
}

func (f *GenericDataProvider) GetConstants() Constants {
	return f.Constants
}

func (f *GenericDataProvider) OwnerRefs() []v1.OwnerReference {
	return []v1.OwnerReference{
		{
			APIVersion: f.Constants.APIVersion,
			Kind:       f.Constants.Kind,
			Name:       f.AgentObject.GetName(),
			UID:        f.AgentObject.GetUID(),
			Controller: util.BoolPointer(true),
		},
	}
}

func (f *GenericDataProvider) ResourceLabels() map[string]string {
	return util.MergeLabels(f.AgentObject.GetLabels(), map[string]string{
		"app.kubernetes.io/name":     f.Constants.Name,
		"app.kubernetes.io/instance": f.AgentObject.GetName(),
	}, generateLoggingRefLabels(f.Constants.LoggingName))
}

func (f *GenericDataProvider) ResourceAnnotations() map[string]string {
	return f.AgentObject.GetAnnotations()
}

func (f *GenericDataProvider) Namespace() string {
	return f.Constants.ControlNamespace
}

func (f *GenericDataProvider) QualifiedName(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf("%s-%s-%s", f.AgentObject.GetName(), f.Constants.Name, s)
	}
	return fmt.Sprintf("%s-%s", f.AgentObject.GetName(), f.Constants.Name)
}

func (f *GenericDataProvider) ConfigFileName() string {
	return f.Constants.ConfigFileName
}
