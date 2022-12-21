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
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/resources/syslogng"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/types"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type fluentbitInputConfig struct {
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
	DisableKubernetesFilter bool
	KubernetesFilter        map[string]string
	AwsFilter               map[string]string
	BufferStorage           map[string]string
	FilterModify            []v1beta1.FilterModify
	FluentForwardOutput     *fluentForwardOutputConfig
	SyslogNGOutput          *syslogNGOutputConfig
}

type fluentForwardOutputConfig struct {
	Network    FluentbitNetwork
	Options    map[string]string
	TargetHost string
	TargetPort int32
	TLS        fluentForwardOutputTLSConfig
	Upstream   fluentForwardOutputUpstreamConfig
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
	Port           int
	JSONDateKey    string
	JSONDateFormat string
	Workers        *int
	Network        FluentbitNetwork
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
	if r.Logging.Spec.FluentbitSpec.CustomConfigSecret != "" {
		return &corev1.Secret{
			ObjectMeta: r.FluentbitObjectMeta(fluentBitSecretConfigName),
		}, reconciler.StateAbsent, nil
	}

	disableKubernetesFilter := r.Logging.Spec.FluentbitSpec.DisableKubernetesFilter != nil && *r.Logging.Spec.FluentbitSpec.DisableKubernetesFilter

	if !disableKubernetesFilter {
		if r.Logging.Spec.FluentbitSpec.FilterKubernetes.BufferSize == "" {
			log.Log.Info("Notice: If the Buffer_Size value is empty we will set it 0. For more information: https://github.com/fluent/fluent-bit/issues/2111")
			r.Logging.Spec.FluentbitSpec.FilterKubernetes.BufferSize = "0"
		} else if r.Logging.Spec.FluentbitSpec.FilterKubernetes.BufferSize != "0" {
			log.Log.Info("Notice: If the kubernetes filter buffer_size parameter is underestimated it can cause log loss. For more information: https://github.com/fluent/fluent-bit/issues/2111")
		}
	}

	input := fluentBitConfig{
		Flush:                   r.Logging.Spec.FluentbitSpec.Flush,
		Grace:                   r.Logging.Spec.FluentbitSpec.Grace,
		LogLevel:                r.Logging.Spec.FluentbitSpec.LogLevel,
		CoroStackSize:           r.Logging.Spec.FluentbitSpec.CoroStackSize,
		Namespace:               r.Logging.Spec.ControlNamespace,
		DisableKubernetesFilter: disableKubernetesFilter,
		FilterModify:            r.Logging.Spec.FluentbitSpec.FilterModify,
	}

	if r.Logging.Spec.FluentbitSpec.Metrics != nil {
		input.Monitor.Enabled = true
		input.Monitor.Port = r.Logging.Spec.FluentbitSpec.Metrics.Port
		input.Monitor.Path = r.Logging.Spec.FluentbitSpec.Metrics.Path
	}

	if r.Logging.Spec.FluentbitSpec.InputTail.Parser == "" {
		switch types.ContainerRuntime {
		case "docker":
			r.Logging.Spec.FluentbitSpec.InputTail.Parser = "docker"
		case "containerd":
			r.Logging.Spec.FluentbitSpec.InputTail.Parser = "cri"
		default:
			r.Logging.Spec.FluentbitSpec.InputTail.Parser = "cri"
		}
	}

	if r.Logging.Spec.FluentdSpec != nil {
		fluentbitTargetHost := r.Logging.Spec.FluentbitSpec.TargetHost
		if fluentbitTargetHost == "" {
			fluentbitTargetHost = fmt.Sprintf("%s.%s.svc%s", r.Logging.QualifiedName(fluentd.ServiceName), r.Logging.Spec.ControlNamespace, r.Logging.ClusterDomainAsSuffix())
		}

		fluentbitTargetPort := r.Logging.Spec.FluentbitSpec.TargetPort
		if fluentbitTargetPort == 0 {
			fluentbitTargetPort = r.Logging.Spec.FluentdSpec.Port
		}

		input.FluentForwardOutput = &fluentForwardOutputConfig{
			TargetHost: fluentbitTargetHost,
			TargetPort: fluentbitTargetPort,
			TLS: fluentForwardOutputTLSConfig{
				Enabled:   *r.Logging.Spec.FluentbitSpec.TLS.Enabled,
				SharedKey: r.Logging.Spec.FluentbitSpec.TLS.SharedKey,
			},
		}
	}

	mapper := types.NewStructToStringMapper(nil)

	// FluentBit input Values
	inputTail := r.Logging.Spec.FluentbitSpec.InputTail

	if len(inputTail.MultilineParser) > 0 {
		input.Input.MultilineParser = inputTail.MultilineParser
		inputTail.MultilineParser = nil

		// If MultilineParser is set, remove other parser fields
		// See https://docs.fluentbit.io/manual/pipeline/inputs/tail#multiline-core-v1.8

		log.Log.Info("Notice: MultilineParser is enabled. Disabling other parser options",
			"logging", r.Logging.Name)

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

	input.KubernetesFilter, err = mapper.StringsMap(r.Logging.Spec.FluentbitSpec.FilterKubernetes)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map kubernetes filter for fluentbit")
	}

	input.BufferStorage, err = mapper.StringsMap(r.Logging.Spec.FluentbitSpec.BufferStorage)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map buffer storage for fluentbit")
	}

	if r.Logging.Spec.FluentbitSpec.FilterAws != nil {
		awsFilter, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.FilterAws)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map aws filter for fluentbit")
		}
		input.AwsFilter = awsFilter
	}

	if input.FluentForwardOutput != nil {
		if r.Logging.Spec.FluentbitSpec.TargetHost != "" {
			input.FluentForwardOutput.TargetHost = r.Logging.Spec.FluentbitSpec.TargetHost
		}
		if r.Logging.Spec.FluentbitSpec.TargetPort != 0 {
			input.FluentForwardOutput.TargetPort = r.Logging.Spec.FluentbitSpec.TargetPort
		}
		if r.Logging.Spec.FluentbitSpec.ForwardOptions != nil {
			forwardOptions, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.ForwardOptions)
			if err != nil {
				return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map forwardOptions for fluentbit")
			}
			input.FluentForwardOutput.Options = forwardOptions
		}
		if r.Logging.Spec.FluentbitSpec.Network != nil {
			input.FluentForwardOutput.Network = newFluentbitNetwork(*r.Logging.Spec.FluentbitSpec.Network)
		}

		fluentdReplicas, err := r.fluentdDataProvider.GetReplicaCount(context.TODO(), r.Logging)
		if err != nil {
			return nil, nil, errors.WrapIf(err, "getting replica count for fluentd")
		}

		if r.Logging.Spec.FluentbitSpec.Network == nil && utils.PointerToInt32(fluentdReplicas) > 1 {
			input.FluentForwardOutput.Network.KeepaliveSet = true
			input.FluentForwardOutput.Network.Keepalive = true
			input.FluentForwardOutput.Network.KeepaliveIdleTimeoutSet = true
			input.FluentForwardOutput.Network.KeepaliveIdleTimeout = 30
			input.FluentForwardOutput.Network.KeepaliveMaxRecycleSet = true
			input.FluentForwardOutput.Network.KeepaliveMaxRecycle = 100
			log.Log.Info("Notice: Because the Fluentd statefulset has been scaled, we've made some changes in the fluentbit network config too. We advice to revise these default configurations.")
		}

		if r.Logging.Spec.FluentbitSpec.EnableUpstream {
			input.FluentForwardOutput.Upstream.Enabled = true
			input.FluentForwardOutput.Upstream.Config.Name = "fluentd-upstream"
			for i := int32(0); i < utils.PointerToInt32(fluentdReplicas); i++ {
				input.FluentForwardOutput.Upstream.Config.Nodes = append(input.FluentForwardOutput.Upstream.Config.Nodes, r.generateUpstreamNode(i))
			}
		}
	}

	if r.Logging.Spec.SyslogNGSpec != nil {
		input.SyslogNGOutput = &syslogNGOutputConfig{}
		input.SyslogNGOutput.Host = fmt.Sprintf("%s.%s.svc.cluster.local", r.Logging.QualifiedName(syslogng.ServiceName), r.Logging.Spec.ControlNamespace)
		input.SyslogNGOutput.Port = syslogng.ServicePort
		input.SyslogNGOutput.JSONDateKey = "ts"
		input.SyslogNGOutput.JSONDateFormat = "iso8601"

		if r.Logging.Spec.FluentbitSpec.SyslogNGOutput != nil {
			input.SyslogNGOutput.JSONDateKey = r.Logging.Spec.FluentbitSpec.SyslogNGOutput.JsonDateKey
			input.SyslogNGOutput.JSONDateFormat = r.Logging.Spec.FluentbitSpec.SyslogNGOutput.JsonDateFormat
		}
	}

	if input.SyslogNGOutput != nil {
		if r.Logging.Spec.FluentbitSpec.Network != nil {
			input.SyslogNGOutput.Network = newFluentbitNetwork(*r.Logging.Spec.FluentbitSpec.Network)
		}
	}

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

	r.configs = confs

	return &corev1.Secret{
		ObjectMeta: r.FluentbitObjectMeta(fluentBitSecretConfigName),
		Data:       confs,
	}, reconciler.StatePresent, nil
}

func generateConfig(input fluentBitConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentBitConfigTemplate)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", err
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
		Port: 24240,
	}
}
