package main

import (
	"context"
	"runtime"

	"github.com/banzaicloud/logging-operator/pkg/stub"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
	"os"
)

func printVersion(namespace string) {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
	logrus.Infof("Operator namespace: %s", namespace)
}

func getConfiguration() *operatorConfig {
	return &operatorConfig{}
}

func main() {

	const (
		operatorNamespace =  "OPERATOR_NAMESPACE"
		operatorResource  =  "logging.banzaicloud.com/v1alpha1"
		configMap         =  "ConfigMap"
		kind              =  "LoggingOperator"
	)

	ns := os.Getenv(operatorNamespace)
	printVersion(ns)
	//logrus.Info("Deploy fluentd")
	//fluentd.InitFluentd()
	resyncPeriod := 0
	logrus.Infof("Watching %s, %s, %s, %d", operatorResource, kind, ns, resyncPeriod)
	sdk.Watch(operatorResource, kind, ns, resyncPeriod)
	logrus.Infof("Watching %s, %s, %s", "v1", configMap, ns)
	sdk.Watch("v1", configMap, ns, resyncPeriod)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}

type operatorConfig struct {
}
