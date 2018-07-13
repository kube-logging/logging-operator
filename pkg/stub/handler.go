package stub

import (
    "context"

    "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
    "github.com/operator-framework/operator-sdk/pkg/sdk"
    "github.com/sirupsen/logrus"
)

func NewHandler() sdk.Handler {
    return &Handler{}
}

type Handler struct {
    // Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
    switch o := event.Object.(type) {
    case *v1alpha1.LoggingOperator:
        logrus.Info("New CRD arrived %#v", o)
    }
    return nil
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
