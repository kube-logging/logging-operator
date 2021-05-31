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

	"github.com/banzaicloud/operator-tools/pkg/utils"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
)

type CRDModifier func(*v1.CustomResourceDefinition) (*v1.CustomResourceDefinition, error)

const (
	DefaultSecurePort                    int32  = 443
	DefaultWebhookPort                   int32  = 9443
	CertManagerInjectCAFromAnnotationKey string = "cert-manager.io/inject-ca-from"
	WebhookNameAffix                     string = "-logging-operator-webhooks"
	WebhookCertDir                       string = "/tmp/k8s-webhook-server/serving-certs"
)

var ConversionReviewVersions = []string{
	"v1beta1",
	"v1alpha1",
}

func ModifierConversionWebhook(svc types.NamespacedName) CRDModifier {
	svcRef := v1.ServiceReference{
		Namespace: svc.Namespace,
		Name:      svc.Name,
		Path:      utils.StringPointer("/convert"),
		Port:      utils.IntPointer(DefaultSecurePort),
	}
	return func(crd *v1.CustomResourceDefinition) (*v1.CustomResourceDefinition, error) {
		crd.Spec.Conversion = &v1.CustomResourceConversion{

			Strategy: v1.WebhookConverter,
			Webhook: &v1.WebhookConversion{

				ClientConfig: &v1.WebhookClientConfig{
					Service: &svcRef,
				},
				ConversionReviewVersions: ConversionReviewVersions,
			},
		}
		return crd, nil
	}
}

func ModifierCAInjectAnnotation(certName types.NamespacedName) CRDModifier {
	return func(crd *v1.CustomResourceDefinition) (*v1.CustomResourceDefinition, error) {
		if crd.Annotations == nil {
			crd.Annotations = make(map[string]string)
		}
		crd.Annotations[CertManagerInjectCAFromAnnotationKey] = fmt.Sprintf("%s/%s", certName.Namespace, certName.Name)
		return crd, nil
	}
}
