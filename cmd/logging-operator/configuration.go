package main

import (
    "github.com/spf13/viper"
    "github.com/banzaicloud/logging-operator/cmd/logging-operator/fluentbit"
    "github.com/sirupsen/logrus"
    "github.com/banzaicloud/logging-operator/cmd/logging-operator/fluentd"
    "github.com/fsnotify/fsnotify"
    "path/filepath"
    "fmt"
)
//Initialize the configuration
const configFile = "/logging-operator/config/config.toml"

func init() {
    logrus.Info("Initializing configuration")
    go handleConfigChanges()
}

func handleConfigChanges() {
    c := make(chan fsnotify.Event, 1)
    viper.SetConfigFile(configFile)
    go func() {
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
            logrus.Fatal(err)
        }
        defer watcher.Close()

        // we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
        configFile := filepath.Clean(configFile)
        configDir, _ := filepath.Split(configFile)

        done := make(chan bool)
        go func() {
            for {
                select {
                case event := <- watcher.Events:
                    // we only care about the config file or the ConfigMap directory (if in Kubernetes)
                    if filepath.Clean(event.Name) == configFile || filepath.Base(event.Name) == "..data" {
                        if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
                            err := viper.ReadInConfig()
                            if err != nil {
                                logrus.Println("error:", err)
                            }
                            c <- event
                        }
                    }
                case err := <-watcher.Errors:
                    logrus.Println("error:", err)
                }
            }
        }()

        watcher.Add(configDir)
        <-done
    }()
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Sprintf("Error during reading config file : %s", err))
    }
    c <- fsnotify.Event{Name: "Initial", Op: fsnotify.Create}

    for e := range c {
        logrus.Infoln("New config file change", e.String())
        configureOperator()
    }
}

func configureOperator() {
   if viper.GetBool("fluent-bit.enabled") {
       logrus.Info("Trying to init fluent-bit")
        fluentbit.InitFluentBit()
   } else if !viper.GetBool("fluent-bit.enabled"){
       logrus.Info("Deleting fluent-bit DaemonSet...")
       fluentbit.DeleteFluentBit()
   }
   if viper.GetBool("fluentd.enabled") {
       logrus.Info("Trying to init fluentd")
       fluentd.InitFluentd()
   } else if !viper.GetBool("fluentd.enabled") {

   }
}