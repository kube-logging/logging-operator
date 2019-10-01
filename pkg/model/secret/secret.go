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

package secret

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:object:generate=true

type Secret struct {
	Value     string     `json:"value,omitempty"`
	ValueFrom *ValueFrom `json:"valueFrom,omitempty"`
	MountFrom *ValueFrom `json:"mountFrom,omitempty"`
}

// +kubebuilder:object:generate=true

type ValueFrom struct {
	SecretKeyRef *KubernetesSecret `json:"secretKeyRef,omitempty"`
}

// +kubebuilder:object:generate=true

type KubernetesSecret struct {
	// Name of the kubernetes secret
	Name string `json:"name"`
	// Secret key for the value
	Key string `json:"key"`
}

type SecretLoader interface {
	Load(secret *Secret) (string, error)
	Mount(secret *Secret) (string, error)
}

type secretLoader struct {
	// secretLoader is limited to a single namespace, to avoid hijacking other namespace's secrets
	namespace   string
	client      client.Client
	mountConfig *MountConfig
}

// MountConfig to attach volume secrets to Fluentd
type MountConfig struct {
	SecretName      string
	SecretNamespace string
	ConfigPath      string
	LoggingRef      string
}

func NewSecretLoader(client client.Client, namespace string, mountConfig *MountConfig) *secretLoader {
	return &secretLoader{
		client:      client,
		namespace:   namespace,
		mountConfig: mountConfig,
	}
}

func (k *secretLoader) Mount(secret *Secret) (string, error) {
	k8sSecret := &corev1.Secret{}
	err := k.client.Get(context.TODO(), types.NamespacedName{
		Name:      secret.MountFrom.SecretKeyRef.Name,
		Namespace: k.namespace}, k8sSecret)
	if err != nil {
		return "", errors.WrapIff(err, "failed to get kubernetes secret %s:%s",
			k.namespace,
			secret.MountFrom.SecretKeyRef.Name)
	}
	secretKey := fmt.Sprintf("%s-%s-%s", k.namespace, secret.MountFrom.SecretKeyRef.Name, secret.MountFrom.SecretKeyRef.Key)
	value, ok := k8sSecret.Data[secret.MountFrom.SecretKeyRef.Key]
	if !ok {
		return "", errors.Errorf("key %q not found in secret %q in namespace %q",
			secret.MountFrom.SecretKeyRef.Key,
			secret.MountFrom.SecretKeyRef.Name,
			k.namespace)
	}
	fluentOutputSecret := &corev1.Secret{}
	err = k.client.Get(context.TODO(), types.NamespacedName{
		Name:      k.mountConfig.SecretName,
		Namespace: k.mountConfig.SecretNamespace}, fluentOutputSecret)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to get fluent output secret: %q", k.mountConfig.SecretName))
	}
	if fluentOutputSecret.Data == nil {
		fluentOutputSecret.Data = make(map[string][]byte)
	}
	fluentOutputSecret.Data[secretKey] = value
	err = k.client.Update(context.TODO(), fluentOutputSecret)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to update fluent output secret: %q", k.mountConfig.SecretName))
	}
	var loggingRef string
	if k.mountConfig.LoggingRef != "" {
		loggingRef = k.mountConfig.LoggingRef
	} else {
		loggingRef = "default"
	}
	annotationKey := fmt.Sprintf("logging.banzaicloud.io/%s", loggingRef)
	if k8sSecret.ObjectMeta.Annotations == nil {
		k8sSecret.ObjectMeta.Annotations = make(map[string]string)
	}
	k8sSecret.ObjectMeta.Annotations[annotationKey] = "watched"
	err = k.client.Update(context.TODO(), k8sSecret)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to update source secret: %q", secret.MountFrom.SecretKeyRef.Name))
	}
	return k.mountConfig.ConfigPath + "/" + secretKey, nil
}

func (k *secretLoader) Load(secret *Secret) (string, error) {
	if secret.Value != "" {
		return secret.Value, nil
	}

	if secret.ValueFrom.SecretKeyRef != nil {
		k8sSecret := &corev1.Secret{}
		err := k.client.Get(context.TODO(), types.NamespacedName{
			Name:      secret.ValueFrom.SecretKeyRef.Name,
			Namespace: k.namespace}, k8sSecret)
		if err != nil {
			return "", errors.WrapIff(err, "failed to get kubernetes secret %s:%s",
				k.namespace,
				secret.ValueFrom.SecretKeyRef.Name)
		}
		value, ok := k8sSecret.Data[secret.ValueFrom.SecretKeyRef.Key]
		if !ok {
			return "", errors.Errorf("key %q not found in secret %q in namespace %q",
				secret.ValueFrom.SecretKeyRef.Key,
				secret.ValueFrom.SecretKeyRef.Name,
				k.namespace)
		}
		return string(value), nil
	}

	return "", errors.New("No secret Value or ValueFrom defined for field")
}
