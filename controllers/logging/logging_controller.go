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

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources"
	"github.com/kube-logging/logging-operator/pkg/resources/fluentbit"
	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/resources/loggingdataprovider"
	"github.com/kube-logging/logging-operator/pkg/resources/model"
	"github.com/kube-logging/logging-operator/pkg/resources/nodeagent"
	"github.com/kube-logging/logging-operator/pkg/resources/syslogng"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	syslogngconfig "github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"

	loggingv1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// NewLoggingReconciler returns a new LoggingReconciler instance
func NewLoggingReconciler(client client.Client, log logr.Logger) *LoggingReconciler {
	return &LoggingReconciler{
		Client: client,
		Log:    log,
	}
}

// LoggingReconciler reconciles a Logging object
type LoggingReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggings;fluentbitagents;flows;clusterflows;outputs;clusteroutputs;nodeagents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=loggings/status;fluentbitagents/status;flows/status;clusterflows/status;outputs/status;clusteroutputs/status;nodeagents/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=syslogngflows;syslogngclusterflows;syslogngoutputs;syslogngclusteroutputs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=syslogngflows/status;syslogngclusterflows/status;syslogngoutputs/status;syslogngclusteroutputs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;policy,resources=podsecuritypolicies,verbs=get;list;watch;create;update;patch;delete;use
// +kubebuilder:rbac:groups=apps,resources=statefulsets;daemonsets;replicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;persistentvolumeclaims;serviceaccounts;pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=nodes;namespaces;endpoints;nodes/proxy,verbs=get;list;watch
// +kubebuilder:rbac:groups="";events.k8s.io,resources=events,verbs=create;get;list;watch
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules;servicemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=*
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile logging resources
func (r *LoggingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("logging", req.Name)

	var logging loggingv1beta1.Logging
	if err := r.Client.Get(ctx, req.NamespacedName, &logging); err != nil {
		// If object is not found, return without error.
		// Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Client.List(ctx, &v1.ServiceMonitorList{}); err == nil {
		//nolint:staticcheck
		ctx = context.WithValue(ctx, resources.ServiceMonitorKey, true)
	} else {
		log.Info("WARNING ServiceMonitor is not supported in the cluster")
	}

	if err := r.Client.List(ctx, &v1.PrometheusRuleList{}); err == nil {
		//nolint:staticcheck
		ctx = context.WithValue(ctx, resources.PrometheusRuleKey, true)
	} else {
		log.Info("WARNING PormetheusRule is not supported in the cluster")
	}

	if err := logging.SetDefaults(); err != nil {
		return reconcile.Result{}, err
	}

	reconcilerOpts := reconciler.ReconcilerOpts{
		RecreateErrorMessageCondition:                reconciler.MatchImmutableErrorMessages,
		EnableRecreateWorkloadOnImmutableFieldChange: logging.Spec.EnableRecreateWorkloadOnImmutableFieldChange,
		EnableRecreateWorkloadOnImmutableFieldChangeHelp: "Object has to be recreated, but refusing to remove without explicitly being told so. " +
			"Use logging.spec.enableRecreateWorkloadOnImmutableFieldChange to move on but make sure to understand the consequences. " +
			"As of fluentd, to avoid data loss, make sure to use a persistent volume for buffers, which is the default, unless explicitly disabled or configured differently. " +
			"As of fluent-bit, to avoid duplicated logs, make sure to configure a hostPath volume for the positions through `logging.spec.fluentbit.spec.positiondb`. ",
	}

	loggingResourceRepo := model.NewLoggingResourceRepository(r.Client, log)

	loggingResources, err := loggingResourceRepo.LoggingResourcesFor(ctx, logging)
	if err != nil {
		return reconcile.Result{}, errors.WrapIfWithDetails(err, "failed to get logging resources", "logging", logging)
	}

	// metrics
	defer func() {
		stateMetrics, problemsMetrics := getResourceStateMetrics(log)
		// reseting the vectors should remove all orphaned metrics
		stateMetrics.Reset()
		problemsMetrics.Reset()
		for _, ob := range loggingResources.Fluentd.Flows {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.Fluentd.ClusterFlows {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.Fluentd.Outputs {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.Fluentd.ClusterOutputs {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.SyslogNG.Flows {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.SyslogNG.ClusterFlows {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.SyslogNG.Outputs {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
		for _, ob := range loggingResources.SyslogNG.ClusterOutputs {
			updateResourceStateMetrics(&ob, utils.PointerToBool(ob.Status.Active), ob.Status.ProblemsCount, stateMetrics, problemsMetrics)
		}
	}()

	reconcilers := []resources.ContextAwareComponentReconciler{
		model.NewValidationReconciler(
			r.Client,
			loggingResources,
			&secretLoaderFactory{Client: r.Client, Path: fluentd.OutputSecretPath},
			log.WithName("validation"),
		),
	}

	if logging.Spec.FluentdSpec != nil && logging.Spec.SyslogNGSpec != nil {
		return ctrl.Result{}, errors.New("fluentd and syslogNG cannot be enabled simultaneously")
	}

	var loggingDataProvider loggingdataprovider.LoggingDataProvider

	if logging.Spec.FluentdSpec != nil {
		fluentdConfig, secretList, err := r.clusterConfigurationFluentd(loggingResources)
		if err != nil {
			// TODO: move config generation into Fluentd reconciler
			reconcilers = append(reconcilers, func(ctx context.Context) (*reconcile.Result, error) {
				return &reconcile.Result{}, err
			})
		} else {
			log.V(1).Info("flow configuration", "config", fluentdConfig)

			reconcilers = append(reconcilers, fluentd.New(r.Client, r.Log, &logging, &fluentdConfig, secretList, reconcilerOpts).Reconcile)
		}
		loggingDataProvider = fluentd.NewDataProvider(r.Client, &logging)
	}

	if logging.Spec.SyslogNGSpec != nil {
		syslogNGConfig, secretList, err := r.clusterConfigurationSyslogNG(loggingResources)
		if err != nil {
			// TODO: move config generation into Syslog-NG reconciler
			reconcilers = append(reconcilers, func(ctx context.Context) (*reconcile.Result, error) {
				return &reconcile.Result{}, err
			})
		} else {
			log.V(1).Info("flow configuration", "config", syslogNGConfig)

			reconcilers = append(reconcilers, syslogng.New(r.Client, r.Log, &logging, syslogNGConfig, secretList, reconcilerOpts).Reconcile)
		}
		loggingDataProvider = syslogng.NewDataProvider(r.Client, &logging)
	}

	switch len(loggingResources.Fluentbits) {
	case 0:
		// check for legacy definition
		log.Info("WARNING fluentbit definition inside the Logging resource is deprecated and will be removed in the next major release")
		if logging.Spec.FluentbitSpec != nil {
			nameProvider := fluentbit.NewLegacyFluentbitNameProvider(&logging)
			reconcilers = append(reconcilers, fluentbit.New(
				r.Client,
				log.WithName("fluentbit-legacy"),
				&logging,
				reconcilerOpts,
				logging.Spec.FluentbitSpec,
				loggingDataProvider,
				nameProvider,
			).Reconcile)
		}
	default:
		if logging.Spec.FluentbitSpec != nil {
			return ctrl.Result{}, errors.New("fluentbit has to be removed from the logging resource before the new FluentbitAgent can be reconciled")
		}
		l := log.WithName("fluentbit")
		for _, f := range loggingResources.Fluentbits {
			f := f
			reconcilers = append(reconcilers, fluentbit.New(
				r.Client,
				l.WithValues("fluentbitagent", f.Name),
				&logging,
				reconcilerOpts,
				&f.Spec,
				loggingDataProvider,
				fluentbit.NewStandaloneFluentbitNameProvider(&f),
			).Reconcile)
		}
	}

	if len(logging.Spec.NodeAgents) > 0 || len(loggingResources.NodeAgents) > 0 {
		// load agents from standalone NodeAgent resources and additionally with inline nodeAgents from the logging resource
		// for compatibility reasons
		agents := make(map[string]loggingv1beta1.NodeAgentConfig)
		for _, a := range loggingResources.NodeAgents {
			agents[a.Name] = a.Spec.NodeAgentConfig
		}
		for _, a := range logging.Spec.NodeAgents {
			if _, exists := agents[a.Name]; !exists {
				agents[a.Name] = a.NodeAgentConfig
			} else {
				problem := fmt.Sprintf("NodeAgent resource overrides inline nodeAgent definition in logging resource %s", a.Name)
				log.Error(errors.New("nodeagent definition conflict"), problem)
			}
		}
		reconcilers = append(reconcilers, nodeagent.New(r.Client, r.Log, &logging, agents, reconcilerOpts, fluentd.NewDataProvider(r.Client, &logging)).Reconcile)
	}

	for _, rec := range reconcilers {
		result, err := rec(ctx)
		if err != nil {
			return reconcile.Result{}, err
		}
		if result != nil {
			// short circuit if requested explicitly
			return *result, err
		}
	}

	return ctrl.Result{}, nil
}

func updateResourceStateMetrics(obj client.Object, active bool, problemsCount int, statusMetric *prometheus.GaugeVec, problemsMetric *prometheus.GaugeVec) {
	statusMetric.With(prometheus.Labels{"name": obj.GetName(), "namespace": obj.GetNamespace(), "status": "active", "kind": obj.GetObjectKind().GroupVersionKind().Kind}).Set(boolToFloat64(active))
	statusMetric.With(prometheus.Labels{"name": obj.GetName(), "namespace": obj.GetNamespace(), "status": "inactive", "kind": obj.GetObjectKind().GroupVersionKind().Kind}).Set(boolToFloat64(!active))

	problemsMetric.With(prometheus.Labels{"name": obj.GetName(), "namespace": obj.GetNamespace(), "kind": obj.GetObjectKind().GroupVersionKind().Kind}).Set(float64(problemsCount))
}

func getResourceStateMetrics(logger logr.Logger) (stateMetrics *prometheus.GaugeVec, problemsMetrics *prometheus.GaugeVec) {
	var err error

	stateMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "logging_resource_state"}, []string{"name", "namespace", "status", "kind"})
	stateMetrics, err = getOrRegisterGaugeVec(metrics.Registry, stateMetrics)
	if err != nil {
		logger.Error(err, "couldn't register metrics vector for resource", "metric", stateMetrics)
	}

	problemsMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "logging_resource_problems"}, []string{"name", "namespace", "kind"})
	problemsMetrics, err = getOrRegisterGaugeVec(metrics.Registry, problemsMetrics)
	if err != nil {
		logger.Error(err, "couldn't register metrics vector for resource", "metric", problemsMetrics)
	}

	return
}

func getOrRegisterGaugeVec(reg prometheus.Registerer, gv *prometheus.GaugeVec) (*prometheus.GaugeVec, error) {
	if err := reg.Register(gv); err != nil {
		if suberr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if gv, ok := suberr.ExistingCollector.(*prometheus.GaugeVec); ok {
				return gv, nil
			} else {
				return nil, errors.WrapIfWithDetails(suberr, "already registered metric name with different type ", "metric", gv)
			}
		} else {
			return nil, err
		}
	}
	return gv, nil
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func (r *LoggingReconciler) clusterConfigurationFluentd(resources model.LoggingResources) (string, *secret.MountSecrets, error) {
	if cfg := resources.Logging.Spec.FlowConfigOverride; cfg != "" {
		return cfg, nil, nil
	}

	slf := secretLoaderFactory{
		Client: r.Client,
		Path:   fluentd.OutputSecretPath,
	}

	fluentConfig, err := model.CreateSystem(resources, &slf, r.Log)
	if err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to build model", "logging", resources.Logging)
	}

	output := &bytes.Buffer{}
	renderer := render.FluentRender{
		Out:    output,
		Indent: 2,
	}
	if err := renderer.Render(fluentConfig); err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to render fluentd config", "logging", resources.Logging)
	}

	return output.String(), &slf.Secrets, nil
}

func (r *LoggingReconciler) clusterConfigurationSyslogNG(resources model.LoggingResources) (string, *secret.MountSecrets, error) {
	if cfg := resources.Logging.Spec.FlowConfigOverride; cfg != "" {
		return cfg, nil, nil
	}

	slf := secretLoaderFactory{
		Client: r.Client,
		Path:   syslogng.OutputSecretPath,
	}

	in := syslogngconfig.Input{
		Logging:             resources.Logging,
		ClusterOutputs:      resources.SyslogNG.ClusterOutputs,
		Outputs:             resources.SyslogNG.Outputs,
		ClusterFlows:        resources.SyslogNG.ClusterFlows,
		Flows:               resources.SyslogNG.Flows,
		SecretLoaderFactory: &slf,
		SourcePort:          syslogng.ServicePort,
	}
	var b strings.Builder
	if err := syslogngconfig.RenderConfigInto(in, &b); err != nil {
		return "", nil, errors.WrapIfWithDetails(err, "failed to render syslog-ng config", "logging", resources.Logging)
	}

	return b.String(), &slf.Secrets, nil
}

type secretLoaderFactory struct {
	Client  client.Client
	Secrets secret.MountSecrets
	Path    string
}

// Deprecated: use SecretLoaderForNamespace instead
func (f *secretLoaderFactory) OutputSecretLoaderForNamespace(namespace string) secret.SecretLoader {
	return f.SecretLoaderForNamespace(namespace)
}

func (f *secretLoaderFactory) SecretLoaderForNamespace(namespace string) secret.SecretLoader {
	return secret.NewSecretLoader(f.Client, namespace, f.Path, &f.Secrets)
}

// SetupLoggingWithManager setup logging manager
func SetupLoggingWithManager(mgr ctrl.Manager, logger logr.Logger) *ctrl.Builder {
	requestMapper := handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		// get all the logging resources from the cache
		var loggingList loggingv1beta1.LoggingList
		if err := mgr.GetCache().List(ctx, &loggingList); err != nil {
			logger.Error(err, "failed to list logging resources")
			return nil
		}

		switch o := obj.(type) {
		case *loggingv1beta1.ClusterOutput:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.Output:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.Flow:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.ClusterFlow:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.SyslogNGClusterOutput:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.SyslogNGOutput:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.SyslogNGClusterFlow:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.SyslogNGFlow:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.NodeAgent:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *loggingv1beta1.FluentbitAgent:
			return reconcileRequestsForLoggingRef(loggingList.Items, o.Spec.LoggingRef)
		case *corev1.Secret:
			r := regexp.MustCompile(`^logging\.banzaicloud\.io/(.*)`)
			var requestList []reconcile.Request
			for key := range o.Annotations {
				if result := r.FindStringSubmatch(key); len(result) > 1 {
					loggingRef := result[1]
					// When loggingRef is "default" we also trigger for the empty ("") loggingRef as well, because the empty string cannot be used in the annotation, thus "default" refers to the empty case.
					if loggingRef == "default" {
						requestList = append(requestList, reconcileRequestsForLoggingRef(loggingList.Items, "")...)
					}
					requestList = append(requestList, reconcileRequestsForLoggingRef(loggingList.Items, loggingRef)...)
				}
			}
			return requestList
		}
		return nil
	})

	// Trigger reconcile for all logging resources on namespace changes that define a watchNamespaceSelector
	namespaceRequestMapper := handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		var loggingList loggingv1beta1.LoggingList
		if err := mgr.GetCache().List(ctx, &loggingList); err != nil {
			logger.Error(err, "failed to list logging resources")
			return nil
		}
		requests := make([]reconcile.Request, 0)
		for _, l := range loggingList.Items {
			if l.Spec.WatchNamespaceSelector != nil {
				requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{
					Name: l.Name,
				}})
			}
		}
		return requests
	})

	builder := ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Logging{}).
		Owns(&corev1.Pod{}).
		Watches(&corev1.Namespace{}, namespaceRequestMapper).
		Watches(&loggingv1beta1.ClusterOutput{}, requestMapper).
		Watches(&loggingv1beta1.ClusterFlow{}, requestMapper).
		Watches(&loggingv1beta1.Output{}, requestMapper).
		Watches(&loggingv1beta1.Flow{}, requestMapper).
		Watches(&loggingv1beta1.SyslogNGClusterOutput{}, requestMapper).
		Watches(&loggingv1beta1.SyslogNGClusterFlow{}, requestMapper).
		Watches(&loggingv1beta1.SyslogNGOutput{}, requestMapper).
		Watches(&loggingv1beta1.SyslogNGFlow{}, requestMapper).
		Watches(&corev1.Secret{}, requestMapper)

	// TODO remove with the next major release
	if os.Getenv("ENABLE_NODEAGENT_CRD") != "" {
		logger.Info("processing NodeAgent CRDs is explicitly disabled (enable: ENABLE_NODEAGENT_CRD=1)")
		builder.Watches(&loggingv1beta1.NodeAgent{}, requestMapper)
	}

	builder.Watches(&loggingv1beta1.FluentbitAgent{}, requestMapper)

	fluentd.RegisterWatches(builder)
	fluentbit.RegisterWatches(builder)
	nodeagent.RegisterWatches(builder)
	syslogng.RegisterWatches(builder)

	return builder
}

func reconcileRequestsForLoggingRef(loggings []loggingv1beta1.Logging, loggingRef string) (reqs []reconcile.Request) {
	for _, l := range loggings {
		if l.Spec.LoggingRef == loggingRef {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: l.Namespace, // this happens to be empty as long as Logging is cluster scoped
					Name:      l.Name,
				},
			})
		}
	}
	return
}
