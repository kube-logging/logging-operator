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
	TLS       struct {
		Enabled   bool
		SharedKey string
	}
	Monitor struct {
		Enabled bool
		Port    int32
		Path    string
	}
	Flush                   int32
	Grace                   int32
	LogLevel                string
	CoroStackSize           int32
	Output                  map[string]string
	TargetHost              string
	TargetPort              int32
	Input                   fluentbitInputConfig
	DisableKubernetesFilter bool
	KubernetesFilter        map[string]string
	AwsFilter               map[string]string
	BufferStorage           map[string]string
	Network                 struct {
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
	ForwardOptions map[string]string
	Upstream       struct {
		Enabled bool
		Config  upstream
	}
}

func (r *Reconciler) configSecret() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.FluentbitSpec.CustomConfigSecret != "" {
		return &corev1.Secret{
			ObjectMeta: r.FluentbitObjectMeta(fluentBitSecretConfigName),
		}, reconciler.StateAbsent, nil
	}
	monitor := struct {
		Enabled bool
		Port    int32
		Path    string
	}{}
	if r.Logging.Spec.FluentbitSpec.Metrics != nil {
		monitor.Enabled = true
		monitor.Port = r.Logging.Spec.FluentbitSpec.Metrics.Port
		monitor.Path = r.Logging.Spec.FluentbitSpec.Metrics.Path
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

	var fluentbitTargetHost string
	if r.Logging.Spec.FluentdSpec != nil && r.Logging.Spec.FluentbitSpec.TargetHost == "" {
		fluentbitTargetHost = fmt.Sprintf("%s.%s.svc.cluster.local", r.Logging.QualifiedName(fluentd.ServiceName), r.Logging.Spec.ControlNamespace)
	} else {
		fluentbitTargetHost = r.Logging.Spec.FluentbitSpec.TargetHost
	}

	var fluentbitTargetPort int32
	if r.Logging.Spec.FluentdSpec != nil && r.Logging.Spec.FluentbitSpec.TargetPort == 0 {
		fluentbitTargetPort = r.Logging.Spec.FluentdSpec.Port
	} else {
		fluentbitTargetPort = r.Logging.Spec.FluentbitSpec.TargetPort
	}

	mapper := types.NewStructToStringMapper(nil)

	// FluentBit input Values
	fluentbitInput := fluentbitInputConfig{}
	inputTail := r.Logging.Spec.FluentbitSpec.InputTail

	if len(inputTail.MultilineParser) > 0 {
		fluentbitInput.MultilineParser = inputTail.MultilineParser
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
		fluentbitInput.ParserN = inputTail.ParserN
		inputTail.ParserN = nil
	}

	fluentbitInputValues, err := mapper.StringsMap(inputTail)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map container tailer config for fluentbit")
	}
	fluentbitInput.Values = fluentbitInputValues

	disableKubernetesFilter := r.Logging.Spec.FluentbitSpec.DisableKubernetesFilter != nil && *r.Logging.Spec.FluentbitSpec.DisableKubernetesFilter

	if !disableKubernetesFilter {
		if r.Logging.Spec.FluentbitSpec.FilterKubernetes.BufferSize == "" {
			log.Log.Info("Notice: If the Buffer_Size value is empty we will set it 0. For more information: https://github.com/fluent/fluent-bit/issues/2111")
			r.Logging.Spec.FluentbitSpec.FilterKubernetes.BufferSize = "0"
		} else if r.Logging.Spec.FluentbitSpec.FilterKubernetes.BufferSize != "0" {
			log.Log.Info("Notice: If the kubernetes filter buffer_size parameter is underestimated it can cause log loss. For more information: https://github.com/fluent/fluent-bit/issues/2111")
		}
	}

	fluentbitKubernetesFilter, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.FilterKubernetes)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map kubernetes filter for fluentbit")
	}

	fluentbitBufferStorage, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.BufferStorage)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map buffer storage for fluentbit")
	}

	input := fluentBitConfig{
		Flush:         r.Logging.Spec.FluentbitSpec.Flush,
		Grace:         r.Logging.Spec.FluentbitSpec.Grace,
		LogLevel:      r.Logging.Spec.FluentbitSpec.LogLevel,
		CoroStackSize: r.Logging.Spec.FluentbitSpec.CoroStackSize,
		Namespace:     r.Logging.Spec.ControlNamespace,
		TLS: struct {
			Enabled   bool
			SharedKey string
		}{
			Enabled:   *r.Logging.Spec.FluentbitSpec.TLS.Enabled,
			SharedKey: r.Logging.Spec.FluentbitSpec.TLS.SharedKey,
		},
		Monitor:                 monitor,
		TargetHost:              fluentbitTargetHost,
		TargetPort:              fluentbitTargetPort,
		Input:                   fluentbitInput,
		DisableKubernetesFilter: disableKubernetesFilter,
		KubernetesFilter:        fluentbitKubernetesFilter,
		BufferStorage:           fluentbitBufferStorage,
	}
	if r.Logging.Spec.FluentbitSpec.FilterAws != nil {
		awsFilter, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.FilterAws)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map aws filter for fluentbit")
		}
		input.AwsFilter = awsFilter
	}
	if r.Logging.Spec.FluentbitSpec.TargetHost != "" {
		input.TargetHost = r.Logging.Spec.FluentbitSpec.TargetHost
	}
	if r.Logging.Spec.FluentbitSpec.TargetPort != 0 {
		input.TargetPort = r.Logging.Spec.FluentbitSpec.TargetPort
	}
	if r.Logging.Spec.FluentbitSpec.ForwardOptions != nil {
		forwardOptions, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.ForwardOptions)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map forwardOptions for fluentbit")
		}
		input.ForwardOptions = forwardOptions
	}

	if r.Logging.Spec.FluentbitSpec.Network != nil {
		if r.Logging.Spec.FluentbitSpec.Network.ConnectTimeout != nil {
			input.Network.ConnectTimeoutSet = true
			input.Network.ConnectTimeout = *r.Logging.Spec.FluentbitSpec.Network.ConnectTimeout
		}

		if r.Logging.Spec.FluentbitSpec.Network.ConnectTimeoutLogError != nil {
			input.Network.ConnectTimeoutLogErrorSet = true
			input.Network.ConnectTimeoutLogError = *r.Logging.Spec.FluentbitSpec.Network.ConnectTimeoutLogError
		}

		if r.Logging.Spec.FluentbitSpec.Network.DNSMode != "" {
			input.Network.DNSMode = r.Logging.Spec.FluentbitSpec.Network.DNSMode
		}

		if r.Logging.Spec.FluentbitSpec.Network.DNSPreferIPV4 != nil {
			input.Network.DNSPreferIPV4Set = true
			input.Network.DNSPreferIPV4 = *r.Logging.Spec.FluentbitSpec.Network.DNSPreferIPV4
		}

		if r.Logging.Spec.FluentbitSpec.Network.DNSResolver != "" {
			input.Network.DNSResolver = r.Logging.Spec.FluentbitSpec.Network.DNSResolver
		}

		if r.Logging.Spec.FluentbitSpec.Network.Keepalive != nil {
			input.Network.KeepaliveSet = true
			input.Network.Keepalive = *r.Logging.Spec.FluentbitSpec.Network.Keepalive
		}

		if r.Logging.Spec.FluentbitSpec.Network.KeepaliveIdleTimeout != nil {
			input.Network.KeepaliveIdleTimeoutSet = true
			input.Network.KeepaliveIdleTimeout = *r.Logging.Spec.FluentbitSpec.Network.KeepaliveIdleTimeout
		}

		if r.Logging.Spec.FluentbitSpec.Network.KeepaliveMaxRecycle != nil {
			input.Network.KeepaliveMaxRecycleSet = true
			input.Network.KeepaliveMaxRecycle = *r.Logging.Spec.FluentbitSpec.Network.KeepaliveMaxRecycle
		}

		if r.Logging.Spec.FluentbitSpec.Network.SourceAddress != "" {
			input.Network.SourceAddress = r.Logging.Spec.FluentbitSpec.Network.SourceAddress
		}
	}

	fluentdReplicas, err := r.fluentdDataProvider.GetReplicaCount(context.TODO(), r.Logging)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "getting replica count for fluentd")
	}

	if r.Logging.Spec.FluentbitSpec.Network == nil && utils.PointerToInt32(fluentdReplicas) > 1 {
		input.Network.KeepaliveSet = true
		input.Network.Keepalive = true
		input.Network.KeepaliveIdleTimeoutSet = true
		input.Network.KeepaliveIdleTimeout = 30
		input.Network.KeepaliveMaxRecycleSet = true
		input.Network.KeepaliveMaxRecycle = 100
		log.Log.Info("Notice: Because the Fluentd statefulset has been scaled, we've made some changes in the fluentbit network config too. We advice to revise these default configurations.")
	}

	if r.Logging.Spec.FluentbitSpec.EnableUpstream {
		input.Upstream.Enabled = true
		input.Upstream.Config.Name = "fluentd-upstream"
		for i := int32(0); i < utils.PointerToInt32(fluentdReplicas); i++ {
			input.Upstream.Config.Nodes = append(input.Upstream.Config.Nodes, r.generateUpstreamNode(i))
		}
	}

	conf, err := generateConfig(input)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to generate config for fluentbit")
	}
	confs := map[string][]byte{
		BaseConfigName: []byte(conf),
	}

	if input.Upstream.Enabled {
		upstreamConfig, err := generateUpstreamConfig(input)
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

func generateUpstreamConfig(input fluentBitConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("upstream").Parse(upstreamConfigTemplate)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(output, input.Upstream)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func (r *Reconciler) generateUpstreamNode(index int32) upstreamNode {
	podName := r.Logging.QualifiedName(fmt.Sprintf("%s-%d", fluentd.ComponentFluentd, index))
	return upstreamNode{
		Name: podName,
		Host: fmt.Sprintf("%s.%s.%s.svc.cluster.local",
			podName,
			r.Logging.QualifiedName(fluentd.ServiceName+"-headless"),
			r.Logging.Spec.ControlNamespace),
		Port: 24240,
	}
}
