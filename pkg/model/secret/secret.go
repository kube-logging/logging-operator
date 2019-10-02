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
	namespace string
	mountPath string
	client    client.Client
	secrets   *MountSecrets
}

type MountSecrets struct {
	Secrets []MountSecret
}

func (m *MountSecrets) Append(secret MountSecret) {
	m.Secrets = append(m.Secrets, secret)
}

func (m *MountSecrets) List() []MountSecret {
	return m.Secrets
}

type MountSecret struct {
	Name      string
	Key       string
	Namespace string
}

func NewSecretLoader(client client.Client, namespace, mountPath string, secrets *MountSecrets) *secretLoader {
	return &secretLoader{
		client:    client,
		mountPath: mountPath,
		namespace: namespace,
		secrets:   secrets,
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
	k.secrets.Append(MountSecret{
		Name:      secret.MountFrom.SecretKeyRef.Name,
		Key:       secret.MountFrom.SecretKeyRef.Key,
		Namespace: k.namespace,
	})
	return k.mountPath + "/" + secretKey, nil
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
