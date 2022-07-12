// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	"bytes"
	"html/template"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type syslogNGConfig struct {
	LogLevel string
	Monitor  struct {
		Enabled bool
		Port    int32
		Path    string
	}
	IgnoreSameLogInterval     string
	IgnoreRepeatedLogInterval string
	Workers                   int32
	RootDir                   string
}

func generateConfig(input syslogNGConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(SyslogNGInputTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}
	return output.String(), nil
}

func (r *Reconciler) generateConfigSecret() (map[string][]byte, error) {
	input := syslogNGConfig{
		IgnoreSameLogInterval:     r.Logging.Spec.SyslogNGSpec.IgnoreSameLogInterval,
		IgnoreRepeatedLogInterval: r.Logging.Spec.SyslogNGSpec.IgnoreRepeatedLogInterval,
		RootDir:                   r.Logging.Spec.SyslogNGSpec.RootDir,
	}

	if r.Logging.Spec.SyslogNGSpec.Metrics != nil {
		input.Monitor.Enabled = true
		input.Monitor.Port = r.Logging.Spec.SyslogNGSpec.Metrics.Port
		input.Monitor.Path = r.Logging.Spec.SyslogNGSpec.Metrics.Path
	}

	input.LogLevel = r.Logging.Spec.SyslogNGSpec.LogLevel
	if input.LogLevel == "" {
		input.LogLevel = "info"
	}

	input.Workers = r.Logging.Spec.SyslogNGSpec.Workers
	if input.Workers <= 0 {
		input.Workers = 1
	}

	inputConfig, err := generateConfig(input)
	if err != nil {
		return nil, err
	}

	configs := map[string][]byte{
		"syslog-ng.conf": []byte(SyslogNGDefaultTemplate),
		"input.conf":     []byte(inputConfig),
		"devnull.conf":   []byte(SyslogNGOutputTemplate),
	}
	return configs, nil
}

func (r *Reconciler) secretConfig() (runtime.Object, reconciler.DesiredState, error) {
	configMap, err := r.generateConfigSecret()
	if err != nil {
		return nil, nil, err
	}
	//configMap["syslog-ng.conf"] = []byte(fmt.Sprintf(SyslogNGLog, r.Logging.Spec.SyslogNGSpec.SyslogNGLogDestination))
	configs := &corev1.Secret{
		ObjectMeta: r.SyslogNGObjectMeta(SecretConfigName, ComponentSyslogNG),
		Data:       configMap,
	}

	return configs, reconciler.StatePresent, nil
}
