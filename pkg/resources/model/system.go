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
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/go-logr/logr"

	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/common"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/input"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/logging-operator/pkg/sdk/plugins"
)

func CreateSystem(resources LoggingResources, secrets SecretLoaderFactory, logger logr.Logger) (*types.System, error) {
	logging := resources.Logging

	var forwardInput *input.ForwardInputConfig
	if logging.Spec.FluentdSpec.ForwardInputConfig != nil {
		forwardInput = logging.Spec.FluentdSpec.ForwardInputConfig
	} else {
		forwardInput = input.NewForwardInputConfig()
	}

	if logging.Spec.FluentdSpec != nil && logging.Spec.FluentdSpec.TLS.Enabled {
		forwardInput.Transport = &common.Transport{
			Version:        "TLSv1_2",
			CaPath:         "/fluentd/tls/ca.crt",
			CertPath:       "/fluentd/tls/tls.crt",
			PrivateKeyPath: "/fluentd/tls/tls.key",
			ClientCertAuth: true,
		}
		forwardInput.Security = &common.Security{
			SelfHostname: "fluentd",
			SharedKey:    logging.Spec.FluentdSpec.TLS.SharedKey,
		}
	}

	rootInput, err := forwardInput.ToDirective(secrets.OutputSecretLoaderForNamespace(logging.Spec.ControlNamespace), "main")
	if err != nil {
		return nil, errors.WrapIf(err, "creating root input")
	}

	router := types.NewRouter("main", types.Params{
		"metrics": strconv.FormatBool(logging.Spec.FluentdSpec.Metrics != nil),
	})

	globalFilters, err := filtersForFilters(
		"globalFilter",
		"globalFilter",
		secrets.OutputSecretLoaderForNamespace(logging.Spec.ControlNamespace),
		logging.Spec.GlobalFilters)

	if err != nil {
		return nil, err
	}

	builder := types.NewSystemBuilder(rootInput, globalFilters, router)

	for _, flowCr := range resources.Flows {
		flow, err := FlowForFlow(flowCr, resources.ClusterOutputs, resources.Outputs, secrets)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = builder.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	for _, flowCr := range resources.ClusterFlows {
		flow, err := FlowForClusterFlow(flowCr, resources.ClusterOutputs, secrets)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = builder.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	if resources.Logging.Spec.DefaultFlowSpec != nil {
		flow, err := FlowForDefaultFlow(resources.Logging, resources.ClusterOutputs, secrets)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = builder.RegisterDefaultFlow(flow)
		if err != nil {
			return nil, err
		}
	}

	system, err := builder.Build()

	if system != nil && len(system.Flows) == 0 {
		logger.Info("no flows found, generating empty model")
	}

	// TODO: wow such hack
	if logging.Spec.FluentdSpec.Workers > 1 {
		for _, flow := range system.Flows {
			for _, output := range flow.Outputs {
				unsetBufferPath(output)
			}
		}
	}

	return system, err
}

func unsetBufferPath(directive types.Directive) {
	if gd, _ := directive.(*types.GenericDirective); gd != nil && gd.Directive == "buffer" {
		delete(gd.Params, "path")
		return
	}
	for _, d := range directive.GetSections() {
		unsetBufferPath(d)
	}
}

type SecretLoaderFactory interface {
	OutputSecretLoaderForNamespace(namespace string) secret.SecretLoader
}

func filtersForFilters(flowID string, flowName string, secretLoader secret.SecretLoader, filters []v1beta1.Filter) ([]types.Filter, error) {
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

func FlowForFlow(flow v1beta1.Flow, clusterOutputs ClusterOutputs, outputs Outputs, secrets SecretLoaderFactory) (*types.Flow, error) {
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

	var allOutputs []types.Output
	for _, outputRef := range flow.Spec.GlobalOutputRefs {
		if clusterOutput := clusterOutputs.FindByName(outputRef); clusterOutput != nil {
			outputID := fmt.Sprintf("%s:clusteroutput:%s:%s", flowID, clusterOutput.Namespace, clusterOutput.Name)
			plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputID, secrets.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
			if err != nil {
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %q", outputRef))
				continue
			}
			allOutputs = append(allOutputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced clusteroutput not found: %s", outputRef))
		}
	}
	for _, outputRef := range flow.Spec.LocalOutputRefs {
		if output := outputs.FindByNamespacedName(flow.Namespace, outputRef); output != nil {
			outputID := fmt.Sprintf("%s:output:%s:%s", flowID, output.Namespace, output.Name)
			plugin, err := plugins.CreateOutput(output.Spec, outputID, secrets.OutputSecretLoaderForNamespace(output.Namespace))
			if err != nil {
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %q", outputRef))
				continue
			}
			allOutputs = append(allOutputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced output not found: %s", outputRef))
		}
	}
	result.WithOutputs(allOutputs...)

	filters, err := filtersForFilters(flowID, flow.Name, secrets.OutputSecretLoaderForNamespace(flow.Namespace), flow.Spec.Filters)
	errs = errors.Append(errs, err)
	result.WithFilters(filters...)

	return result, errs
}

func FlowForClusterFlow(flow v1beta1.ClusterFlow, clusterOutputs ClusterOutputs, secrets SecretLoaderFactory) (*types.Flow, error) {
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
		if clusterOutput := clusterOutputs.FindByName(outputRef); clusterOutput != nil {
			outputID := fmt.Sprintf("%s:clusteroutput:%s:%s", flowID, clusterOutput.Namespace, clusterOutput.Name)
			plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputID, secrets.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
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

	filters, err := filtersForFilters(flowID, flow.Name, secrets.OutputSecretLoaderForNamespace(flow.Namespace), flow.Spec.Filters)
	errs = errors.Append(errs, err)
	result.WithFilters(filters...)

	return result, errs
}

func FlowForDefaultFlow(logging v1beta1.Logging, clusterOutputs ClusterOutputs, secrets SecretLoaderFactory) (*types.Flow, error) {
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
		if clusterOutput := clusterOutputs.FindByName(outputRef); clusterOutput != nil {
			outputID := fmt.Sprintf("%s:clusteroutput:%s:%s", flowID, clusterOutput.Namespace, clusterOutput.Name)
			plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, outputID, secrets.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
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

	filters, err := filtersForFilters(flowID, logging.Name, secrets.OutputSecretLoaderForNamespace(logging.Namespace), logging.Spec.DefaultFlowSpec.Filters)
	errs = errors.Append(errs, err)
	result.WithFilters(filters...)

	return result, errs
}
