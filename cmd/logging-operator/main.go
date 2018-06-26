package main

import (
	"context"
	"runtime"

	stub "github.com/banzaicloud/logging-operator/pkg/stub"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	k8sutil "github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()

	resource := "logging.banzaicloud.com/v1alpha1"
	kind := "LoggingOperator"
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("Failed to get watch namespace: %v", err)
	}
	resyncPeriod := 5
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)
	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}

func initFluentBit() {
	// Create fluntBit daemonset
	// Basic configuration and tagging
	// TLS?
	// Output?
}

func initFluentd() {
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