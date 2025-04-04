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
	"os/signal"
	"runtime/coverage"
	"strings"
	"syscall"
	"time"

	"emperror.dev/errors"
	prometheusOperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/spf13/cast"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	telemetryv1alpha1 "github.com/kube-logging/telemetry-controller/api/telemetry/v1alpha1"

	extensionsControllers "github.com/kube-logging/logging-operator/controllers/extensions"
	loggingControllers "github.com/kube-logging/logging-operator/controllers/logging"
	extensionsv1alpha1 "github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/v1alpha1"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
	loggingv1alpha1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1alpha1"
	loggingv1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
	"github.com/kube-logging/logging-operator/pkg/webhook/podhandler"
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
	_ = telemetryv1alpha1.AddToScheme(scheme)
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var verboseLogging bool
	var loggingOutputFormat string
	var enableprofile bool
	var namespace string
	var loggingRef string
	var watchLabeledChildren bool
	var watchLabeledSecrets bool
	var finalizerCleanup bool
	var enableTelemetryControllerRoute bool
	var klogLevel int
	var syncPeriod string

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&verboseLogging, "verbose", false, "Enable verbose logging")
	flag.StringVar(&loggingOutputFormat, "output-format", "", "Logging output format (json, console)")
	flag.IntVar(&klogLevel, "klogLevel", 0, "Global log level for klog (0-9)")
	flag.BoolVar(&enableprofile, "pprof", false, "Enable pprof")
	flag.StringVar(&namespace, "watch-namespace", "", "Namespace to filter the list of watched objects")
	flag.StringVar(&loggingRef, "watch-logging-name", "", "Logging resource name to optionally filter the list of watched objects based on which logging they belong to by checking the app.kubernetes.io/managed-by label")
	flag.BoolVar(&watchLabeledChildren, "watch-labeled-children", false, "Only watch child resources with Logging operator's name label selector: app.kubernetes.io/name: fluentd|fluentbit|syslog-ng")
	flag.BoolVar(&watchLabeledSecrets, "watch-labeled-secrets", false, "Only watch secrets with the following label selector: logging.banzaicloud.io/watch: enabled")
	flag.BoolVar(&finalizerCleanup, "finalizer-cleanup", false, "Remove finalizers from Logging resources during operator shutdown, useful for Helm uninstallation")
	flag.BoolVar(&enableTelemetryControllerRoute, "enable-telemetry-controller-route", false, "Enable the Telemetry Controller route for Logging resources")
	flag.StringVar(&syncPeriod, "sync-period", "", "SyncPeriod determines the minimum frequency at which watched resources are reconciled. Defaults to 10 hours. Parsed using time.ParseDuration.")
	flag.Parse()

	ctx := context.Background()

	zapLogger := zap.New(func(o *zap.Options) {
		o.Development = verboseLogging

		switch loggingOutputFormat {
		case "json":
			encoder := zap.JSONEncoder()
			encoder(o)
		case "console":
			encoder := zap.ConsoleEncoder()
			encoder(o)
		case "":
			break
		default:
			fmt.Printf("invalid encoder value \"%s\"", loggingOutputFormat)
			os.Exit(1)
		}
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
		Scheme:           scheme,
		Metrics:          metricsserver.Options{BindAddress: metricsAddr},
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "logging-operator." + loggingv1beta1.GroupVersion.Group,
	}

	if os.Getenv("ENABLE_WEBHOOKS") == "true" {
		webhookServerOptions := webhook.Options{
			Port:    config.TailerWebhook.ServerPort,
			CertDir: config.TailerWebhook.CertDir,
		}
		if port, ok := os.LookupEnv("WEBHOOK_PORT"); ok {
			webhookServerOptions.Port = cast.ToInt(port)
		}
		webhookServer := webhook.NewServer(webhookServerOptions)
		mgrOptions.WebhookServer = webhookServer
	}

	customMgrOptions, err := setupCustomCache(&mgrOptions, syncPeriod, namespace, loggingRef, watchLabeledChildren)
	if err != nil {
		setupLog.Error(err, "unable to set up custom cache settings")
		os.Exit(1)
	}
	if watchLabeledSecrets {
		if customMgrOptions.Cache.ByObject == nil {
			customMgrOptions.Cache.ByObject = make(map[client.Object]cache.ByObject)
		}
		customMgrOptions.Cache.ByObject[&corev1.Secret{}] = cache.ByObject{
			Label: labels.Set{"logging.banzaicloud.io/watch": "enabled"}.AsSelector(),
		}
	}

	if enableprofile {
		setupLog.Info("enabling pprof")
		pprofxIndexPath := "/debug/pprof"
		customMgrOptions.Metrics.ExtraHandlers = map[string]http.Handler{
			pprofxIndexPath: http.HandlerFunc(pprof.Index),
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), *customMgrOptions)

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := detectContainerRuntime(ctx, mgr.GetAPIReader()); err != nil {
		setupLog.Error(err, "failed to detect container runtime")
		os.Exit(1)
	}

	loggingReconciler := loggingControllers.NewLoggingReconciler(mgr.GetClient(), mgr.GetEventRecorderFor("logging-operator"), ctrl.Log.WithName("logging"))

	if err := (&extensionsControllers.EventTailerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("event-tailer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "EventTailer")
		os.Exit(1)
	}
	if err := (&extensionsControllers.HostTailerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("host-tailer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HostTailer")
		os.Exit(1)
	}

	if err := loggingControllers.SetupLoggingWithManager(mgr, ctrl.Log.WithName("manager")).Complete(loggingReconciler); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Logging")
		os.Exit(1)
	}

	if err := loggingControllers.SetupLoggingRouteWithManager(mgr, ctrl.Log.WithName("logging-route")); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "LoggingRoute")
		os.Exit(1)
	}

	if enableTelemetryControllerRoute {
		if err := loggingControllers.SetupTelemetryControllerWithManager(mgr, ctrl.Log.WithName("telemetry-controller")); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "TelemetryController")
			os.Exit(1)
		}
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

		setupLog.Info("Registering webhooks...")
		webhookHandler := podhandler.NewPodHandler(ctrl.Log.WithName("webhook-tailer"))
		webhookHandler.Decoder = admission.NewDecoder(mgr.GetScheme())
		webhookServer.Register(config.TailerWebhook.ServerPath, &webhook.Admission{Handler: webhookHandler})
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")

	if err := mgr.Start(setupSignalHandler(mgr, finalizerCleanup)); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// Extends sigs.k8s.io/controller-runtime@v0.17.2/pkg/manager/signals/signal.go with
// SIGUSR1 handler for saving test coverage files
var onlyOneSignalHandler = make(chan struct{})

func setupSignalHandler(mgr ctrl.Manager, finalizerCleanup bool) context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()

		// Due to the way Helm handles uninstallation,
		// the operator might be terminated before the finalizers are removed.
		if finalizerCleanup {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cleanupCancel()

			cleanupFinalizers(cleanupCtx, mgr.GetClient())
		}

		os.Exit(1) // second signal. Exit directly.
	}()

	coverDir, exists := os.LookupEnv("GOCOVERDIR")
	if !exists {
		return ctx
	}
	coverChan := make(chan os.Signal, 1)
	signal.Notify(coverChan, syscall.SIGUSR1)
	go func() {
		for {
			<-coverChan
			if err := coverage.WriteCountersDir(coverDir); err != nil {
				setupLog.Error(err, "Could not write coverage profile data files to the directory")
				os.Exit(1)
			}
			if err := coverage.ClearCounters(); err != nil {
				setupLog.Error(err, "Could not reset coverage counter variables")
				os.Exit(1)
			}
		}
	}()

	return ctx
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
	} else {
		setupLog.Info("Unable to detect container runtime, keeping default value", "runtime", types.ContainerRuntime)
	}

	return nil
}

func setupCustomCache(mgrOptions *ctrl.Options, syncPeriod string, namespace string, loggingRef string, watchLabeledChildren bool) (*ctrl.Options, error) {
	if syncPeriod != "" {
		duration, err := time.ParseDuration(syncPeriod)
		if err != nil {
			return mgrOptions, err
		}
		mgrOptions.Cache.SyncPeriod = &duration
	}

	if namespace != "" || loggingRef != "" || watchLabeledChildren {
		var namespaceSelector fields.Selector
		var labelSelector labels.Selector
		if namespace != "" {
			namespaceSelector = fields.Set{"metadata.namespace": namespace}.AsSelector()
		}
		if loggingRef != "" {
			labelSelector = labels.Set{"app.kubernetes.io/managed-by": loggingRef}.AsSelector()
		}
		if watchLabeledChildren {
			if labelSelector == nil {
				labelSelector = labels.NewSelector()
			}
			// It would be much better to watch for a common label, but we don't have that yet.
			// Adding a new label would recreate statefulsets and daemonsets which would be undesirable.
			// Let's see how this works in the wild. We can optimize in a subsequent iteration.
			req, err := labels.NewRequirement("app.kubernetes.io/name", selection.In, []string{"fluentd", "syslog-ng", "fluentbit"})
			if err != nil {
				return nil, err
			}
			labelSelector = labelSelector.Add(*req)
		}

		if mgrOptions.Cache.ByObject == nil {
			mgrOptions.Cache.ByObject = make(map[client.Object]cache.ByObject)
		}
		objectsToWatch := map[client.Object]cache.ByObject{
			&corev1.Pod{}:                   {Field: namespaceSelector, Label: labelSelector},
			&batchv1.Job{}:                  {Field: namespaceSelector, Label: labelSelector},
			&corev1.Service{}:               {Field: namespaceSelector, Label: labelSelector},
			&corev1.Secret{}:                {Field: namespaceSelector, Label: labelSelector},
			&rbacv1.Role{}:                  {Field: namespaceSelector, Label: labelSelector},
			&rbacv1.ClusterRole{}:           {Label: labelSelector},
			&rbacv1.RoleBinding{}:           {Field: namespaceSelector, Label: labelSelector},
			&rbacv1.ClusterRoleBinding{}:    {Label: labelSelector},
			&corev1.ServiceAccount{}:        {Field: namespaceSelector, Label: labelSelector},
			&appsv1.DaemonSet{}:             {Field: namespaceSelector, Label: labelSelector},
			&appsv1.StatefulSet{}:           {Field: namespaceSelector, Label: labelSelector},
			&appsv1.Deployment{}:            {Field: namespaceSelector, Label: labelSelector},
			&corev1.PersistentVolumeClaim{}: {Field: namespaceSelector, Label: labelSelector},
			&corev1.ConfigMap{}:             {Field: namespaceSelector, Label: labelSelector},
		}

		for obj, config := range objectsToWatch {
			mgrOptions.Cache.ByObject[obj] = config
		}
	}

	return mgrOptions, nil
}

func cleanupFinalizers(ctx context.Context, client client.Client) {
	log := ctrl.Log.WithName("finalizer-cleanup")
	log.Info("Removing finalizers during operator shutdown")

	// List all Logging resources
	loggingList := &loggingv1beta1.LoggingList{}
	if err := client.List(ctx, loggingList); err != nil {
		log.Error(err, "Failed to list Logging resources")
		return
	}

	finalizers := []string{
		loggingControllers.FluentdConfigFinalizer,
		loggingControllers.SyslogNGConfigFinalizer,
		loggingControllers.TelemetryControllerFinalizer,
	}
	for _, logging := range loggingList.Items {
		for _, finalizer := range finalizers {
			if controllerutil.ContainsFinalizer(&logging, finalizer) {
				log.Info(fmt.Sprintf("Removing finalizer: %s from: %s during operator shutdown",
					finalizer,
					logging.Name))

				controllerutil.RemoveFinalizer(&logging, finalizer)
				if err := client.Update(ctx, &logging); err != nil {
					log.Error(err, fmt.Sprintf("Failed to remove finalizer: %s from: %s during operator shutdown",
						finalizer,
						logging.Name))
					// continue trying to remove finalizers for other resources
					continue
				}
			}
		}
	}
}
