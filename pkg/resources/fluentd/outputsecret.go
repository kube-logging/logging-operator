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

package fluentd

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/secret"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func (r *Reconciler) markSecrets(secrets *secret.MountSecrets) ([]runtime.Object, reconciler.DesiredState, error) {
	var loggingRef string
	if r.Logging.Spec.LoggingRef != "" {
		loggingRef = r.Logging.Spec.LoggingRef
	} else {
		loggingRef = "default"
	}
	annotationKey := fmt.Sprintf("logging.banzaicloud.io/%s", loggingRef)
	var markedSecrets []runtime.Object
	for _, secret := range *secrets {
		secretItem := &corev1.Secret{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{
			Name:      secret.Name,
			Namespace: secret.Namespace}, secretItem)
		if err != nil {
			return nil, reconciler.StatePresent, errors.WrapIfWithDetails(
				err, "failed to load secret", "secret", secret.Name, "namespace", secret.Namespace)
		}
		if secretItem.ObjectMeta.Annotations == nil {
			secretItem.ObjectMeta.Annotations = make(map[string]string)
		}
		secretItem.ObjectMeta.Annotations[annotationKey] = "watched"
		markedSecrets = append(markedSecrets, secretItem)
	}
	return markedSecrets, reconciler.StatePresent, nil
}

func (r *Reconciler) outputSecret(secrets *secret.MountSecrets, mountPath string) (runtime.Object, reconciler.DesiredState, error) {
	// Initialise output secret
	fluentOutputSecret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      r.Logging.QualifiedName(OutputSecretName),
			Namespace: r.Logging.Spec.ControlNamespace,
		},
	}
	if fluentOutputSecret.Data == nil {
		fluentOutputSecret.Data = make(map[string][]byte)
	}
	for _, secret := range *secrets {
		fluentOutputSecret.Data[secret.MappedKey] = secret.Value
	}
	return fluentOutputSecret, reconciler.StatePresent, nil
}
