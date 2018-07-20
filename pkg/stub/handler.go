package stub

import (
	"context"

	"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
    "k8s.io/api/core/v1"
    "github.com/banzaicloud/logging-operator/cmd/logging-operator/config"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) (err error) {
	switch o := event.Object.(type) {
	case *v1alpha1.LoggingOperator:
		logrus.Infof("New CRD arrived %#v", o)
    case *v1.ConfigMap:
        // Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
        // All secondary resources must have the CR set as their OwnerReference for this to be the case
        if event.Deleted {
            logrus.Info("Inside deleted event")
            return nil
        }
        if o.Labels["app"] == "fluentd" {
            logrus.Info("Fluentd config modified")
        }
        if o.Labels["app"] == "fluent-bit" {
            logrus.Info("Fluent-bit config modified")
        }
        if o.Labels["app"] == "logging-operator" {
            logrus.Info("Logging operator config modified")
            config.ConfigureOperator()
        }
	}
	return
}

//
//generateFluentdConfig(*v1alpha1.LoggingOperator) {
//
//}

type loggingOperatorCRD struct {
	Input  inputFluentd
	Filter filterFluentd
}

type inputFluentd struct {
	Label string
	Value string
}

type filterFluentd struct {
	Name       string
	Format     string
	TimeFormat string
}

type outputFluentd struct {
	S3 outputS3
}

type outputS3 struct {
	Parameters Parameters
}

type Parameters struct {
	Name      string
	ValueFrom ValueFrom
	Value     string
}

type ValueFrom struct {
	SecretKeyRef kubernetesSecret
}

type kubernetesSecret struct {
	Name string
	Key  string
}
