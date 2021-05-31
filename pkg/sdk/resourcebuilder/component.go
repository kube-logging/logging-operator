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
	"fmt"
	"io/ioutil"

	"emperror.dev/errors"
	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/static/gen/crds"
	"github.com/banzaicloud/logging-operator/pkg/sdk/static/gen/rbac"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

const (
	Image            = "ghcr.io/banzaicloud/logging-operator:3.9.5-dev"
	defaultNamespace = "logging-system"
)

// +kubebuilder:object:generate=true

type ComponentConfig struct {
	Namespace             string               `json:"namespace,omitempty"`
	Enabled               *bool                `json:"enabled,omitempty"`
	MetaOverrides         *types.MetaBase      `json:"metaOverrides,omitempty"`
	WorkloadMetaOverrides *types.MetaBase      `json:"workloadMetaOverrides,omitempty"`
	WorkloadOverrides     *types.PodSpecBase   `json:"workloadOverrides,omitempty"`
	ContainerOverrides    *types.ContainerBase `json:"containerOverrides,omitempty"`
}

func (c *ComponentConfig) IsEnabled() bool {
	return utils.PointerToBool(c.Enabled)
}

func (c *ComponentConfig) IsSkipped() bool {
	return c.Enabled == nil
}

func (c *ComponentConfig) build(parent reconciler.ResourceOwner, fn func(reconciler.ResourceOwner, ComponentConfig) (runtime.Object, reconciler.DesiredState, error)) reconciler.ResourceBuilder {
	return reconciler.ResourceBuilder(func() (runtime.Object, reconciler.DesiredState, error) {
		return fn(parent, *c)
	})
}

func ResourceBuilders(parent reconciler.ResourceOwner, object interface{}) []reconciler.ResourceBuilder {
	config := &ComponentConfig{}
	if object != nil {
		config = object.(*ComponentConfig)
	}
	if config.Namespace == "" {
		config.Namespace = defaultNamespace
	}
	resources := []reconciler.ResourceBuilder{
		config.build(parent, Namespace),
		config.build(parent, Operator),
		config.build(parent, ClusterRole),
		config.build(parent, ClusterRoleBinding),
		config.build(parent, ServiceAccount),
		config.build(parent, Service),
		config.build(parent, Issuer),
		config.build(parent, Certificate),
	}
	// We don't return with an absent state since we don't want them to be removed
	if config.IsEnabled() {
		webhookModifiers := ConversionWebhookModifiers(parent, config)
		resources = append(resources,
			func() (runtime.Object, reconciler.DesiredState, error) {
				return CRD(config, loggingv1beta1.GroupVersion.Group, "loggings", webhookModifiers...)
			},
			func() (runtime.Object, reconciler.DesiredState, error) {
				return CRD(config, loggingv1beta1.GroupVersion.Group, "flows", webhookModifiers...)
			},
			func() (runtime.Object, reconciler.DesiredState, error) {
				return CRD(config, loggingv1beta1.GroupVersion.Group, "clusterflows", webhookModifiers...)
			},
			func() (runtime.Object, reconciler.DesiredState, error) {
				return CRD(config, loggingv1beta1.GroupVersion.Group, "outputs", webhookModifiers...)
			},
			func() (runtime.Object, reconciler.DesiredState, error) {
				return CRD(config, loggingv1beta1.GroupVersion.Group, "clusteroutputs", webhookModifiers...)
			},
		)
	}
	return resources
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

func CRD(config *ComponentConfig, group string, kind string, modifiers ...CRDModifier) (runtime.Object, reconciler.DesiredState, error) {
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
		CRDmodified, err := modifier(crd)
		if err != nil {
			return nil, nil, err
		}
		crd = CRDmodified
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
						Args:    []string{"--enable-leader-election"},
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
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      parent.GetName() + WebhookNameAffix,
								MountPath: WebhookCertDir,
							},
						},
					}),
				},
				Volumes: []corev1.Volume{
					{
						Name: parent.GetName() + WebhookNameAffix,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: parent.GetName() + WebhookNameAffix,
							},
						},
					},
				},
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

func (c *ComponentConfig) labelSelector(parent reconciler.ResourceOwner) map[string]string {
	return map[string]string{
		"banzaicloud.io/operator": parent.GetName() + "-logging-operator",
	}
}

func Service(parent reconciler.ResourceOwner, config ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
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

	spec := corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "conversion-webhook",
				Port:       DefaultSecurePort,
				TargetPort: intstr.IntOrString{IntVal: DefaultWebhookPort},
				Protocol:   corev1.ProtocolTCP,
			},
		},
		Selector: config.labelSelector(parent),
		Type:     corev1.ServiceTypeClusterIP,
	}
	svc.Spec = spec

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

	return &cert, reconciler.StatePresent, nil
}

func ConversionWebhookModifiers(parent reconciler.ResourceOwner, config *ComponentConfig) []CRDModifier {
	return []CRDModifier{
		ModifierConversionWebhook(k8stypes.NamespacedName{Name: parent.GetName() + WebhookNameAffix, Namespace: config.Namespace}),
		ModifierCAInjectAnnotation(k8stypes.NamespacedName{Name: parent.GetName() + WebhookNameAffix, Namespace: config.Namespace}),
	}
}
