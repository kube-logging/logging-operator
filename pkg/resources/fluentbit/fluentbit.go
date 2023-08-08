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

package fluentbit

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/resources/loggingdataprovider"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	defaultServiceAccountName      = "fluentbit"
	clusterRoleBindingName         = "fluentbit"
	clusterRoleName                = "fluentbit"
	fluentBitSecretConfigName      = "fluentbit"
	fluentbitDaemonSetName         = "fluentbit"
	fluentbitPodSecurityPolicyName = "fluentbit"
	fluentbitServiceName           = "fluentbit"
	containerName                  = "fluent-bit"
	defaultBufferVolumeMetricsPort = 9200
)

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func (r *Reconciler) getFluentBitLabels() map[string]string {
	return util.MergeLabels(
		r.fluentbitSpec.Labels,
		map[string]string{
			"app.kubernetes.io/instance": r.nameProvider.Name(),
			"app.kubernetes.io/name":     "fluentbit",
		},
		generateLoggingRefLabels(r.Logging.GetName()))
}

func (r *Reconciler) getServiceAccount() string {
	if r.fluentbitSpec.Security.ServiceAccount != "" {
		return r.fluentbitSpec.Security.ServiceAccount
	}
	return r.nameProvider.ComponentName(defaultServiceAccountName)
}

type NameProvider interface {
	// ComponentName provides a qualified name using (Name + "-" + name)
	ComponentName(name string) string
	// Name returns the name of the resource, that is owning fluentbit
	// It is Logging.Name for legacy but the resource's name for FluentbitAgent
	Name() string
	OwnerRef() v1.OwnerReference
}

type DesiredObject struct {
	Object runtime.Object
	State  reconciler.DesiredState
}

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	resourceReconciler  *reconciler.GenericResourceReconciler
	logger              logr.Logger
	Logging             *v1beta1.Logging
	configs             map[string][]byte
	fluentbitSpec       *v1beta1.FluentbitSpec
	loggingDataProvider loggingdataprovider.LoggingDataProvider
	nameProvider        NameProvider
}

// NewReconciler creates a new FluentbitAgent reconciler
func New(client client.Client,
	logger logr.Logger,
	logging *v1beta1.Logging,
	opts reconciler.ReconcilerOpts,
	fluentbitSpec *v1beta1.FluentbitSpec,
	loggingDataProvider loggingdataprovider.LoggingDataProvider,
	nameProvider NameProvider) *Reconciler {
	return &Reconciler{
		Logging:             logging,
		logger:              logger,
		resourceReconciler:  reconciler.NewGenericReconciler(client, logger.WithName("reconciler"), opts),
		fluentbitSpec:       fluentbitSpec,
		loggingDataProvider: loggingDataProvider,
		nameProvider:        nameProvider,
	}
}

// Reconcile reconciles the fluentBit resource
func (r *Reconciler) Reconcile(ctx context.Context) (*reconcile.Result, error) {
	if err := v1beta1.FluentBitDefaults(r.fluentbitSpec); err != nil {
		return nil, err
	}

	objects := []resources.Resource{
		r.serviceAccount,
		r.clusterRole,
		r.clusterRoleBinding,
		r.configSecret,
		r.daemonSet,
		r.serviceMetrics,
		r.serviceBufferMetrics,
	}
	if resources.PSPEnabled {
		objects = append(objects, r.clusterPodSecurityPolicy, r.pspClusterRole, r.pspClusterRoleBinding)
	}
	if resources.IsSupported(ctx, resources.ServiceMonitorKey) {
		objects = append(objects, r.monitorServiceMetrics, r.monitorBufferServiceMetrics)
	}
	if resources.IsSupported(ctx, resources.PrometheusRuleKey) {
		objects = append(objects, r.prometheusRules, r.bufferVolumePrometheusRules)
	}
	for _, factory := range objects {
		o, state, err := factory()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", factory)
		}
		result, err := r.resourceReconciler.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapWithDetails(err,
				"failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

type FluentbitNameProvider struct {
	logging   *v1beta1.Logging
	fluentbit *v1beta1.FluentbitAgent
}

func (l *FluentbitNameProvider) ComponentName(name string) string {
	if l.logging != nil {
		return l.logging.QualifiedName(name)
	}
	return fmt.Sprintf("%s-%s", l.fluentbit.Name, name)
}

func (l *FluentbitNameProvider) Name() string {
	if l.logging != nil {
		return l.logging.Name
	}
	return l.fluentbit.Name
}

func (l *FluentbitNameProvider) OwnerRef() v1.OwnerReference {
	if l.logging != nil {
		return v1.OwnerReference{
			APIVersion: l.logging.APIVersion,
			Kind:       l.logging.Kind,
			Name:       l.logging.Name,
			UID:        l.logging.UID,
			Controller: util.BoolPointer(true),
		}
	}
	return v1.OwnerReference{
		APIVersion: l.fluentbit.APIVersion,
		Kind:       l.fluentbit.Kind,
		Name:       l.fluentbit.Name,
		UID:        l.fluentbit.UID,
		Controller: util.BoolPointer(true),
	}
}

func NewLegacyFluentbitNameProvider(logging *v1beta1.Logging) *FluentbitNameProvider {
	return &FluentbitNameProvider{
		logging: logging,
	}
}

func NewStandaloneFluentbitNameProvider(agent *v1beta1.FluentbitAgent) *FluentbitNameProvider {
	return &FluentbitNameProvider{
		fluentbit: agent,
	}
}

func RegisterWatches(builder *builder.Builder) *builder.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{})
}
