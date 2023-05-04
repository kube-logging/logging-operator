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

package syslogng_agent

import (
	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources/nodeagent"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type SyslogNGAgentReconciler struct {
	agent           v1beta1.SyslogNGAgent
	agentReconciler *nodeagent.GenericAgentReconciler
	dataProvider    nodeagent.AgentDataProvider
	log             logr.Logger
}

func NewSyslogNGAgentReconciler(
	reconciler *nodeagent.GenericAgentReconciler,
	dataProvider nodeagent.AgentDataProvider,
	log logr.Logger,
	agent v1beta1.SyslogNGAgent,
) *SyslogNGAgentReconciler {
	return &SyslogNGAgentReconciler{
		agentReconciler: reconciler,
		agent:           agent,
		dataProvider:    dataProvider,
		log:             log,
	}
}

func (s *SyslogNGAgentReconciler) Reconcile() (*reconcile.Result, error) {
	metricsEnabled := s.agent.Spec.Metrics != nil
	prometheusAnnotationsEnabled := metricsEnabled && s.agent.Spec.Metrics.PrometheusAnnotations
	if spec, err := nodeagent.NodeAgentSyslogNGDefaults(metricsEnabled, prometheusAnnotationsEnabled); err != nil {
		s.agent.Spec = *spec
		return nil, errors.Wrap(err, "applying syslogNG defaults")
	} else {
		err = merge.Merge(&s.agent.Spec, spec)
		if err != nil {
			return nil, err
		}
	}

	s.log.Info("SyslogNGAgent", "syslogng-agent", s.agent)

	resourceBuilders := []reconciler.ResourceBuilder{
		s.configSecret,
	}

	result, err := s.agentReconciler.Reconcile(resourceBuilders)
	return &result, err
}

func (s *SyslogNGAgentReconciler) configSecret() (runtime.Object, reconciler.DesiredState, error) {
	input := syslogNGConfig{
		TargetHost: s.dataProvider.TargetHost(),
	}

	secret, err := nodeagent.GenerateConfigSecret(input, syslogNGConfigTemplate, s.dataProvider.ConfigFileName())
	secret.SetName(s.dataProvider.QualifiedName("config"))
	s.agentReconciler.ChildObjectMeta(&secret)
	return &secret, reconciler.StatePresent, err
}
