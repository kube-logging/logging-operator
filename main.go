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
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"emperror.dev/errors"
	extensionsControllers "github.com/banzaicloud/logging-operator/controllers/extensions"
	loggingControllers "github.com/banzaicloud/logging-operator/controllers/logging"
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	extensionsv1alpha1 "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/api/v1alpha1"
	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/extensionsconfig"
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1alpha1"
	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/types"
	"github.com/banzaicloud/logging-operator/pkg/webhook/podhandler"
	prometheusOperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/spf13/cast"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
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
	_ = extensionsv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
	_ = prometheusOperator.AddToScheme(scheme)
	_ = apiextensions.AddToScheme(scheme)
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var verboseLogging bool
	var enableprofile bool
	var namespace string
	var loggingRef string
	var klogLevel int

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&verboseLogging, "verbose", false, "Enable verbose logging")
	flag.IntVar(&klogLevel, "klogLevel", 0, "Global log level for klog (0-9)")
	flag.BoolVar(&enableprofile, "pprof", false, "Enable pprof")
	flag.StringVar(&namespace, "watch-namespace", "", "Namespace to filter the list of watched objects")
	flag.StringVar(&loggingRef, "watch-logging-name", "", "Logging resource name to optionally filter the list of watched objects based on which logging they belong to by checking the app.kubernetes.io/managed-by label")
	flag.Parse()

	ctx := context.Background()

	zapLogger := zap.New(func(o *zap.Options) {
		o.Development = verboseLogging
	})
	ctrl.SetLogger(zapLogger)

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	err := klogFlags.Set("v", cast.ToString(klogLevel))
	if err != nil {
		fmt.Printf("%s - failed to set log level for klog, moving on.\n", err)
	}
	klog.SetLogger(zapLogger)

	mgrOptions := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "logging-operator." + loggingv1beta1.GroupVersion.Group,
		MapperProvider:     k8sutil.NewCached,
		Port:               9443,
	}

	customMgrOptions, err := setupCustomCache(&mgrOptions, namespace, loggingRef)
	if err != nil {
		setupLog.Error(err, "unable to set up custom cache settings")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), *customMgrOptions)

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

	loggingReconciler := loggingControllers.NewLoggingReconciler(mgr.GetClient(), ctrl.Log.WithName("controllers").WithName("Logging"))

	if err := (&extensionsControllers.EventTailerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("EventTailer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "EventTailer")
		os.Exit(1)
	}
	if err := (&extensionsControllers.HostTailerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("HostTailer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HostTailer")
		os.Exit(1)
	}

	if err := loggingControllers.SetupLoggingWithManager(mgr, ctrl.Log.WithName("manager")).Complete(loggingReconciler); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Logging")
		os.Exit(1)
	}

	if os.Getenv("ENABLE_WEBHOOKS") == "true" {
		if err := loggingv1beta1.SetupWebhookWithManager(mgr, loggingv1beta1.APITypes()...); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "v1beta1.logging")
			os.Exit(1)
		}
		if err := loggingv1beta1.SetupWebhookWithManager(mgr, loggingv1alpha1.APITypes()...); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "v1alpha1.logging")
			os.Exit(1)
		}

		// Webhook server registration
		setupLog.Info("Setting up webhook server...")
		webhookServer := mgr.GetWebhookServer()
		if config.TailerWebhook.ServerPort != 0 {
			webhookServer.Port = config.TailerWebhook.ServerPort
		}
		if config.TailerWebhook.CertDir != "" {
			webhookServer.CertDir = config.TailerWebhook.CertDir
		}

		setupLog.Info("Registering webhooks...")
		webhookServer.Register(config.TailerWebhook.ServerPath, &webhook.Admission{Handler: podhandler.NewPodHandler(mgr.GetClient())})

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

func setupCustomCache(mgrOptions *ctrl.Options, namespace string, loggingRef string) (*ctrl.Options, error) {
	if namespace == "" && loggingRef == "" {
		return mgrOptions, nil
	}

	var namespaceSelector fields.Selector
	var labelSelector labels.Selector
	if namespace != "" {
		namespaceSelector = fields.Set{"metadata.namespace": namespace}.AsSelector()
	}
	if loggingRef != "" {
		labelSelector = labels.Set{"app.kubernetes.io/managed-by": loggingRef}.AsSelector()
	}

	selectorsByObject := cache.SelectorsByObject{
		&corev1.Pod{}: {
			Field: namespaceSelector,
			Label: labelSelector,
		},
		&appsv1.DaemonSet{}: {
			Field: namespaceSelector,
			Label: labelSelector,
		},
		&appsv1.StatefulSet{}: {
			Field: namespaceSelector,
			Label: labelSelector,
		},
		&appsv1.Deployment{}: {
			Field: namespaceSelector,
			Label: labelSelector,
		},
		&corev1.PersistentVolumeClaim{}: {
			Field: namespaceSelector,
			Label: labelSelector,
		},
	}

	mgrOptions.NewCache = cache.BuilderWithOptions(cache.Options{SelectorsByObject: selectorsByObject})

	return mgrOptions, nil
}
