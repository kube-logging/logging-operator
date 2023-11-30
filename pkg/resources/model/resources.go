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
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type LoggingResources struct {
	AllLoggings     []v1beta1.Logging
	Logging         v1beta1.Logging
	Fluentd         FluentdLoggingResources
	SyslogNG        SyslogNGLoggingResources
	NodeAgents      []v1beta1.NodeAgent
	Fluentbits      []v1beta1.FluentbitAgent
	LoggingRoutes   []v1beta1.LoggingRoute
	WatchNamespaces []string
}

func (l LoggingResources) GetFluentd() *v1beta1.Fluentd {
	if l.Fluentd.Configuration != nil {
		return l.Fluentd.Configuration
	}
	return nil
}

func (l LoggingResources) GetFluentdSpec() *v1beta1.FluentdSpec {

	if detachedFluentd := l.GetFluentd(); detachedFluentd != nil {
		return &detachedFluentd.Spec
	}
	if l.Logging.Spec.FluentdSpec != nil {
		return l.Logging.Spec.FluentdSpec
	}

	return nil
}

type FluentdLoggingResources struct {
	ClusterFlows   []v1beta1.ClusterFlow
	ClusterOutputs ClusterOutputs
	Flows          []v1beta1.Flow
	Outputs        Outputs
	Configuration  *v1beta1.Fluentd
	ExcessFluentds []v1beta1.Fluentd
}

type SyslogNGLoggingResources struct {
	ClusterFlows   []v1beta1.SyslogNGClusterFlow
	ClusterOutputs SyslogNGClusterOutputs
	Flows          []v1beta1.SyslogNGFlow
	Outputs        SyslogNGOutputs
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

type SyslogNGClusterOutputs []v1beta1.SyslogNGClusterOutput

func (c SyslogNGClusterOutputs) FindByName(name string) *v1beta1.SyslogNGClusterOutput {
	for i := range c {
		output := &c[i]
		if output.Name == name {
			return output
		}
	}
	return nil
}

type SyslogNGOutputs []v1beta1.SyslogNGOutput

func (c SyslogNGOutputs) FindByNamespacedName(namespace string, name string) *v1beta1.SyslogNGOutput {
	for i := range c {
		output := &c[i]
		if output.Namespace == namespace && output.Name == name {
			return output
		}
	}
	return nil
}
