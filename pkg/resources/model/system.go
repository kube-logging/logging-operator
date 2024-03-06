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
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/common"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/input"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/plugins"
)

func CreateSystem(resources LoggingResources, secrets SecretLoaderFactory, logger logr.Logger) (*types.System, error) {
	logging := resources.Logging
	_, fluentdSpec := resources.GetFluentd()

	var forwardInput *input.ForwardInputConfig
	if fluentdSpec != nil && fluentdSpec.ForwardInputConfig != nil {
		forwardInput = fluentdSpec.ForwardInputConfig
	} else {
		forwardInput = input.NewForwardInputConfig()
	}

	if fluentdSpec != nil && fluentdSpec.TLS.Enabled {
		forwardInput.Transport = &common.Transport{
			Version:        "TLSv1_2",
			CaPath:         "/fluentd/tls/ca.crt",
			CertPath:       "/fluentd/tls/tls.crt",
			PrivateKeyPath: "/fluentd/tls/tls.key",
			ClientCertAuth: true,
		}
		forwardInput.Security = &common.Security{
			SelfHostname: "fluentd",
			SharedKey:    fluentdSpec.TLS.SharedKey,
		}
	}

	rootInput, err := forwardInput.ToDirective(secrets.OutputSecretLoaderForNamespace(logging.Spec.ControlNamespace), "main")
	if err != nil {
		return nil, errors.WrapIf(err, "creating root input")
	}

	router := types.NewRouter("main", types.Params{
		"metrics": strconv.FormatBool(fluentdSpec.Metrics != nil),
	})

	var globalFilters []types.Filter
	globalFilters, err = filtersForFilters(
		"globalFilter",
		"globalFilter",
		secrets.OutputSecretLoaderForNamespace(logging.Spec.ControlNamespace),
		logging.Spec.GlobalFilters)

	if err != nil {
		return nil, err
	}

	builder := types.NewSystemBuilder(rootInput, globalFilters, router)

	for _, flowCr := range resources.Fluentd.Flows {
		flow, err := FlowForFlow(flowCr, resources.Fluentd.ClusterOutputs, resources.Fluentd.Outputs, secrets)
		if err != nil {
			if logging.Spec.SkipInvalidResources {
				logger.Error(err, "Flow contains errors, skipping.")
				continue
			} else {
				return nil, err
			}
		}
		err = builder.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	for _, flowCr := range resources.Fluentd.ClusterFlows {
		flow, err := FlowForClusterFlow(flowCr, resources.Fluentd.ClusterOutputs, secrets)
		if err != nil {
			if logging.Spec.SkipInvalidResources {
				logger.Error(err, "ClusterFlow contains errors, skipping.")
				continue
			} else {
				return nil, err
			}
		}
		err = builder.RegisterFlow(flow)
		if err != nil {
			return nil, err
		}
	}
	if resources.Logging.Spec.DefaultFlowSpec != nil {
		flow, err := FlowForDefaultFlow(resources.Logging, resources.Fluentd.ClusterOutputs, secrets)
		if err != nil {
			// TODO set flow status to error?
			return nil, err
		}
		err = builder.RegisterDefaultFlow(flow)
		if err != nil {
			return nil, err
		}
	}

	// Set ErrorOutput
	var errorFlow *types.Flow
	if resources.Logging.Spec.ErrorOutputRef != "" {
		errorFlow, err = FlowForError(resources.Logging.Spec.ErrorOutputRef, resources.Fluentd.ClusterOutputs, secrets)
		if err != nil {
			return nil, err
		}
	} else {
		errorFlow = &types.Flow{
			PluginMeta: types.PluginMeta{
				Directive: "label",
				Tag:       "@ERROR",
			},
			FlowLabel: "@ERROR",
		}
		plugin, err := output.NewNullOutputConfig().ToDirective(nil, "main-fluentd-error")
		if err != nil {
			return nil, err
		}
		errorFlow.WithOutputs(plugin)
	}
	err = builder.RegisterErrorFlow(errorFlow)
	if err != nil {
		return nil, err
	}

	system, err := builder.Build()

	if system != nil && len(system.Flows) == 0 {
		logger.Info("no flows found, generating empty model")
	}

	// TODO: wow such hack
	if fluentdSpec.Workers > 1 {
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

func FlowForError(outputRef string, clusterOutputs ClusterOutputs, secrets SecretLoaderFactory) (*types.Flow, error) {
	errorFlow := &types.Flow{
		PluginMeta: types.PluginMeta{
			Directive: "label",
			Tag:       "@ERROR",
		},
		FlowLabel: "@ERROR",
	}

	if clusterOutput := clusterOutputs.FindByName(outputRef); clusterOutput != nil {
		plugin, err := plugins.CreateOutput(clusterOutput.Spec.OutputSpec, "main-fluentd-error", secrets.OutputSecretLoaderForNamespace(clusterOutput.Namespace))
		if err != nil {
			return nil, errors.WrapIff(err, "failed to create configured output %q", outputRef)
		}
		return errorFlow.WithOutputs(plugin), nil
	}

	return nil, errors.Errorf("there is no ClusterOutput named %s", outputRef)
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

	result, err := types.NewFlow(matches, flowID, flow.Name, flow.Namespace, flow.Spec.FlowLabel, flow.Spec.IncludeLabelInRouter)
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
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %s", outputRef))
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
				errs = errors.Append(errs, errors.WrapIff(err, "failed to create configured output %s/%s", output.Namespace, output.Name))
				continue
			}
			allOutputs = append(allOutputs, plugin)
		} else {
			errs = errors.Append(errs, errors.Errorf("referenced output %s not found for flow %s/%s", outputRef, flow.Namespace, flow.Name))
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

	result, err := types.NewFlow(matches, flowID, flow.Name, flow.Namespace, flow.Spec.FlowLabel, flow.Spec.IncludeLabelInRouter)
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

	result, err := types.NewFlow([]types.FlowMatch{}, flowID, logging.Name, logging.Namespace, logging.Spec.DefaultFlowSpec.FlowLabel, logging.Spec.DefaultFlowSpec.IncludeLabelInRouter)
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
