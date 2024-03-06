// Copyright © 2022 Cisco Systems, Inc. and/or its affiliates
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

package syslogng

import (
	"context"

	"emperror.dev/errors"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type DataProvider struct {
	client          client.Client
	logging         *v1beta1.Logging
	syslogNGSConfig *v1beta1.SyslogNGConfig
}

func NewDataProvider(client client.Client, logging *v1beta1.Logging, syslogNGSConfig *v1beta1.SyslogNGConfig) *DataProvider {
	return &DataProvider{
		client:          client,
		logging:         logging,
		syslogNGSConfig: syslogNGSConfig,
	}
}

func (p *DataProvider) GetReplicaCount(ctx context.Context) (*int32, error) {
	sts := &v1.StatefulSet{}
	om := p.logging.SyslogNGObjectMeta(StatefulSetName, ComponentSyslogNG, p.syslogNGSConfig)
	err := p.client.Get(ctx, types.NamespacedName{Namespace: om.Namespace, Name: om.Name}, sts)
	if err != nil {
		return nil, errors.WrapIf(client.IgnoreNotFound(err), "getting syslog-ng statefulset")
	}
	return sts.Spec.Replicas, nil
}
