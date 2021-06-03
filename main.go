// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/controllers"
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1alpha1"
	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	prometheusOperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = loggingv1beta1.AddToScheme(scheme)
	_ = loggingv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
	_ = prometheusOperator.AddToScheme(scheme)
	_ = apiextensions.AddToScheme(scheme)
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var verboseLogging bool
	var enableprofile bool

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&verboseLogging, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&enableprofile, "pprof", false, "enable pprof")
	flag.Parse()

	ctx := context.Background()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = verboseLogging
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "logging-operator." + loggingv1beta1.GroupVersion.Group,
		MapperProvider:     k8sutil.NewCached,
		Port:               9443,
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if enableprofile {
		setupLog.Info("enabling pprof")
		err = mgr.AddMetricsExtraHandler("/debug/pprof/", http.HandlerFunc(pprof.Index))
		if err != nil {
			setupLog.Error(err, "unable to attach pprof to webserver")
			os.Exit(1)
		}
	}

	if err := detectContainerRuntime(ctx, mgr.GetAPIReader()); err != nil {
		setupLog.Error(err, "failed to detect container runtime")
		os.Exit(1)
	}

	loggingReconciler := controllers.NewLoggingReconciler(mgr.GetClient(), ctrl.Log.WithName("controllers").WithName("Logging"))

	if err := controllers.SetupLoggingWithManager(mgr, ctrl.Log.WithName("manager")).Complete(loggingReconciler); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Logging")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	if os.Getenv("ENABLE_WEBHOOKS") == "true" {
		if err = loggingv1beta1.SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "logging")
			os.Exit(1)
		}
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func detectContainerRuntime(ctx context.Context, c client.Reader) error {
	var nodeList corev1.NodeList
	if err := c.List(ctx, &nodeList, client.Limit(1)); err != nil {
		return errors.WithStackIf(err)
	}

	if len(nodeList.Items) > 0 {
		runtimeWithVersion := nodeList.Items[0].Status.NodeInfo.ContainerRuntimeVersion
		runtime := strings.Split(runtimeWithVersion, "://")[0]
		setupLog.Info("Detected container runtime", "runtime", runtime)

		types.ContainerRuntime = runtime
	}

	setupLog.Info("Unable to detect container runtime, keeping default value", "runtime", types.ContainerRuntime)
	return nil
}
