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
	"fmt"
	"text/template"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type fluentbitInputConfig struct {
	Values  map[string]string
	ParserN []string
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
	Output           map[string]string
	TargetHost       string
	TargetPort       int32
	Input            fluentbitInputConfig
	KubernetesFilter map[string]string
	AwsFilter        map[string]string
	BufferStorage    map[string]string
	Network          struct {
		ConnectTimeoutSet       bool
		ConnectTimeout          uint32
		Keepalive               bool
		KeepaliveSet            bool
		KeepaliveIdleTimeout    uint32
		KeepaliveIdleTimeoutSet bool
		KeepaliveMaxRecycle     uint32
		KeepaliveMaxRecycleSet  bool
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

	mapper := types.NewStructToStringMapper(nil)

	// FluentBit input Values
	fluentbitInput := fluentbitInputConfig{}
	inputTail := r.Logging.Spec.FluentbitSpec.InputTail
	if len(inputTail.ParserN) > 0 {
		fluentbitInput.ParserN = r.Logging.Spec.FluentbitSpec.InputTail.ParserN
		inputTail.ParserN = nil
	}
	fluentbitInputValues, err := mapper.StringsMap(inputTail)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map container tailer config for fluentbit")
	}
	fluentbitInput.Values = fluentbitInputValues

	fluentbitKubernetesFilter, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.FilterKubernetes)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map kubernetes filter for fluentbit")
	}

	fluentbitBufferStorage, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.BufferStorage)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to map buffer storage for fluentbit")
	}

	input := fluentBitConfig{
		Namespace: r.Logging.Spec.ControlNamespace,
		TLS: struct {
			Enabled   bool
			SharedKey string
		}{
			Enabled:   r.Logging.Spec.FluentbitSpec.TLS.Enabled,
			SharedKey: r.Logging.Spec.FluentbitSpec.TLS.SharedKey,
		},
		Monitor:          monitor,
		TargetHost:       fmt.Sprintf("%s.%s.svc", r.Logging.QualifiedName(fluentd.ServiceName), r.Logging.Spec.ControlNamespace),
		TargetPort:       r.Logging.Spec.FluentdSpec.Port,
		Input:            fluentbitInput,
		KubernetesFilter: fluentbitKubernetesFilter,
		BufferStorage:    fluentbitBufferStorage,
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

	if r.Logging.Spec.FluentbitSpec.Network != nil {
		if r.Logging.Spec.FluentbitSpec.Network.ConnectTimeout != nil {
			input.Network.ConnectTimeoutSet = true
			input.Network.ConnectTimeout = *r.Logging.Spec.FluentbitSpec.Network.ConnectTimeout
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
	}

	r.desiredConfig, err = generateConfig(input)
	if err != nil {
		return nil, reconciler.StatePresent, errors.WrapIf(err, "failed to generate config for fluentbit")
	}

	return &corev1.Secret{
		ObjectMeta: r.FluentbitObjectMeta(fluentBitSecretConfigName),
		Data: map[string][]byte{
			"fluent-bit.conf": []byte(r.desiredConfig),
		},
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
