// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package nodeagent

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
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

type fluentbitInputConfig struct {
	Values          map[string]string
	ParserN         []string
	MultilineParser []string
}

type upstreamNode struct {
	Name string
	Host string
	Port int32
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
		ConnectTimeoutSet       bool
		ConnectTimeout          uint32
		Keepalive               bool
		KeepaliveSet            bool
		KeepaliveIdleTimeout    uint32
		KeepaliveIdleTimeoutSet bool
		KeepaliveMaxRecycle     uint32
		KeepaliveMaxRecycleSet  bool
	}
	ForwardOptions map[string]string
	Upstream       struct {
		Enabled bool
		Config  upstream
	}
}

func (n *nodeAgentInstance) configSecret() (runtime.Object, reconciler.DesiredState, error) {
	if n.nodeAgent.FluentbitSpec.CustomConfigSecret != "" {
		return &corev1.Secret{
			ObjectMeta: n.NodeAgentObjectMeta(fluentBitSecretConfigName),
		}, reconciler.StateAbsent, nil
	}
	monitor := struct {
		Enabled bool
		Port    int32
		Path    string
	}{}
	if n.nodeAgent.FluentbitSpec.Metrics != nil {
		monitor.Enabled = true
		monitor.Port = n.nodeAgent.FluentbitSpec.Metrics.Port
		monitor.Path = n.nodeAgent.FluentbitSpec.Metrics.Path
	}

	if n.nodeAgent.FluentbitSpec.InputTail.Parser == "" {
		switch types.ContainerRuntime {
		case "docker":
			n.nodeAgent.FluentbitSpec.InputTail.Parser = "docker"
		case "containerd":
			n.nodeAgent.FluentbitSpec.InputTail.Parser = "cri"
		default:
			n.nodeAgent.FluentbitSpec.InputTail.Parser = "cri"
		}
	}

	mapper := types.NewStructToStringMapper(nil)

	// FluentBit input Values
	fluentbitInput := fluentbitInputConfig{}
	inputTail := n.nodeAgent.FluentbitSpec.InputTail

	if len(inputTail.MultilineParser) > 0 {
		fluentbitInput.MultilineParser = inputTail.MultilineParser
		inputTail.MultilineParser = nil

		// If MultilineParser is set, remove other parser fields
		// See https://docs.fluentbit.io/manual/pipeline/inputs/tail#multiline-core-v1.8

		log.Log.Info("Notice: MultilineParser is enabled. Disabling other parser options")
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

	disableKubernetesFilter := n.nodeAgent.FluentbitSpec.DisableKubernetesFilter != nil && *n.nodeAgent.FluentbitSpec.DisableKubernetesFilter == true

	if !disableKubernetesFilter {
		if n.nodeAgent.FluentbitSpec.FilterKubernetes.BufferSize == "" {
			log.Log.Info("Notice: If the Buffer_Size value is empty we will set it 0. For more information: https://github.com/fluent/fluent-bit/issues/2111")
			n.nodeAgent.FluentbitSpec.FilterKubernetes.BufferSize = "0"
		} else if n.nodeAgent.FluentbitSpec.FilterKubernetes.BufferSize != "0" {
			log.Log.Info("Notice: If the kubernetes filter buffer_size parameter is underestimated it can cause log loss. For more information: https://github.com/fluent/fluent-bit/issues/2111")
		}
	}

	fluentbitKubernetesFilter, err := mapper.StringsMap(n.nodeAgent.FluentbitSpec.FilterKubernetes)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map kubernetes filter for fluentbit")
	}

	fluentbitBufferStorage, err := mapper.StringsMap(n.nodeAgent.FluentbitSpec.BufferStorage)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map buffer storage for fluentbit")
	}

	input := fluentBitConfig{
		Flush:                   n.nodeAgent.FluentbitSpec.Flush,
		Grace:                   n.nodeAgent.FluentbitSpec.Grace,
		LogLevel:                n.nodeAgent.FluentbitSpec.LogLevel,
		CoroStackSize:           n.nodeAgent.FluentbitSpec.CoroStackSize,
		Namespace:               n.logging.Spec.ControlNamespace,
		Monitor:                 monitor,
		TargetHost:              fmt.Sprintf("%s.%s.svc%s", n.FluentdQualifiedName(fluentd.ServiceName), n.logging.Spec.ControlNamespace, n.logging.ClusterDomainAsSuffix()),
		TargetPort:              fluentd.ServicePort,
		Input:                   fluentbitInput,
		DisableKubernetesFilter: disableKubernetesFilter,
		KubernetesFilter:        fluentbitKubernetesFilter,
		BufferStorage:           fluentbitBufferStorage,
	}

	if n.nodeAgent.FluentbitSpec != nil && n.nodeAgent.FluentbitSpec.TLS != nil {
		input.TLS = struct {
			Enabled   bool
			SharedKey string
		}{
			Enabled:   *n.nodeAgent.FluentbitSpec.TLS.Enabled,
			SharedKey: n.nodeAgent.FluentbitSpec.TLS.SharedKey,
		}
	}
	if n.nodeAgent.FluentbitSpec.FilterAws != nil {
		awsFilter, err := mapper.StringsMap(n.nodeAgent.FluentbitSpec.FilterAws)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map aws filter for fluentbit")
		}
		input.AwsFilter = awsFilter
	}
	if n.nodeAgent.FluentbitSpec.TargetHost != "" {
		input.TargetHost = n.nodeAgent.FluentbitSpec.TargetHost
	}
	if n.nodeAgent.FluentbitSpec.TargetPort != 0 {
		input.TargetPort = n.nodeAgent.FluentbitSpec.TargetPort
	}
	if n.nodeAgent.FluentbitSpec.ForwardOptions != nil {
		forwardOptions, err := mapper.StringsMap(n.nodeAgent.FluentbitSpec.ForwardOptions)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map forwardOptions for fluentbit")
		}
		input.ForwardOptions = forwardOptions
	}

	if n.nodeAgent.FluentbitSpec.Network != nil {
		if n.nodeAgent.FluentbitSpec.Network.ConnectTimeout != nil {
			input.Network.ConnectTimeoutSet = true
			input.Network.ConnectTimeout = *n.nodeAgent.FluentbitSpec.Network.ConnectTimeout
		}

		if n.nodeAgent.FluentbitSpec.Network.Keepalive != nil {
			input.Network.KeepaliveSet = true
			input.Network.Keepalive = *n.nodeAgent.FluentbitSpec.Network.Keepalive
		}

		if n.nodeAgent.FluentbitSpec.Network.KeepaliveIdleTimeout != nil {
			input.Network.KeepaliveIdleTimeoutSet = true
			input.Network.KeepaliveIdleTimeout = *n.nodeAgent.FluentbitSpec.Network.KeepaliveIdleTimeout
		}

		if n.nodeAgent.FluentbitSpec.Network.KeepaliveMaxRecycle != nil {
			input.Network.KeepaliveMaxRecycleSet = true
			input.Network.KeepaliveMaxRecycle = *n.nodeAgent.FluentbitSpec.Network.KeepaliveMaxRecycle
		}
	}

	if nil == n.loggingDataProvider {
		return nil, nil, errors.WrapIf(err, "nil fluent data provider")
	}

	fluentdReplicas, err := n.loggingDataProvider.GetReplicaCount(context.TODO())
	if err != nil {
		return nil, nil, errors.WrapIf(err, "getting replica count for fluentd")
	}

	if n.nodeAgent.FluentbitSpec.Network == nil && utils.PointerToInt32(fluentdReplicas) > 1 {
		input.Network.KeepaliveSet = true
		input.Network.Keepalive = true
		input.Network.KeepaliveIdleTimeoutSet = true
		input.Network.KeepaliveIdleTimeout = 30
		input.Network.KeepaliveMaxRecycleSet = true
		input.Network.KeepaliveMaxRecycle = 100
		log.Log.Info("Notice: Because the Fluentd statefulset has been scaled, we've made some changes in the fluentbit network config too. We advice to revise these default configurations.")
	}

	if utils.PointerToBool(n.nodeAgent.FluentbitSpec.EnableUpstream) {
		input.Upstream.Enabled = true
		input.Upstream.Config.Name = "fluentd-upstream"

		for i := int32(0); i < utils.PointerToInt32(fluentdReplicas); i++ {
			input.Upstream.Config.Nodes = append(input.Upstream.Config.Nodes, n.generateUpstreamNode(i))
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

	n.configs = confs

	return &corev1.Secret{
		ObjectMeta: n.NodeAgentObjectMeta(fluentBitSecretConfigName),
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

func (n *nodeAgentInstance) generateUpstreamNode(index int32) upstreamNode {
	podName := n.FluentdQualifiedName(fmt.Sprintf("%s-%d", fluentd.ComponentFluentd, index))
	return upstreamNode{
		Name: podName,
		Host: fmt.Sprintf("%s.%s.%s.svc%s",
			podName,
			n.FluentdQualifiedName(fluentd.ServiceName+"-headless"),
			n.logging.Spec.ControlNamespace,
			n.logging.ClusterDomainAsSuffix()),
		Port: n.fluentdSpec.Port,
	}
}
