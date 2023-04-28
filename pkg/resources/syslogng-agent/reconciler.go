// Copyright © 2023 Kube logging authors
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

package syslogng_agent

import (
	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"

	"github.com/kube-logging/logging-operator/pkg/resources/nodeagent"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func Reconcile(agent v1beta1.SyslogNGAgent) error {
	metricsEnabled := agent.Spec.Metrics != nil
	prometheusAnnotationsEnabled := metricsEnabled && agent.Spec.Metrics.PrometheusAnnotations
	if spec, err := nodeagent.NodeAgentSyslogNGDefaults(metricsEnabled, prometheusAnnotationsEnabled); err != nil {
		agent.Spec = *spec
		return errors.Wrap(err, "applying syslogNG defaults")
	} else {
		err = merge.Merge(&agent.Spec, spec)
		if err != nil {
			return err
		}
	}

	// 2. generate resources
	//   2.1 generate most of static resources
	//   2.2 generate config

	return nil
}
