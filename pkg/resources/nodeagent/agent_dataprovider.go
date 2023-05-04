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
	"fmt"

	util "github.com/cisco-open/operator-tools/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type AgentDataProvider interface {
	TargetHost() string
	ConfigFileName() string
	QualifiedName(string) string
	Namespace() string
	ResourceLabels() map[string]string
	ResourceAnnotations() map[string]string
	OwnerRefs() []v1.OwnerReference
}

type GenericDataProvider struct {
	Logging               v1beta1.Logging
	AggregatorServiceName string
	AgentType             string
	AgentConfigFileName   string
	AgentObject           client.Object
	AgentTypeMeta         v1.TypeMeta
}

func (f *GenericDataProvider) OwnerRefs() []v1.OwnerReference {
	return []v1.OwnerReference{
		{
			APIVersion: f.AgentTypeMeta.APIVersion,
			Kind:       f.AgentTypeMeta.Kind,
			Name:       f.AgentObject.GetName(),
			UID:        f.AgentObject.GetUID(),
			Controller: util.BoolPointer(true),
		},
	}
}

func (f *GenericDataProvider) ResourceLabels() map[string]string {
	return util.MergeLabels(f.AgentObject.GetLabels(), map[string]string{
		"app.kubernetes.io/name":     f.AgentType,
		"app.kubernetes.io/instance": f.AgentObject.GetName(),
	}, generateLoggingRefLabels(f.Logging.GetName()))
}

func (f *GenericDataProvider) ResourceAnnotations() map[string]string {
	return f.AgentObject.GetAnnotations()
}

func (f *GenericDataProvider) Namespace() string {
	return f.Logging.Spec.ControlNamespace
}

func (f *GenericDataProvider) QualifiedName(s string) string {
	return fmt.Sprintf("%s-%s-%s", f.AgentObject.GetName(), f.AgentType, s)
}

func (f *GenericDataProvider) ConfigFileName() string {
	return f.AgentConfigFileName
}

func (f *GenericDataProvider) TargetHost() string {
	return fmt.Sprintf("%s.%s.svc%s", f.Logging.QualifiedName(f.AggregatorServiceName), f.Logging.Spec.ControlNamespace, f.Logging.ClusterDomainAsSuffix())
}
