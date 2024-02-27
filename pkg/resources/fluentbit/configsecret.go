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

package fluentbit

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/resources/syslogng"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

type fluentbitInputConfig struct {
	Values          map[string]string
	ParserN         []string
	MultilineParser []string
}

type fluentbitInputConfigWithTenant struct {
	Tenant          string
	Values          map[string]string
	ParserN         []string
	MultilineParser []string
}

type upstreamNode struct {
	Name string
	Host string
	Port int
}

type upstream struct {
	Name  string
	Path  string
	Nodes []upstreamNode
}

type fluentBitConfig struct {
	Namespace string
	Monitor   struct {
		Enabled bool
		Port    int32
		Path    string
	}
	Flush                   int32
	Grace                   int32
	LogLevel                string
	CoroStackSize           int32
	Output                  map[string]string
	Input                   fluentbitInputConfig
	Inputs                  []fluentbitInputConfigWithTenant
	DisableKubernetesFilter bool
	KubernetesFilter        map[string]string
	AwsFilter               map[string]string
	BufferStorage           map[string]string
	FilterModify            []v1beta1.FilterModify
	FluentForwardOutput     *fluentForwardOutputConfig
	SyslogNGOutput          *syslogNGOutputConfig
	DefaultParsers          string
	CustomParsers           string
	HealthCheck             *v1beta1.HealthCheck
}

type fluentForwardOutputConfig struct {
	Network    FluentbitNetwork
	Options    map[string]string
	Targets    []forwardTargetConfig
	TargetHost string
	TargetPort int32
	TLS        fluentForwardOutputTLSConfig
	Upstream   fluentForwardOutputUpstreamConfig
}

type forwardTargetConfig struct {
	NamespaceRegex string
	Match          string
	Host           string
	Port           int32
}

type fluentForwardOutputTLSConfig struct {
	Enabled   bool
	SharedKey string
}

type fluentForwardOutputUpstreamConfig struct {
	Enabled bool
	Config  upstream
}
type FluentbitNetwork struct {
	ConnectTimeoutSet         bool
	ConnectTimeout            uint32
	ConnectTimeoutLogErrorSet bool
	ConnectTimeoutLogError    bool
	DNSMode                   string
	DNSPreferIPV4Set          bool
	DNSPreferIPV4             bool
	DNSResolver               string
	KeepaliveSet              bool
	Keepalive                 bool
	KeepaliveIdleTimeoutSet   bool
	KeepaliveIdleTimeout      uint32
	KeepaliveMaxRecycleSet    bool
	KeepaliveMaxRecycle       uint32
	SourceAddress             string
}

// https://docs.fluentbit.io/manual/pipeline/outputs/tcp-and-tls
type syslogNGOutputConfig struct {
	Host           string
	Port           int32
	Targets        []forwardTargetConfig
	JSONDateKey    string
	JSONDateFormat string
	Workers        *int
	Network        FluentbitNetwork
}

func newSyslogNGOutputConfig() *syslogNGOutputConfig {
	return &syslogNGOutputConfig{
		JSONDateKey:    "ts",
		JSONDateFormat: "iso8601",
	}
}

func newFluentbitNetwork(network v1beta1.FluentbitNetwork) (result FluentbitNetwork) {
	if network.ConnectTimeout != nil {
		result.ConnectTimeoutSet = true
		result.ConnectTimeout = *network.ConnectTimeout
	}

	if network.ConnectTimeoutLogError != nil {
		result.ConnectTimeoutLogErrorSet = true
		result.ConnectTimeoutLogError = *network.ConnectTimeoutLogError
	}

	if network.DNSMode != "" {
		result.DNSMode = network.DNSMode
	}

	if network.DNSPreferIPV4 != nil {
		result.DNSPreferIPV4Set = true
		result.DNSPreferIPV4 = *network.DNSPreferIPV4
	}

	if network.DNSResolver != "" {
		result.DNSResolver = network.DNSResolver
	}

	if network.Keepalive != nil {
		result.KeepaliveSet = true
		result.Keepalive = *network.Keepalive
	}

	if network.KeepaliveIdleTimeout != nil {
		result.KeepaliveIdleTimeoutSet = true
		result.KeepaliveIdleTimeout = *network.KeepaliveIdleTimeout
	}

	if network.KeepaliveMaxRecycle != nil {
		result.KeepaliveMaxRecycleSet = true
		result.KeepaliveMaxRecycle = *network.KeepaliveMaxRecycle
	}

	if network.SourceAddress != "" {
		result.SourceAddress = network.SourceAddress
	}
	return
}

func (r *Reconciler) configSecret() (runtime.Object, reconciler.DesiredState, error) {
	ctx := context.TODO()
	if r.fluentbitSpec.CustomConfigSecret != "" {
		return &corev1.Secret{
			ObjectMeta: r.FluentbitObjectMeta(fluentBitSecretConfigName),
		}, reconciler.StateAbsent, nil
	}

	disableKubernetesFilter := r.fluentbitSpec.DisableKubernetesFilter != nil && *r.fluentbitSpec.DisableKubernetesFilter

	if !disableKubernetesFilter {
		if r.fluentbitSpec.FilterKubernetes.BufferSize == "" {
			r.logger.Info("Notice: If the Buffer_Size value is empty we will set it 0. For more information: https://github.com/fluent/fluent-bit/issues/2111")
			r.fluentbitSpec.FilterKubernetes.BufferSize = "0"
		} else if r.fluentbitSpec.FilterKubernetes.BufferSize != "0" {
			r.logger.Info("Notice: If the kubernetes filter buffer_size parameter is underestimated it can cause log loss. For more information: https://github.com/fluent/fluent-bit/issues/2111")
		}
		if r.fluentbitSpec.FilterKubernetes.K8SLoggingExclude == "" {
			r.fluentbitSpec.FilterKubernetes.K8SLoggingExclude = "On"
		}
	}

	input := fluentBitConfig{
		Flush:                   r.fluentbitSpec.Flush,
		Grace:                   r.fluentbitSpec.Grace,
		LogLevel:                r.fluentbitSpec.LogLevel,
		CoroStackSize:           r.fluentbitSpec.CoroStackSize,
		Namespace:               r.Logging.Spec.ControlNamespace,
		DisableKubernetesFilter: disableKubernetesFilter,
		FilterModify:            r.fluentbitSpec.FilterModify,
		HealthCheck:             r.fluentbitSpec.HealthCheck,
	}

	input.DefaultParsers = fmt.Sprintf("%s/%s", StockConfigPath, "parsers.conf")

	if r.fluentbitSpec.CustomParsers != "" {
		input.CustomParsers = fmt.Sprintf("%s/%s", OperatorConfigPath, CustomParsersConfigName)
	}

	if r.fluentbitSpec.Metrics != nil {
		input.Monitor.Enabled = true
		input.Monitor.Port = r.fluentbitSpec.Metrics.Port
		input.Monitor.Path = r.fluentbitSpec.Metrics.Path
	}

	if r.fluentbitSpec.InputTail.Parser == "" {
		switch types.ContainerRuntime {
		case "docker":
			r.fluentbitSpec.InputTail.Parser = "docker"
		case "containerd":
			r.fluentbitSpec.InputTail.Parser = "cri"
		default:
			r.fluentbitSpec.InputTail.Parser = "cri"
		}
	}

	loggingResources, errs := r.loggingResourcesRepo.LoggingResourcesFor(ctx, *r.Logging)
	if errs != nil {
		return nil, nil, errs
	}

	_, fluentdSpec := loggingResources.GetFluentd()

	if fluentdSpec != nil {
		fluentbitTargetHost := r.fluentbitSpec.TargetHost
		if fluentbitTargetHost == "" {
			fluentbitTargetHost = aggregatorEndpoint(r.Logging, fluentd.ServiceName)
		}

		fluentbitTargetPort := int32(fluentd.ServicePort)
		if r.fluentbitSpec.TargetPort != 0 {
			fluentbitTargetPort = r.fluentbitSpec.TargetPort
		}

		input.FluentForwardOutput = &fluentForwardOutputConfig{
			TargetHost: fluentbitTargetHost,
			TargetPort: fluentbitTargetPort,
			TLS: fluentForwardOutputTLSConfig{
				Enabled:   *r.fluentbitSpec.TLS.Enabled,
				SharedKey: r.fluentbitSpec.TLS.SharedKey,
			},
		}
	}

	mapper := types.NewStructToStringMapper(nil)

	// FluentBit input Values
	inputTail := r.fluentbitSpec.InputTail

	if len(inputTail.MultilineParser) > 0 {
		input.Input.MultilineParser = inputTail.MultilineParser
		inputTail.MultilineParser = nil

		// If MultilineParser is set, remove other parser fields
		// See https://docs.fluentbit.io/manual/pipeline/inputs/tail#multiline-core-v1.8

		r.logger.Info("Notice: MultilineParser is enabled. Disabling other parser options")

		inputTail.Parser = ""
		inputTail.ParserFirstline = ""
		inputTail.ParserN = nil
		inputTail.Multiline = ""
		inputTail.MultilineFlush = ""
		inputTail.DockerMode = ""
		inputTail.DockerModeFlush = ""
		inputTail.DockerModeParser = ""

	} else if len(inputTail.ParserN) > 0 {
		input.Input.ParserN = inputTail.ParserN
		inputTail.ParserN = nil
	}

	fluentbitInputValues, err := mapper.StringsMap(inputTail)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map container tailer config for fluentbit")
	}
	input.Input.Values = fluentbitInputValues

	input.KubernetesFilter, err = mapper.StringsMap(r.fluentbitSpec.FilterKubernetes)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map kubernetes filter for fluentbit")
	}

	input.BufferStorage, err = mapper.StringsMap(r.fluentbitSpec.BufferStorage)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map buffer storage for fluentbit")
	}

	if r.fluentbitSpec.FilterAws != nil {
		awsFilter, err := mapper.StringsMap(r.fluentbitSpec.FilterAws)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map aws filter for fluentbit")
		}
		input.AwsFilter = awsFilter
	}

	// Fluentd aggregator path
	if input.FluentForwardOutput != nil {
		if r.fluentbitSpec.TargetHost != "" {
			input.FluentForwardOutput.TargetHost = r.fluentbitSpec.TargetHost
		}
		if r.fluentbitSpec.TargetPort != 0 {
			input.FluentForwardOutput.TargetPort = r.fluentbitSpec.TargetPort
		}

		if r.fluentbitSpec.ForwardOptions != nil {
			forwardOptions, err := mapper.StringsMap(r.fluentbitSpec.ForwardOptions)
			if err != nil {
				return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map forwardOptions for fluentbit")
			}
			input.FluentForwardOutput.Options = forwardOptions
		}
		r.applyNetworkSettings(input)

		aggregatorReplicas, err := r.loggingDataProvider.GetReplicaCount(context.TODO())
		if err != nil {
			return nil, nil, errors.WrapIf(err, "getting replica count for fluentd")
		}

		if r.fluentbitSpec.Network == nil && utils.PointerToInt32(aggregatorReplicas) > 1 {
			input.FluentForwardOutput.Network.KeepaliveSet = true
			input.FluentForwardOutput.Network.Keepalive = true
			input.FluentForwardOutput.Network.KeepaliveIdleTimeoutSet = true
			input.FluentForwardOutput.Network.KeepaliveIdleTimeout = 30
			input.FluentForwardOutput.Network.KeepaliveMaxRecycleSet = true
			input.FluentForwardOutput.Network.KeepaliveMaxRecycle = 100
			r.logger.Info("Notice: fluentbit `network` settings have been configured automatically to adapt to multiple aggregator replicas. Configure it manually to avoid this notice.")
		}

		if r.fluentbitSpec.EnableUpstream {
			input.FluentForwardOutput.Upstream.Enabled = true
			input.FluentForwardOutput.Upstream.Config.Path = fmt.Sprintf("%s/%s", OperatorConfigPath, UpstreamConfigName)
			input.FluentForwardOutput.Upstream.Config.Name = "fluentd-upstream"
			for i := int32(0); i < utils.PointerToInt32(aggregatorReplicas); i++ {
				input.FluentForwardOutput.Upstream.Config.Nodes = append(input.FluentForwardOutput.Upstream.Config.Nodes, r.generateUpstreamNode(i))
			}
		}
	}

	_, syslogNGSPec := loggingResources.GetSyslogNGSpec()
	if syslogNGSPec != nil {
		input.SyslogNGOutput = newSyslogNGOutputConfig()
		input.SyslogNGOutput.Host = aggregatorEndpoint(r.Logging, syslogng.ServiceName)
		input.SyslogNGOutput.Port = syslogng.ServicePort
	}

	if len(loggingResources.LoggingRoutes) > 0 {
		var tenants []v1beta1.Tenant
		for _, a := range loggingResources.LoggingRoutes {
			tenants = append(tenants, a.Status.Tenants...)
		}
		if err := r.configureInputsForTenants(tenants, &input); err != nil {
			return nil, nil, errors.WrapIf(err, "configuring inputs for target tenants")
		}
		if err := r.configureOutputsForTenants(ctx, tenants, &input); err != nil {
			return nil, nil, errors.WrapIf(err, "configuring outputs for target tenants")
		}
	} else {
		// compatibility with existing configuration
		if input.FluentForwardOutput != nil {
			input.FluentForwardOutput.Targets = append(input.FluentForwardOutput.Targets, forwardTargetConfig{
				Match: "*",
				Host:  input.FluentForwardOutput.TargetHost,
				Port:  input.FluentForwardOutput.TargetPort,
			})
		} else if input.SyslogNGOutput != nil {
			input.SyslogNGOutput.Targets = append(input.SyslogNGOutput.Targets, forwardTargetConfig{
				Match: "*",
				Host:  input.SyslogNGOutput.Host,
				Port:  input.SyslogNGOutput.Port,
			})
		}
	}

	if input.SyslogNGOutput != nil && r.fluentbitSpec.SyslogNGOutput != nil {
		input.SyslogNGOutput.JSONDateKey = r.fluentbitSpec.SyslogNGOutput.JsonDateKey
		input.SyslogNGOutput.JSONDateFormat = r.fluentbitSpec.SyslogNGOutput.JsonDateFormat
		input.SyslogNGOutput.Workers = r.fluentbitSpec.SyslogNGOutput.Workers
	}

	r.applyNetworkSettings(input)

	conf, err := generateConfig(input)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to generate config for fluentbit")
	}
	confs := map[string][]byte{
		BaseConfigName: []byte(conf),
	}

	if input.FluentForwardOutput != nil && input.FluentForwardOutput.Upstream.Enabled {
		upstreamConfig, err := generateUpstreamConfig(input.FluentForwardOutput.Upstream)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to generate upstream config for fluentbit")
		}
		confs[UpstreamConfigName] = []byte(upstreamConfig)
	}

	if r.fluentbitSpec.CustomParsers != "" {
		confs[CustomParsersConfigName] = []byte(r.fluentbitSpec.CustomParsers)
	}

	r.configs = confs

	return &corev1.Secret{
		ObjectMeta: r.FluentbitObjectMeta(fluentBitSecretConfigName),
		Data:       confs,
	}, reconciler.StatePresent, nil
}

func (r *Reconciler) applyNetworkSettings(input fluentBitConfig) {
	if r.fluentbitSpec.Network != nil {
		if input.FluentForwardOutput != nil {
			input.FluentForwardOutput.Network = newFluentbitNetwork(*r.fluentbitSpec.Network)
		}
		if input.SyslogNGOutput != nil {
			input.SyslogNGOutput.Network = newFluentbitNetwork(*r.fluentbitSpec.Network)
		}
	}
}

func aggregatorEndpoint(l *v1beta1.Logging, svc string) string {
	return fmt.Sprintf("%s.%s.svc%s", l.QualifiedName(svc), l.Spec.ControlNamespace, l.ClusterDomainAsSuffix())
}

func generateConfig(input fluentBitConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl := template.New("fluentbit")
	tmpl, err := tmpl.Parse(fluentBitConfigTemplate)
	if err != nil {
		return "", errors.WrapIf(err, "parsing fluentbit template")
	}
	tmpl, err = tmpl.Parse(fluentbitNetworkTemplate)
	if err != nil {
		return "", errors.WrapIf(err, "parsing fluentbit network nested template")
	}
	tmpl, err = tmpl.Parse(fluentbitInputTemplate)
	if err != nil {
		return "", errors.WrapIf(err, "parsing fluentbit input nested template")
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", errors.WrapIf(err, "executing fluentbit config template")
	}
	outputString := output.String()
	return outputString, nil
}

func generateUpstreamConfig(input fluentForwardOutputUpstreamConfig) (string, error) {
	tmpl, err := template.New("upstream").Parse(upstreamConfigTemplate)
	if err != nil {
		return "", err
	}
	output := new(bytes.Buffer)
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func (r *Reconciler) generateUpstreamNode(index int32) upstreamNode {
	podName := r.Logging.QualifiedName(fmt.Sprintf("%s-%d", fluentd.ComponentFluentd, index))
	return upstreamNode{
		Name: podName,
		Host: fmt.Sprintf("%s.%s.%s.svc%s",
			podName,
			r.Logging.QualifiedName(fluentd.ServiceName+"-headless"),
			r.Logging.Spec.ControlNamespace,
			r.Logging.ClusterDomainAsSuffix()),
		Port: fluentd.ServicePort,
	}
}
