package config

import (
    "github.com/spf13/viper"
    "github.com/banzaicloud/logging-operator/cmd/logging-operator/fluentbit"
    "github.com/sirupsen/logrus"
    "github.com/banzaicloud/logging-operator/cmd/logging-operator/fluentd"
)
//Initialize the configuration

func init() {
    logrus.Info("Initializing configuration")
    viper.AddConfigPath("/logging-operator/config/")
    viper.SetConfigName("config")
}

func ConfigureOperator() {
   if viper.GetBool("fluent-bit.enabled") {
        fluentbit.InitFluentBit()
   }
   if viper.GetBool("fluentd.enabled") {
       fluentd.InitFluentd()
   }
}