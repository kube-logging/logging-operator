// Copyright © 2023 Kube logging authors
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

package controllers_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentbit"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestFindTenants(t *testing.T) {
	defer beforeEach(t)()

	currentLoggingName := "current"

	testData := []struct {
		name     string
		target   v1.LabelSelector
		loggings []*v1beta1.Logging
		wantErr  bool
		satisfy  func([]fluentbit.Tenant) bool
	}{
		{
			name: "static logging target with a static watch namespace list",
			target: v1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "a",
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "a",
						Labels: map[string]string{
							"name": "a",
						},
					},
					Spec: v1beta1.LoggingSpec{
						WatchNamespaces: []string{"asd"},
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 1) &&
					assert.Contains(t, tenants, fluentbit.Tenant{
						Name:         "a",
						AllNamespace: false,
						Namespaces:   []string{"asd"},
					})
			},
		},
		{
			name: "watching all namespaces will result in an empty namespace list, with the allNamespaces flag set to true",
			target: v1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "a",
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "a",
						Labels: map[string]string{
							"name": "a",
						},
					},
					Spec: v1beta1.LoggingSpec{
						WatchNamespaces:        nil,
						WatchNamespaceSelector: nil,
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 1) &&
					assert.Contains(t, tenants, fluentbit.Tenant{
						Name:         "a",
						AllNamespace: true,
					})
			},
		},
		{
			name: "static logging target with an empty watch namespace will be omitted",
			target: v1.LabelSelector{
				MatchLabels: map[string]string{
					"name": currentLoggingName,
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "a",
						Labels: map[string]string{
							"name": "a",
						},
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 0)
			},
		},
		{
			name: "dynamic logging targets with a static watch namespace list",
			target: v1.LabelSelector{
				MatchExpressions: []v1.LabelSelectorRequirement{
					{
						Key:      "tenant",
						Operator: v1.LabelSelectorOpExists,
					},
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "b",
						Labels: map[string]string{
							"tenant": "x",
						},
					},
					Spec: v1beta1.LoggingSpec{
						WatchNamespaces: []string{"bsd"},
					},
				},
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "c",
						Labels: map[string]string{
							"tenant": "y",
						},
					},
					Spec: v1beta1.LoggingSpec{
						WatchNamespaces: []string{"csd"},
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 2) &&
					assert.Contains(t, tenants, fluentbit.Tenant{
						Name:         "b",
						AllNamespace: false,
						Namespaces:   []string{"bsd"},
					}) &&
					assert.Contains(t, tenants, fluentbit.Tenant{
						Name:         "c",
						AllNamespace: false,
						Namespaces:   []string{"csd"},
					})
			},
		},
	}

	for _, td := range testData {
		td := td
		deferred := ensureCreatedAll(t, td.loggings)
		assert.Eventually(t, func() bool {
			tenants, err := fluentbit.FindTenants(context.TODO(), td.target, mgr.GetClient())
			if td.wantErr {
				assert.NoError(t, err)
			}
			return td.satisfy(tenants)
		}, time.Second, 100*time.Millisecond, td.name)
		deferred()
	}
}
