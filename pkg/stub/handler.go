package stub

import (
	"context"
	"github.com/banzaicloud/logging-operator/cmd/logging-operator/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/banzaicloud/logging-operator/pkg/plugins"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHandler creates a new Handler struct
func NewHandler(namepsace string) sdk.Handler {
	return &Handler{
		NameSpace: namepsace,
	}
}

// Handler struct
type Handler struct {
	NameSpace string
}

// Handle every event set up by the watcher
func (h *Handler) Handle(ctx context.Context, event sdk.Event) (err error) {
	switch o := event.Object.(type) {
	case *v1alpha1.LoggingOperator:
		if event.Deleted {
			logrus.Infof("Delete CRD: %s", o.Name)
			deleteFromConfigMap(o.Name)
			return
		}
		logrus.Infof("New CRD arrived %#v", o)
		logrus.Info("Generating configuration.")
		name, config := generateFluentdConfig(o, h.NameSpace)
		if config != "" && name != "" {
			fluentd.CreateOrUpdateAppConfig(name, config)
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

//
func generateFluentdConfig(crd *v1alpha1.LoggingOperator, namespace string) (string, string) {
	var finalConfig string
	// Generate filters
	for _, filter := range crd.Spec.Filter {
		logrus.Info("Applying filter")
		values, err := plugins.GetDefaultValues(filter.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		values["pattern"] = crd.Spec.Input.Label["app"]
		config, err := v1alpha1.RenderPlugin(filter, values, namespace)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		finalConfig += config
	}

	// Generate output
	for _, output := range crd.Spec.Output {
		values, err := plugins.GetDefaultValues(output.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		values["pattern"] = crd.Spec.Input.Label["app"]
		config, err := v1alpha1.RenderPlugin(output, values, namespace)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		finalConfig += config
	}
	return crd.Name, finalConfig

}
