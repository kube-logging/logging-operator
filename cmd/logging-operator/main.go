package main

import (
	"context"
	"runtime"

	stub "github.com/banzaicloud/logging-operator/pkg/stub"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func getConfiguration() *operatorConfig {
	return &operatorConfig{}
}

func main() {
	printVersion()

	resource := "logging.banzaicloud.com/v1alpha1"
	kind := "LoggingOperator"
	//namespace, err := k8sutil.GetWatchNamespace()
	//if err != nil {
	//	logrus.Fatalf("Failed to get watch namespace: %v", err)
	//}
	logrus.Info("Deploy fluent-bit")
	initFluentBit()
	resyncPeriod := 5
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, "", resyncPeriod)
	sdk.Watch(resource, kind, "", resyncPeriod)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}

type operatorConfig struct {
}

func initFluentBit() {
	cfg := &fluentBitDeploymentConfig{
		Namespace: "default",
	}
	sdk.Create(newServiceAccount(cfg))
	sdk.Create(newClusterRole(cfg))
	sdk.Create(newClusterRoleBinding(cfg))
	cfgMap, _ := newFluentBitConifg(cfg)
	sdk.Create(cfgMap)
	sdk.Create(newFluentBitDaemonSet(cfg))
}

func initFluentd() {
	sdk.Create(newFluentdConfigmap())
	sdk.Create(newFluentdPVC())
	sdk.Create(newFluentdDeployment())
	// Create fluentd services
	// Possible options
	//  replica: x
	//  tag_rewrite config: ? it should be possible to give labels
	//  input port
	//  TLS?
	//  monitoring
	//    enabled:
	//    port:
	//    path:
}
