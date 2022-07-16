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

package model

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type LoggingResources struct {
	Logging        v1beta1.Logging
	Outputs        Outputs
	Flows          []v1beta1.Flow
	ClusterOutputs ClusterOutputs
	ClusterFlows   []v1beta1.ClusterFlow
}

type ClusterOutputs []v1beta1.ClusterOutput

func (c ClusterOutputs) FindByName(name string) *v1beta1.ClusterOutput {
	for i := range c {
		output := &c[i]
		if output.Name == name {
			return output
		}
	}
	return nil
}

type Outputs []v1beta1.Output

func (c Outputs) FindByNamespacedName(namespace string, name string) *v1beta1.Output {
	for i := range c {
		output := &c[i]
		if output.Namespace == namespace && output.Name == name {
			return output
		}
	}
	return nil
}

type SyslogNGLoggingResources struct {
	Logging        v1beta1.Logging
	Flows          []v1beta1.SyslogNGFlow
	Outputs        []v1beta1.SyslogNGOutput
	ClusterFlows   []v1beta1.SyslogNGClusterFlow
	ClusterOutputs []v1beta1.SyslogNGClusterOutput
}
