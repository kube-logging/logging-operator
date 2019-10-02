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
	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/model/common"
	"github.com/banzaicloud/logging-operator/pkg/model/input"
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
	"github.com/banzaicloud/logging-operator/pkg/plugins"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LoggingResources struct {
	client         client.Client
	logger         logr.Logger
	logging        *v1beta1.Logging
	Outputs        []v1beta1.Output
	Flows          []v1beta1.Flow
	ClusterOutputs []v1beta1.ClusterOutput
	ClusterFlows   []v1beta1.ClusterFlow
	Secrets        *secret.MountSecrets
}

func NewLoggingResources(logging *v1beta1.Logging, client client.Client, logger logr.Logger) *LoggingResources {
	return &LoggingResources{
		client:         client,
		logger:         logger,
		logging:        logging,
		Outputs:        make([]v1beta1.Output, 0),
		ClusterOutputs: make([]v1beta1.ClusterOutput, 0),
		Flows:          make([]v1beta1.Flow, 0),
		ClusterFlows:   make([]v1beta1.ClusterFlow, 0),
		Secrets:        &secret.MountSecrets{},
	}
}

func (l *LoggingResources) CreateModel() (*types.Builder, error) {
	forwardInput := input.NewForwardInputConfig()
	if l.logging.Spec.FluentdSpec != nil && l.logging.Spec.FluentdSpec.TLS.Enabled {
		forwardInput.Transport = &common.Transport{
			Version:        "TLSv1_2",
			CaPath:         "/fluentd/tls/ca.crt",
			CertPath:       "/fluentd/tls/tls.crt",
			PrivateKeyPath: "/fluentd/tls/tls.key",
			ClientCertAuth: true,
		}
		forwardInput.Security = &common.Security{
			SelfHostname: "fluentd",
			SharedKey:    l.logging.Spec.FluentdSpec.TLS.SharedKey,
		}
	}
	rootInput, err := forwardInput.ToDirective(secret.NewSecretLoader(l.client, l.logging.Spec.ControlNamespace, fluentd.OutputSecretPath, l.Secrets))
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create root input")
	}
	system := types.NewSystem(rootInput, types.NewRouter())
	for _, flow := range l.Flows {
		flow, err := l.CreateFlowFromCustomResource(flow, flow.Namespace)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = system.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	for _, flowCr := range l.ClusterFlows {
		flow, err := l.CreateFlowFromCustomResource(v1beta1.Flow{
			TypeMeta:   flowCr.TypeMeta,
			ObjectMeta: flowCr.ObjectMeta,
			Spec:       flowCr.Spec,
			Status:     flowCr.Status,
		}, "")
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = system.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	if len(l.Flows) == 0 && len(l.ClusterFlows) == 0 {
		l.logger.Info("no flows found, generating empty model")
	}
	return system, nil
}

func (l *LoggingResources) CreateFlowFromCustomResource(flowCr v1beta1.Flow, namespace string) (*types.Flow, error) {
	flow, err := types.NewFlow(namespace, flowCr.Spec.Selectors)
	if err != nil {
		return nil, err
	}
	outputs := []types.Output{}
	var multierr error
FindOutputForAllRefs:
	for _, outputRef := range flowCr.Spec.OutputRefs {
		// only namespaced flows should use namespaced outputs
		if namespace != "" {
			for _, output := range l.Outputs {
				// only an output from the same namespace can be used with a matching name
				if output.Namespace == namespace && outputRef == output.Name {
					plugin, err := plugins.CreateOutput(output.Spec, secret.NewSecretLoader(l.client, output.Namespace, fluentd.OutputSecretPath, l.Secrets))
					if err != nil {
						multierr = errors.Combine(multierr, errors.WrapIff(err, "failed to create configured output %s", outputRef))
						continue FindOutputForAllRefs
					}
					outputs = append(outputs, plugin)
					continue FindOutputForAllRefs
				}
			}
		}
		for _, clusterOutput := range l.ClusterOutputs {
			if outputRef == clusterOutput.Name {
				plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, secret.NewSecretLoader(l.client, clusterOutput.Namespace, fluentd.OutputSecretPath, l.Secrets))
				if err != nil {
					multierr = errors.Combine(multierr, errors.WrapIff(err, "failed to create configured output %s", outputRef))
					continue FindOutputForAllRefs
				}
				outputs = append(outputs, plugin)
				continue FindOutputForAllRefs
			}
		}
		multierr = errors.Combine(multierr, errors.Errorf("referenced output not found: %s", outputRef))
	}
	flow.WithOutputs(outputs...)

	// Filter
	var filters []types.Filter
	for i, f := range flowCr.Spec.Filters {
		filter, err := plugins.CreateFilter(f, secret.NewSecretLoader(l.client, flowCr.Namespace, fluentd.OutputSecretPath, l.Secrets))
		if err != nil {
			multierr = errors.Combine(multierr, errors.Errorf("failed to create filter with index %d for flow %s", i, flowCr.Name))
			continue
		}
		filters = append(filters, filter)
	}
	flow.WithFilters(filters...)

	return flow, multierr
}

func isEnabled(namespace string, output v1beta1.ClusterOutputSpec) bool {
	for _, enabledNs := range output.EnabledNamespaces {
		if enabledNs == namespace {
			return true
		}
	}
	return false
}
