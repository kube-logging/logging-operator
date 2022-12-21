// Copyright © 2021 Banzai Cloud
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

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DataProvider struct {
	client client.Client
}

func NewDataProvider(client client.Client) *DataProvider {
	return &DataProvider{
		client: client,
	}
}

func (p *DataProvider) GetReplicaCount(ctx context.Context, logging *v1beta1.Logging) (*int32, error) {
	if logging.Spec.FluentdSpec != nil {
		sts := &v1.StatefulSet{}
		om := logging.FluentdObjectMeta(StatefulSetName, ComponentFluentd)
		err := p.client.Get(ctx, types.NamespacedName{Namespace: om.Namespace, Name: om.Name}, sts)
		if err != nil {
			return nil, errors.WrapIf(client.IgnoreNotFound(err), "getting fluentd statefulset")
		}
		return sts.Spec.Replicas, nil
	}
	return nil, nil
}
