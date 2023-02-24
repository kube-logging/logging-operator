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

package eventtailer

import (
	"encoding/json"
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Config for configmap json serialization
type Config struct {
	Sink                            string `json:"sink"`
	LastResourceVersionPositionPath string `json:"lastResourceVersionPositionPath"`
}

func (e *EventTailer) positionFile() string {
	return fmt.Sprintf("%s/%s", config.Global.FluentBitPosFilePath, config.EventTailer.TailerAffix)
}

func (e *EventTailer) makeJSONString() (string, error) {
	c := Config{
		Sink:                            "stdout",
		LastResourceVersionPositionPath: e.positionFile(),
	}

	config, err := json.Marshal(c)

	return string(config), err
}

// ConfigMap resource for reconciler
func (e *EventTailer) ConfigMap() (runtime.Object, reconciler.DesiredState, error) {
	conf, err := e.makeJSONString()
	configMap := corev1.ConfigMap{
		ObjectMeta: e.objectMeta(),
		Data: map[string]string{
			config.EventTailer.ConfigurationFileName: conf,
		},
	}
	return &configMap, reconciler.StatePresent, err
}
