package config

import (
    "github.com/spf13/viper"
    "github.com/banzaicloud/logging-operator/cmd/logging-operator/fluentbit"
    "github.com/sirupsen/logrus"
)
//Initialize the configuration
func init() {
    logrus.Info("Initializing configuration")
    viper.AddConfigPath("/logging-operator/config/")
    viper.SetConfigName("config")
    err := viper.ReadInConfig()
    if err != nil {
        logrus.Errorf("Error during reading in config file : %s", err)
    }
}

func ConfigureOperator() {
    if viper.GetBool("fluent-bit.enabled") && !fluentbit.CheckIfDeamonSetExist(){
        fluentbit.InitFluentBit()
   }
}