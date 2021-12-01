// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package eventtailer

import (
	"encoding/json"
	"fmt"

	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensionsconfig"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
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
