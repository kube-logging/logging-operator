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
        operatorNamespace = "OPERATOR_NAMESPACE"
        resource = "logging.banzaicloud.com/v1alpha1"
        kind = "LoggingOperator"
    )

    ns := os.Getenv(operatorNamespace)
    printVersion(ns)

    logrus.Info("Deploy fluent-bit")
    initFluentBit()
    logrus.Info("Deploy fluentd")
    initFluentd()
    resyncPeriod := 5
    logrus.Infof("Watching %s, %s, %s, %d", resource, kind, ns, resyncPeriod)
    sdk.Watch(resource, kind, ns, resyncPeriod)
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
    cfgMap, _ := newFluentBitConfig(cfg)
    sdk.Create(cfgMap)
    sdk.Create(newFluentBitDaemonSet(cfg))
}

func initFluentd() {
    sdk.Create(newFluentdConfigmap())
    sdk.Create(newFluentdPVC())
    sdk.Create(newFluentdDeployment())
    sdk.Create(newFluentdService())
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
