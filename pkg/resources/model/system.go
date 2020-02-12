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
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta2"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/common"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/input"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/logging-operator/pkg/sdk/plugins"
	"github.com/banzaicloud/operator-tools/pkg/secret"
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
	rootInput, err := forwardInput.ToDirective(secret.NewSecretLoader(l.client, l.logging.Spec.ControlNamespace, fluentd.OutputSecretPath, l.Secrets), "main")
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create root input")
	}
	system := types.NewSystem(rootInput, types.NewRouter("main"))
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

type CommonFlow struct {
	Name       string
	Namespace  string
	OutputRefs []string
	Filters    []v1beta2.Filter
	Flow       *types.Flow
}

func FlowDispatcher(flowCr interface{}) (*CommonFlow, error) {
	var commonFlow *CommonFlow
	switch f := flowCr.(type) {
	case v1beta1.Flow:
		//TODO transform old format into CommonFlow and FlowMatch
		return nil, nil
	case v1beta2.Flow:
		commonFlow = &CommonFlow{
			Name:       f.Name,
			Namespace:  f.Namespace,
			OutputRefs: f.Spec.OutputRefs,
			Filters:    f.Spec.Filters,
		}

		flowMatches := types.GetFlowMatchFromSpec(f.Namespace, f.Spec.Match)
		flow, err := types.NewFlow(f.Namespace, flowMatches)
		if err != nil {
			return nil, err
		}
		commonFlow.Flow = flow
		return commonFlow, nil
	case v1beta2.ClusterFlow:
		commonFlow = &CommonFlow{
			Name:       f.Name,
			Namespace:  f.Namespace,
			OutputRefs: f.Spec.OutputRefs,
			Filters:    f.Spec.Filters,
		}

		flowMatches := types.GetFlowMatchFromSpec(f.Namespace, f.Spec.Match)
		flow, err := types.NewFlow(f.Namespace, flowMatches)
		if err != nil {
			return nil, err
		}
		commonFlow.Flow = flow
		return commonFlow, nil
	}
	return nil, fmt.Errorf("unsupported type: %t", flowCr)
}

func (l *LoggingResources) CreateFlowFromCustomResource(flowCr interface{}, namespace string) (*types.Flow, error) {
	commonFlow, err := FlowDispatcher(flowCr)
	if err != nil {
		return nil, err
	}
	flow := commonFlow.Flow
	outputs := []types.Output{}
	var multierr error
FindOutputForAllRefs:
	for _, outputRef := range commonFlow.OutputRefs {
		// only namespaced flows should use namespaced outputs
		if namespace != "" {
			for _, output := range l.Outputs {
				// only an output from the same namespace can be used with a matching name
				if output.Namespace == namespace && outputRef == output.Name {
					outputId := namespace + "_" + commonFlow.Name + "_" + output.Name
					plugin, err := plugins.CreateOutput(output.Spec, outputId, secret.NewSecretLoader(l.client, output.Namespace, fluentd.OutputSecretPath, l.Secrets))
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
				outputId := namespace + "_" + commonFlow.Name + "_" + clusterOutput.Name
				plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputId, secret.NewSecretLoader(l.client, clusterOutput.Namespace, fluentd.OutputSecretPath, l.Secrets))
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
	for i, f := range commonFlow.Filters {
		filter, err := plugins.CreateFilter(f, commonFlow.Name, i, secret.NewSecretLoader(l.client, commonFlow.Namespace, fluentd.OutputSecretPath, l.Secrets))
		if err != nil {
			multierr = errors.Combine(multierr, errors.WrapIff(err, "failed to create filter with index %d for flow %s", i, commonFlow.Name))
			continue
		}
		filters = append(filters, filter)
	}
	flow.WithFilters(filters...)

	return flow, multierr
}
