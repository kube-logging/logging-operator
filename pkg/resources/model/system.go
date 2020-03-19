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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/common"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/input"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/logging-operator/pkg/sdk/plugins"
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/banzaicloud/operator-tools/pkg/utils"
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
	for _, flowCr := range l.Flows {
		flow, err := l.CreateFlowFromCustomResource(flowCr)
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
		flow, err := l.CreateFlowFromCustomResource(flowCr)
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
	Scope      string
	OutputRefs []string
	Filters    []v1beta1.Filter
	Flow       *types.Flow
}

// Create []FlowMatch from v1beta1.Match and v1beta1.ClusterMatch
func GetFlowMatchFromSpec(namespace string, matches interface{}) ([]types.FlowMatch, error) {
	var flowMatches []types.FlowMatch
	switch matchList := matches.(type) {
	case []v1beta1.Match:
		for _, match := range matchList {
			if match.Select != nil && match.Exclude != nil {
				return nil, errors.Errorf("select and exclude cannot be set simultaneously")
			}
			if match.Select != nil {
				flowMatches = append(flowMatches, types.FlowMatch{
					Labels:     match.Select.Labels,
					Namespaces: []string{namespace},
					Negate:     false,
				})
			}
			if match.Exclude != nil {
				flowMatches = append(flowMatches, types.FlowMatch{
					Labels:     match.Exclude.Labels,
					Namespaces: []string{namespace},
					Negate:     true,
				})
			}
		}
	case []v1beta1.ClusterMatch:
		for _, match := range matchList {
			if match.ClusterSelect != nil && match.ClusterExclude != nil {
				return nil, errors.Errorf("select and exclude cannot be set simultaneously")
			}
			if match.ClusterSelect != nil {
				flowMatches = append(flowMatches, types.FlowMatch{
					Labels:     match.ClusterSelect.Labels,
					Namespaces: match.ClusterSelect.Namespaces,
					Negate:     false,
				})
			}
			if match.ClusterExclude != nil {
				flowMatches = append(flowMatches, types.FlowMatch{
					Labels:     match.ClusterExclude.Labels,
					Namespaces: match.ClusterExclude.Namespaces,
					Negate:     true,
				})
			}
		}
	}
	return flowMatches, nil
}

func FlowDispatcher(flowCr interface{}) (*CommonFlow, error) {
	var commonFlow *CommonFlow
	var err error
	switch f := flowCr.(type) {
	case v1beta1.ClusterFlow:
		var matches []types.FlowMatch
		commonFlow = &CommonFlow{
			Name:       f.Name,
			Namespace:  f.Namespace,
			Scope:      "",
			OutputRefs: f.Spec.OutputRefs,
			Filters:    f.Spec.Filters,
		}
		if f.Spec.Match != nil && f.Spec.Selectors != nil {
			return nil, errors.Errorf("match and selectors cannot be defined simultaneously for clusterflow %s",
				utils.ObjectKeyFromObjectMeta(&f).String())
		}
		if f.Spec.Match != nil {
			matches, err = GetFlowMatchFromSpec(f.Namespace, f.Spec.Match)
			if err != nil {
				return nil, errors.WrapIff(err, "failed to process match for %s", utils.ObjectKeyFromObjectMeta(&f).String())
			}
		} else {
			matches = []types.FlowMatch{
				{
					Labels:     f.Spec.Selectors,
					Namespaces: []string{""},
					Negate:     false,
				},
			}
		}
		flow, err := types.NewFlow(matches, f.Name, f.Namespace)
		if err != nil {
			return nil, err
		}
		commonFlow.Flow = flow
		return commonFlow, nil
	case v1beta1.Flow:
		var matches []types.FlowMatch
		commonFlow = &CommonFlow{
			Name:       f.Name,
			Namespace:  f.Namespace,
			Scope:      f.Namespace,
			OutputRefs: f.Spec.OutputRefs,
			Filters:    f.Spec.Filters,
		}
		if f.Spec.Match != nil && f.Spec.Selectors != nil {
			return nil, errors.Errorf("match and selectors cannot be defined simultaneously for flow %s",
				utils.ObjectKeyFromObjectMeta(&f).String())
		}
		if f.Spec.Match != nil {
			matches, err = GetFlowMatchFromSpec(f.Namespace, f.Spec.Match)
			if err != nil {
				return nil, errors.WrapIff(err, "failed to process match for %s", utils.ObjectKeyFromObjectMeta(&f).String())
			}
		} else {
			matches = []types.FlowMatch{
				{
					Labels:     f.Spec.Selectors,
					Namespaces: []string{f.Namespace},
					Negate:     false,
				},
			}
		}
		flow, err := types.NewFlow(matches, f.Name, f.Namespace)
		commonFlow.Flow = flow
		if err != nil {
			return nil, err
		}
		return commonFlow, nil
	}
	return nil, fmt.Errorf("unsupported type: %t", flowCr)
}

func (l *LoggingResources) CreateFlowFromCustomResource(flowCr interface{}) (*types.Flow, error) {
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
		if commonFlow.Scope != "" {
			for _, output := range l.Outputs {
				// only an output from the same namespace can be used with a matching name
				if output.Namespace == commonFlow.Scope && outputRef == output.Name {
					outputId := commonFlow.Scope + "_" + commonFlow.Name + "_" + output.Name
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
				outputId := commonFlow.Namespace + "_" + commonFlow.Name + "_" + clusterOutput.Name
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
