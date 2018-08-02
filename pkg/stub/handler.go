package stub

import (
	"context"

	"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHandler creates a new Handler struct
func NewHandler() sdk.Handler {
	return &Handler{}
}

// Handler struct
type Handler struct {
	// Fill me
}

// Handle every event set up by the watcher
func (h *Handler) Handle(ctx context.Context, event sdk.Event) (err error) {
	switch o := event.Object.(type) {
	case *v1alpha1.LoggingOperator:
		if event.Deleted {
			deleteFromConfigMap(o.Spec.Input.Label["app"])
			return
		}
		logrus.Infof("New CRD arrived %#v", o)
		logrus.Info("Generating configuration.")
		name, config := generateFluentdConfig(o)
		if config != "" && name != "" {
			updateConfigMap(name, config)
		}
	}
	return
}

func deleteFromConfigMap(name string) {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-app-config",
			Namespace: "default",
		},
	}
	err := sdk.Get(configMap)
	if err != nil {
		logrus.Error(err)
	}
	if configMap.Data == nil {
		configMap.Data = map[string]string{}
	}
	delete(configMap.Data, name+".conf")
	err = sdk.Update(configMap)
	if err != nil {
		logrus.Error(err)
	}
}

func updateConfigMap(name, config string) {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-app-config",
			Namespace: "default",
		},
	}
	err := sdk.Get(configMap)
	if err != nil {
		logrus.Error(err)
	}
	if configMap.Data == nil {
		configMap.Data = map[string]string{}
	}
	configMap.Data[name+".conf"] = config
	err = sdk.Update(configMap)
	if err != nil {
		logrus.Error(err)
	}
}

//
func generateFluentdConfig(crd *v1alpha1.LoggingOperator) (string, string) {
	var finalConfig string
	// Create pattern for match
	baseMap := map[string]string{}
	baseMap["pattern"] = crd.Spec.Input.Label["app"]
	// Generate filters
	for _, filter := range crd.Spec.Filter {
		logrus.Info("Applying filter")
		values := filter.GetMap()
		values["pattern"] = crd.Spec.Input.Label["app"]
		config, err := filter.Render(values)
		if err != nil {
			logrus.Error("Error in rendering template.")
			return "", ""
		}
		finalConfig += config
	}

	// Generate output
	for _, output := range crd.Spec.Output {
		values := output.S3.GetMap()
		values["pattern"] = crd.Spec.Input.Label["app"]
		config, err := output.S3.Render(values)
		if err != nil {
			logrus.Info("Error in rendering template.")
			return "", ""
		}
		finalConfig += config
	}
	return baseMap["pattern"], finalConfig

}
