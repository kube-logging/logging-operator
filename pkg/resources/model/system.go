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
	"strconv"

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
	rootInput, err := forwardInput.ToDirective(l.OutputSecretLoaderForNamespace(l.logging.Spec.ControlNamespace), "main")
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create root input")
	}
	router := types.NewRouter("main",
		map[string]string{
			"metrics": strconv.FormatBool(l.logging.Spec.FluentdSpec.Metrics != nil),
		})
	system := types.NewSystem(rootInput, router)
	for _, flowCr := range l.Flows {
		flow, err := l.FlowFromFlow(flowCr)
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
		flow, err := l.FlowFromClusterFlow(flowCr)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = system.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	if l.logging.Spec.DefaultFlowSpec != nil {
		flow, err := l.FlowFromDefaultFlow(*l.logging)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = system.RegisterDefaultFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	if len(l.Flows) == 0 && len(l.ClusterFlows) == 0 && l.logging.Spec.DefaultFlowSpec == nil {
		l.logger.Info("no flows found, generating empty model")
	}
	return system, nil
}

func (l *LoggingResources) ClusterOutputByName(name string) *v1beta1.ClusterOutput {
	for _, output := range l.ClusterOutputs {
		if output.Name == name {
			return &output
		}
	}
	return nil
}

func (l *LoggingResources) OutputByNamespacedName(namespace string, name string) *v1beta1.Output {
	for _, output := range l.Outputs {
		if output.Namespace == namespace && output.Name == name {
			return &output
		}
	}
	return nil
}

func (l *LoggingResources) FiltersForFilters(flowID string, flowName string, secretLoader secret.SecretLoader, filters []v1beta1.Filter) ([]types.Filter, error) {
	var (
		result []types.Filter
		errs   error
	)
	for i, f := range filters {
		id := fmt.Sprintf("%s:%d", flowID, i)
		filter, err := plugins.CreateFilter(f, id, secretLoader)
		if err != nil {
			errs = errors.Append(errs, errors.WrapIff(err, "failed to create filter with index %d for flow %s", i, flowName))
			continue
		}
		result = append(result, filter)
	}
	return result, errs
}

func (l *LoggingResources) OutputSecretLoaderForNamespace(namespace string) secret.SecretLoader {
	return secret.NewSecretLoader(l.client, namespace, fluentd.OutputSecretPath, l.Secrets)
}

func (l *LoggingResources) FlowFromFlow(flow v1beta1.Flow) (*types.Flow, error) {
	if flow.Spec.Match != nil && flow.Spec.Selectors != nil {
		return nil, errors.Errorf("match and selectors cannot be defined simultaneously for flow %s",
			utils.ObjectKeyFromObjectMeta(&flow).String())
	}

	var matches []types.FlowMatch
	if flow.Spec.Match != nil {
		for _, match := range flow.Spec.Match {
			if match.Select != nil && match.Exclude != nil {
				return nil, errors.Errorf("select and exclude cannot be set simultaneously for flow %s",
					utils.ObjectKeyFromObjectMeta(&flow).String())
			}

			if match.Select != nil {
				matches = append(matches, types.FlowMatch{
					Labels:         match.Select.Labels,
					ContainerNames: match.Select.ContainerNames,
					Hosts:          match.Select.Hosts,
					Namespaces:     []string{flow.Namespace},
					Negate:         false,
				})
			}
			if match.Exclude != nil {
				matches = append(matches, types.FlowMatch{
					Labels:         match.Exclude.Labels,
					ContainerNames: match.Exclude.ContainerNames,
					Hosts:          match.Exclude.Hosts,
					Namespaces:     []string{flow.Namespace},
					Negate:         true,
				})
			}
		}
	} else {
		matches = []types.FlowMatch{
			{
				Labels:     flow.Spec.Selectors,
				Namespaces: []string{flow.Namespace},
				Negate:     false,
			},
		}
	}

	flowID := fmt.Sprintf("flow:%s:%s", flow.Namespace, flow.Name)

	result, err := types.NewFlow(matches, flowID, flow.Name, flow.Namespace)
	if err != nil {
		return nil, err
	}

	var errs error

	var outputs []types.Output
	for _, outputRef := range flow.Spec.GlobalOutputRefs {
		if clusterOutput := l.ClusterOutputByName(outputRef); clusterOutput != nil {
			outputID := fmt.Sprintf("%s:clusteroutput:%s:%s", flowID, clusterOutput.Namespace, clusterOutput.Name)
			plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputID, l.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
			if err != nil {
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %q", outputRef))
				continue
			}
			outputs = append(outputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced clusteroutput not found: %s", outputRef))
		}
	}
	for _, outputRef := range flow.Spec.LocalOutputRefs {
		if output := l.OutputByNamespacedName(flow.Namespace, outputRef); output != nil {
			outputID := fmt.Sprintf("%s:output:%s:%s", flowID, output.Namespace, output.Name)
			plugin, err := plugins.CreateOutput(output.Spec, outputID, l.OutputSecretLoaderForNamespace(output.Namespace))
			if err != nil {
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %q", outputRef))
				continue
			}
			outputs = append(outputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced output not found: %s", outputRef))
		}
	}
	result.WithOutputs(outputs...)

	filters, err := l.FiltersForFilters(flowID, flow.Name, l.OutputSecretLoaderForNamespace(flow.Namespace), flow.Spec.Filters)
	errs = errors.Append(errs, err)
	result.WithFilters(filters...)

	return result, errs
}

func (l *LoggingResources) FlowFromClusterFlow(flow v1beta1.ClusterFlow) (*types.Flow, error) {
	if flow.Spec.Match != nil && flow.Spec.Selectors != nil {
		return nil, errors.Errorf("match and selectors cannot be defined simultaneously for clusterflow %s",
			utils.ObjectKeyFromObjectMeta(&flow).String())
	}

	var matches []types.FlowMatch
	if flow.Spec.Match != nil {
		for _, match := range flow.Spec.Match {
			if match.ClusterSelect != nil && match.ClusterExclude != nil {
				return nil, errors.Errorf("select and exclude cannot be set simultaneously for clusterflow %s",
					utils.ObjectKeyFromObjectMeta(&flow).String())
			}

			if match.ClusterSelect != nil {
				matches = append(matches, types.FlowMatch{
					Labels:         match.ClusterSelect.Labels,
					ContainerNames: match.ClusterSelect.ContainerNames,
					Hosts:          match.ClusterSelect.Hosts,
					Namespaces:     match.ClusterSelect.Namespaces,
					Negate:         false,
				})
			}
			if match.ClusterExclude != nil {
				matches = append(matches, types.FlowMatch{
					Labels:         match.ClusterExclude.Labels,
					ContainerNames: match.ClusterExclude.ContainerNames,
					Hosts:          match.ClusterExclude.Hosts,
					Namespaces:     match.ClusterExclude.Namespaces,
					Negate:         true,
				})
			}
		}
	} else {
		matches = []types.FlowMatch{
			{
				Labels:     flow.Spec.Selectors,
				Namespaces: []string{""},
				Negate:     false,
			},
		}
	}

	flowID := fmt.Sprintf("clusterflow:%s:%s", flow.Namespace, flow.Name)

	result, err := types.NewFlow(matches, flowID, flow.Name, flow.Namespace)
	if err != nil {
		return nil, err
	}

	var errs error

	var outputs []types.Output
	for _, outputRef := range flow.Spec.GlobalOutputRefs {
		if clusterOutput := l.ClusterOutputByName(outputRef); clusterOutput != nil {
			outputID := fmt.Sprintf("%s:clusteroutput:%s:%s", flowID, clusterOutput.Namespace, clusterOutput.Name)
			plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputID, l.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
			if err != nil {
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %q", outputRef))
				continue
			}
			outputs = append(outputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced clusteroutput not found: %s", outputRef))
		}
	}
	result.WithOutputs(outputs...)

	filters, err := l.FiltersForFilters(flowID, flow.Name, l.OutputSecretLoaderForNamespace(flow.Namespace), flow.Spec.Filters)
	errs = errors.Append(errs, err)
	result.WithFilters(filters...)

	return result, errs
}

func (l *LoggingResources) FlowFromDefaultFlow(logging v1beta1.Logging) (*types.Flow, error) {
	if logging.Spec.DefaultFlowSpec == nil {
		return nil, nil
	}

	flowID := fmt.Sprintf("logging:%s:%s", logging.Namespace, logging.Name)

	result, err := types.NewFlow([]types.FlowMatch{}, flowID, logging.Name, logging.Namespace)
	if err != nil {
		return nil, err
	}

	var errs error

	var outputs []types.Output
	for _, outputRef := range logging.Spec.DefaultFlowSpec.GlobalOutputRefs {
		if clusterOutput := l.ClusterOutputByName(outputRef); clusterOutput != nil {
			outputID := fmt.Sprintf("%s:clusteroutput:%s:%s", flowID, clusterOutput.Namespace, clusterOutput.Name)
			plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputID, l.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
			if err != nil {
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %q", outputRef))
				continue
			}
			outputs = append(outputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced clusteroutput not found: %s", outputRef))
		}
	}
	result.WithOutputs(outputs...)

	filters, err := l.FiltersForFilters(flowID, logging.Name, l.OutputSecretLoaderForNamespace(logging.Namespace), logging.Spec.DefaultFlowSpec.Filters)
	errs = errors.Append(errs, err)
	result.WithFilters(filters...)

	return result, errs
}
