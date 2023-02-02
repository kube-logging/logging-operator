// Copyright © 2020 Banzai Cloud
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

package resourcebuilder

import (
	"context"
	"fmt"
	"io/ioutil"

	"emperror.dev/errors"
	extensionsv1alpha1 "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/api/v1alpha1"
	extensionsconfig "github.com/banzaicloud/logging-operator/pkg/sdk/extensions/extensionsconfig"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/static/gen/crds"
	"github.com/banzaicloud/logging-operator/pkg/sdk/static/gen/rbac"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	admissionregistration "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	Image            = "ghcr.io/banzaicloud/logging-operator:4.0.0-rc19"
	defaultNamespace = "logging-system"
)

// +kubebuilder:object:generate=true

type ComponentConfig struct {
	types.EnabledComponent `json:",inline"`
	Namespace              string               `json:"namespace,omitempty"`
	MetaOverrides          *types.MetaBase      `json:"metaOverrides,omitempty"`
	WorkloadMetaOverrides  *types.MetaBase      `json:"workloadMetaOverrides,omitempty"`
	WorkloadOverrides      *types.PodSpecBase   `json:"workloadOverrides,omitempty"`
	ContainerOverrides     *types.ContainerBase `json:"containerOverrides,omitempty"`
	WatchNamespace         string               `json:"watchNamespace,omitempty"`
	WatchLoggingName       string               `json:"watchLoggingName,omitempty"`
	DisableWebhook         bool                 `json:"disableWebhook,omitempty"`
	Metrics                *v1beta1.Metrics     `json:"installServiceMonitor,omitempty"`
	InstallPrometheusRules bool                 `json:"-"`
}

func (c *ComponentConfig) build(parent reconciler.ResourceOwner, fn func(reconciler.ResourceOwner, ComponentConfig) (runtime.Object, reconciler.DesiredState, error)) reconciler.ResourceBuilder {
	return reconciler.ResourceBuilder(func() (runtime.Object, reconciler.DesiredState, error) {
		return fn(parent, *c)
	})
}

func ResourceBuildersWithReader(reader client.Reader) reconciler.ResourceBuilders {
	return func(parent reconciler.ResourceOwner, object interface{}) (resources []reconciler.ResourceBuilder) {
		config := &ComponentConfig{}
		if object != nil {
			config = object.(*ComponentConfig)
		}
		if config.Namespace == "" {
			config.Namespace = defaultNamespace
		}
		var modifiers []CRDModifier
		// We don't return with an absent state since we don't want them to be removed
		// however we return the CRDs without modified with webhook configuration to set the conversion strategy to default (None)
		if config.IsEnabled() && !config.DisableWebhook {
			modifiers = ConversionWebhookModifiers(parent, config)
		}
		resources = AppendCRDResourceBuilders(resources, modifiers...)
		resources = AppendOperatorResourceBuilders(resources, parent, config)

		if ok, _ := CRDExists(context.Background(), reader, "issuers.cert-manager.io"); ok {
			resources = AppendWebhookResourceBuilders(resources, parent, config)
		}

		if config.InstallPrometheusRules {
			resources = AppendPrometheusRulesResourceBuilders(resources, parent, config)
		}
		if config.Metrics != nil {
			resources = AppendServiceMonitorBuilder(resources, parent, config)
		}
		return resources
	}
}

func AppendOperatorResourceBuilders(rbs []reconciler.ResourceBuilder, parent reconciler.ResourceOwner, config *ComponentConfig) []reconciler.ResourceBuilder {
	return append(rbs,
		config.build(parent, Namespace),
		config.build(parent, Operator),
		config.build(parent, ClusterRole),
		config.build(parent, ClusterRoleBinding),
		config.build(parent, ServiceAccount),
	)
}

func AppendCRDResourceBuilders(rbs []reconciler.ResourceBuilder, modifiers ...CRDModifier) []reconciler.ResourceBuilder {
	return append(rbs,
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "loggings", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "flows", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "clusterflows", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "outputs", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "clusteroutputs", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "syslogngflows", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "syslogngclusterflows", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "syslogngoutputs", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(loggingv1beta1.GroupVersion.Group, "syslogngclusteroutputs", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(extensionsv1alpha1.GroupVersion.Group, "hosttailers", modifiers...)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(extensionsv1alpha1.GroupVersion.Group, "eventtailers", modifiers...)
		},
	)
}

func AppendWebhookResourceBuilders(rbs []reconciler.ResourceBuilder, parent reconciler.ResourceOwner, config *ComponentConfig) []reconciler.ResourceBuilder {
	return append(rbs,
		config.build(parent, MutatingWebhookConfiguration),
		config.build(parent, WebhookService),
		config.build(parent, Issuer),
		config.build(parent, Certificate),
	)
}

func AppendPrometheusRulesResourceBuilders(rbs []reconciler.ResourceBuilder, parent reconciler.ResourceOwner, config *ComponentConfig) []reconciler.ResourceBuilder {
	return append(rbs,
		func() (runtime.Object, reconciler.DesiredState, error) {
			return &monitoringv1.PrometheusRule{
				ObjectMeta: config.objectMeta(parent),
				Spec: monitoringv1.PrometheusRuleSpec{
					Groups: []monitoringv1.RuleGroup{
						{
							Name: "fluentd",
							Rules: []monitoringv1.Rule{
								{
									Record: "fluentd_buffer_size_bytes",
									Expr:   intstr.FromString(`avg (node_filesystem_size_bytes{container="buffer-metrics-sidecar",mountpoint="/buffers"}) without(container,mountpoint)`),
									Labels: map[string]string{
										"service": "fluentd",
									},
								},
								{
									Record: "fluentd_buffer_avail_bytes",
									Expr:   intstr.FromString(`avg (node_filesystem_avail_bytes{container="buffer-metrics-sidecar",mountpoint="/buffers"}) without(container,mountpoint)`),
									Labels: map[string]string{
										"service": "fluentd",
									},
								},
								{
									Record: "fluentd_buffer_used_bytes",
									Expr:   intstr.FromString(`fluentd_buffer_size_bytes - fluentd_buffer_avail_bytes`),
									Labels: map[string]string{
										"service": "fluentd",
									},
								},
								{
									Record: "fluentd_buffer_usage",
									Expr:   intstr.FromString(`fluentd_buffer_used_bytes / fluentd_buffer_size_bytes`),
									Labels: map[string]string{
										"service": "fluentd",
									},
								},
							},
						},
					},
				},
			}, reconciler.StatePresent, nil
		},
	)
}

func SetupWithBuilder(builder *builder.Builder) {
	builder.Owns(&appsv1.Deployment{})
}

func Namespace(_ reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	return &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: config.Namespace,
		},
	}, reconciler.StateCreated, nil
}

func CRD(group string, kind string, modifiers ...CRDModifier) (runtime.Object, reconciler.DesiredState, error) {
	crd := &crdv1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", kind, group),
		},
	}
	crdFile, err := crds.Root.Open(fmt.Sprintf("/%s_%s.yaml", group, kind))
	if err != nil {
		return nil, nil, errors.WrapIff(err, "failed to open %s crd", kind)
	}
	bytes, err := ioutil.ReadAll(crdFile)
	if err != nil {
		return nil, nil, errors.WrapIff(err, "failed to read %s crd", kind)
	}

	// apply modifiers
	for _, modifier := range modifiers {
		crdModified, err := modifier(crd)
		if err != nil {
			return nil, nil, err
		}
		crd = crdModified
	}

	scheme := runtime.NewScheme()
	_ = crdv1.AddToScheme(scheme)

	_, _, err = serializer.NewSerializerWithOptions(serializer.DefaultMetaFactory, scheme, scheme, serializer.SerializerOptions{
		Yaml: true,
	}).Decode(bytes, &schema.GroupVersionKind{}, crd)

	if err != nil {
		return nil, nil, errors.WrapIff(err, "failed to unmarshal %s crd", kind)
	}

	// clear the TypeMeta to avoid objectmatcher diffing on it every time,
	// because the current object coming from the API Server will not have TypeMeta set
	crd.TypeMeta.Kind = ""
	crd.TypeMeta.APIVersion = ""

	return crd, reconciler.DesiredStateHook(func(object runtime.Object) error {
		current := object.(*crdv1.CustomResourceDefinition)
		// simply copy the existing status over, so that we don't diff because of it
		crd.Status = current.Status
		return nil
	}), nil
}

func Operator(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: config.MetaOverrides.Merge(config.objectMeta(parent)),
	}
	if !config.IsEnabled() {
		return deployment, reconciler.StateAbsent, nil
	}
	var env []corev1.EnvVar
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	if !config.DisableWebhook {
		env = append(env, corev1.EnvVar{Name: "ENABLE_WEBHOOKS", Value: "true"})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      parent.GetName() + WebhookNameAffix,
			MountPath: WebhookCertDir,
		})
		volumes = append(volumes, corev1.Volume{
			Name: parent.GetName() + WebhookNameAffix,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: parent.GetName() + WebhookNameAffix,
				},
			},
		})
	}
	deployment.Spec = appsv1.DeploymentSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: config.WorkloadMetaOverrides.Merge(v1.ObjectMeta{
				Labels: config.labelSelector(parent),
			}),
			Spec: config.WorkloadOverrides.Override(corev1.PodSpec{
				ServiceAccountName: config.objectMeta(parent).Name,
				Containers: []corev1.Container{
					config.ContainerOverrides.Override(corev1.Container{
						Name:    "logging-operator",
						Image:   Image,
						Command: []string{"/manager"},
						Args:    OperatorArgs(config),
						Env:     env,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("300m"),
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("50m"),
								corev1.ResourceMemory: resource.MustParse("20Mi"),
							},
						},
						VolumeMounts: volumeMounts,
					}),
				},
				Volumes: volumes,
			}),
		},
		Selector: &v1.LabelSelector{
			MatchLabels: config.labelSelector(parent),
		},
	}
	return deployment, reconciler.StatePresent, nil
}

func ServiceAccount(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: config.MetaOverrides.Merge(config.objectMeta(parent)),
	}
	if !config.IsEnabled() {
		return sa, reconciler.StateAbsent, nil
	}
	// remove internal sa in case an externally provided service account is used
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		return sa, reconciler.StateAbsent, nil
	}
	return sa, reconciler.StatePresent, nil
}

func ClusterRoleBinding(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	rb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: config.MetaOverrides.Merge(config.clusterObjectMeta(parent)),
	}

	if !config.IsEnabled() {
		return rb, reconciler.StateAbsent, nil
	}

	sa := config.objectMeta(parent).Name
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		sa = config.WorkloadOverrides.ServiceAccountName
	}

	rb.Subjects = []rbacv1.Subject{
		{
			Kind:      rbacv1.ServiceAccountKind,
			Name:      sa,
			Namespace: config.Namespace,
		},
	}
	rb.RoleRef = rbacv1.RoleRef{
		APIGroup: rbacv1.GroupName,
		Kind:     "ClusterRole",
		Name:     config.objectMeta(parent).Name,
	}

	return rb, reconciler.StatePresent, nil
}

func ClusterRole(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	role := &rbacv1.ClusterRole{
		ObjectMeta: config.MetaOverrides.Merge(config.clusterObjectMeta(parent)),
	}
	if !config.IsEnabled() {
		return role, reconciler.StateAbsent, nil
	}
	// remove internal sa in case an externally provided service account is used
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		return role, reconciler.StateAbsent, nil
	}
	roleFile, err := rbac.Root.Open("/role.yaml")
	if err != nil {
		return nil, nil, errors.WrapIf(err, "failed to open role.yaml")
	}
	roleAsByte, err := ioutil.ReadAll(roleFile)
	if err != nil {
		return nil, nil, err
	}
	scheme := runtime.NewScheme()
	err = rbacv1.AddToScheme(scheme)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "failed to extend scheme with rbacv1 types")
	}
	_, _, err = serializer.NewSerializerWithOptions(serializer.DefaultMetaFactory, scheme, scheme, serializer.SerializerOptions{
		Yaml: true,
	}).Decode(roleAsByte, &schema.GroupVersionKind{}, role)

	// overwrite the objectmeta that has been read from file
	role.ObjectMeta = config.MetaOverrides.Merge(config.clusterObjectMeta(parent))
	role.TypeMeta.Kind = ""
	role.TypeMeta.APIVersion = ""

	return role, reconciler.StatePresent, err
}

func AppendServiceMonitorBuilder(rbs []reconciler.ResourceBuilder, parent reconciler.ResourceOwner, config *ComponentConfig) []reconciler.ResourceBuilder {
	rbs = append(rbs,
		func() (runtime.Object, reconciler.DesiredState, error) {
			return &corev1.Service{
				ObjectMeta: config.MetaOverrides.Merge(config.objectMeta(parent)),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Protocol:   corev1.ProtocolTCP,
							Name:       "http-metrics",
							Port:       8080,
							TargetPort: intstr.IntOrString{IntVal: 8080},
						},
					},
					Selector:  config.labelSelector(parent),
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "None",
				},
			}, reconciler.StatePresent, nil
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return &monitoringv1.ServiceMonitor{
				ObjectMeta: config.MetaOverrides.Merge(config.objectMeta(parent)),
				Spec: monitoringv1.ServiceMonitorSpec{
					JobLabel:        "",
					TargetLabels:    nil,
					PodTargetLabels: nil,
					Endpoints: []monitoringv1.Endpoint{{
						Port:                 "http-metrics",
						Path:                 config.Metrics.Path,
						Interval:             monitoringv1.Duration(config.Metrics.Interval),
						ScrapeTimeout:        monitoringv1.Duration(config.Metrics.Timeout),
						HonorLabels:          config.Metrics.ServiceMonitorConfig.HonorLabels,
						RelabelConfigs:       config.Metrics.ServiceMonitorConfig.Relabelings,
						MetricRelabelConfigs: config.Metrics.ServiceMonitorConfig.MetricsRelabelings,
						Scheme:               config.Metrics.ServiceMonitorConfig.Scheme,
						TLSConfig:            config.Metrics.ServiceMonitorConfig.TLSConfig,
					}},
					Selector:          v1.LabelSelector{MatchLabels: config.MetaOverrides.Merge(config.objectMeta(parent)).Labels},
					NamespaceSelector: monitoringv1.NamespaceSelector{MatchNames: []string{config.MetaOverrides.Merge(config.objectMeta(parent)).Namespace}},
					SampleLimit:       0,
				},
			}, reconciler.StatePresent, nil
		},
	)

	return rbs

}

func (c *ComponentConfig) objectMeta(parent reconciler.ResourceOwner) v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:      parent.GetName() + "-logging-operator",
		Namespace: c.Namespace,
		Labels:    c.labelSelector(parent),
	}
	return meta
}

func (c *ComponentConfig) clusterObjectMeta(parent reconciler.ResourceOwner) v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:   parent.GetName() + "-logging-operator",
		Labels: c.labelSelector(parent),
	}
	return meta
}

func (c *ComponentConfig) clusterObjectMetaExt(parent reconciler.ResourceOwner, annotations map[string]string) v1.ObjectMeta {
	meta := c.clusterObjectMeta(parent)
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string, len(annotations))
	}
	for k, v := range annotations {
		meta.Annotations[k] = v
	}
	return meta
}

func (c *ComponentConfig) labelSelector(parent reconciler.ResourceOwner) map[string]string {
	return map[string]string{
		"banzaicloud.io/operator": parent.GetName() + "-logging-operator",
	}
}

func WebhookService(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	svc := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      parent.GetName() + WebhookNameAffix,
			Namespace: config.Namespace,
			Labels:    config.labelSelector(parent),
		},
	}
	if !config.IsEnabled() {
		return svc, reconciler.StateAbsent, nil
	}

	svc.Spec = corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "logging-webhooks",
				Port:       DefaultSecurePort,
				TargetPort: intstr.IntOrString{IntVal: DefaultWebhookPort},
				Protocol:   corev1.ProtocolTCP,
			},
		},
		Selector: config.labelSelector(parent),
		Type:     corev1.ServiceTypeClusterIP,
	}

	return svc, reconciler.StateCreated, nil
}

type strIfMap = map[string]interface{}

func Issuer(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	issuer := unstructured.Unstructured{}
	issuer.SetAPIVersion("cert-manager.io/v1")
	issuer.SetKind("Issuer")
	issuer.SetName(parent.GetName() + WebhookNameAffix)
	issuer.SetNamespace(config.Namespace)
	issuer.SetLabels(config.labelSelector(parent))
	issuer.Object["spec"] = strIfMap{"selfSigned": strIfMap{}}

	if config.DisableWebhook {
		return &issuer, reconciler.StateAbsent, nil
	}
	return &issuer, reconciler.StatePresent, nil
}

func Certificate(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	cert := unstructured.Unstructured{}
	cert.SetAPIVersion("cert-manager.io/v1")
	cert.SetKind("Certificate")
	cert.SetName(parent.GetName() + WebhookNameAffix)
	cert.SetNamespace(config.Namespace)
	cert.SetLabels(config.labelSelector(parent))
	cert.Object["spec"] = strIfMap{
		"secretName": parent.GetName() + WebhookNameAffix,
		"commonName": parent.GetName() + WebhookNameAffix + "-ca",
		"isCA":       false,
		"issuerRef": strIfMap{
			"name": parent.GetName() + WebhookNameAffix,
			"kind": "Issuer",
		},
		"dnsNames": []interface{}{
			parent.GetName() + WebhookNameAffix,
			parent.GetName() + WebhookNameAffix + "." + config.Namespace,
			parent.GetName() + WebhookNameAffix + "." + config.Namespace + ".svc",
		},
	}
	if config.DisableWebhook {
		return &cert, reconciler.StateAbsent, nil
	}
	return &cert, reconciler.StatePresent, nil
}

func MutatingWebhookConfiguration(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	scope := admissionregistration.AllScopes
	failurePolicy := admissionregistration.Ignore
	sideEffects := admissionregistration.SideEffectClassNone

	scheme := runtime.NewScheme()
	_ = loggingv1beta1.AddToScheme(scheme)

	webhooks := []admissionregistration.MutatingWebhook{}
	for _, apiType := range loggingv1beta1.APITypes() {
		gvk, err := apiutil.GVKForObject(apiType, scheme)
		if err != nil {
			return nil, nil, err
		}

		if _, isDefaulter := apiType.(admission.Defaulter); !isDefaulter {
			// TODO implement proper info logging
			// log.Info(fmt.Sprintf("missing Defaulter implementation in %v, skipping mutationwebhook generation\n", gvk.Kind))
			continue
		}

		plural, _ := meta.UnsafeGuessKindToResource(gvk)

		webhook := admissionregistration.MutatingWebhook{
			Name: GVKDomainName(gvk),
			ClientConfig: admissionregistration.WebhookClientConfig{
				Service: &admissionregistration.ServiceReference{
					Name:      parent.GetName() + WebhookNameAffix,
					Namespace: config.Namespace,
					Path:      utils.StringPointer(generateMutatePath(gvk)),
				},
			},
			Rules: []admissionregistration.RuleWithOperations{
				{
					Operations: []admissionregistration.OperationType{
						admissionregistration.Create,
					},
					Rule: admissionregistration.Rule{
						APIGroups:   []string{"logging.banzaicloud.io"},
						APIVersions: []string{"v1beta1"},
						Resources:   []string{plural.Resource},
						Scope:       &scope,
					},
				},
			},
			FailurePolicy:           &failurePolicy,
			SideEffects:             &sideEffects,
			AdmissionReviewVersions: []string{"v1"},
		}

		webhooks = append(webhooks, webhook, ExtensionsMutatingWebhook(parent, config))
	}

	mutatingWebhookConfiguration := &admissionregistration.MutatingWebhookConfiguration{
		ObjectMeta: config.MetaOverrides.Merge(config.clusterObjectMetaExt(parent,
			map[string]string{"cert-manager.io/inject-ca-from": config.Namespace + "/" + parent.GetName() + WebhookNameAffix})),
		Webhooks: webhooks,
	}

	if !config.IsEnabled() || config.DisableWebhook {
		return mutatingWebhookConfiguration, reconciler.StateAbsent, nil
	}

	return mutatingWebhookConfiguration, reconciler.StatePresent, nil
}

func ConversionWebhookModifiers(parent reconciler.ResourceOwner, config *ComponentConfig) []CRDModifier {
	return []CRDModifier{
		ModifierConversionWebhook(k8stypes.NamespacedName{Name: parent.GetName() + WebhookNameAffix, Namespace: config.Namespace}),
		ModifierCAInjectAnnotation(k8stypes.NamespacedName{Name: parent.GetName() + WebhookNameAffix, Namespace: config.Namespace}),
	}
}

func OperatorArgs(config ComponentConfig) (args []string) {
	args = append(args, "--enable-leader-election")
	if config.WatchNamespace != "" {
		args = append(args, "--watch-namespace", config.WatchNamespace)
	}
	if config.WatchLoggingName != "" {
		args = append(args, "--watch-logging-name", config.WatchLoggingName)
	}
	return
}

func ExtensionsMutatingWebhook(parent reconciler.ResourceOwner, config ComponentConfig) admissionregistration.MutatingWebhook {
	scope := admissionregistration.AllScopes
	failurePolicy := admissionregistration.Ignore
	sideEffects := admissionregistration.SideEffectClassNone

	return admissionregistration.MutatingWebhook{
		Name: "logging-extensions.banzaicloud.io",
		ClientConfig: admissionregistration.WebhookClientConfig{
			Service: &admissionregistration.ServiceReference{
				Name:      parent.GetName() + WebhookNameAffix,
				Namespace: config.Namespace,
				Path:      utils.StringPointer(extensionsconfig.TailerWebhook.ServerPath),
			},
		},
		Rules: []admissionregistration.RuleWithOperations{
			{
				Operations: []admissionregistration.OperationType{
					admissionregistration.Create,
				},
				Rule: admissionregistration.Rule{
					APIGroups:   []string{corev1.SchemeGroupVersion.Group},
					APIVersions: []string{corev1.SchemeGroupVersion.Version},
					Resources:   []string{"pods"},
					Scope:       &scope,
				},
			},
		},
		FailurePolicy:           &failurePolicy,
		SideEffects:             &sideEffects,
		AdmissionReviewVersions: []string{"v1"},
	}
}

func CRDExists(ctx context.Context, reader client.Reader, crdName string) (bool, error) {
	if err := reader.Get(ctx, client.ObjectKey{Name: crdName}, new(crdv1.CustomResourceDefinition)); err != nil {
		return false, errors.WrapIff(client.IgnoreNotFound(err), "retrieving custom resource definition %q", crdName)
	}
	return true, nil
}
