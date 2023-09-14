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
		targets  []v1beta1.Target
		loggings []*v1beta1.Logging
		wantErr  bool
		satisfy  func([]fluentbit.Tenant) bool
	}{
		{
			name: "static logging target with a static watch namespace list",
			targets: []v1beta1.Target{
				{
					LoggingName: "a",
				},
				{
					LoggingName: "b",
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "a",
					},
					Spec: v1beta1.LoggingSpec{
						WatchNamespaces: []string{"asd"},
					},
				},
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "b",
					},
					Spec: v1beta1.LoggingSpec{
						WatchNamespaces: []string{"bsd"},
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 2) &&
					assert.Equal(t, "a", tenants[0].Logging.Name) &&
					assert.Equal(t, "b", tenants[1].Logging.Name) &&
					assert.Equal(t, []string{"asd"}, tenants[0].Namespaces) &&
					assert.Equal(t, []string{"bsd"}, tenants[1].Namespaces) &&
					assert.False(t, tenants[0].AllNamespace) &&
					assert.False(t, tenants[1].AllNamespace)
			},
		},
		{
			name: "static logging target with an empty watch namespace will be omitted",
			targets: []v1beta1.Target{
				{
					LoggingName: "a",
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: "a",
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 0)
			},
		},
		{
			name: "dynamic logging targets with a static watch namespace list",
			targets: []v1beta1.Target{
				{
					LoggingSelector: &v1.LabelSelector{
						MatchExpressions: []v1.LabelSelectorRequirement{
							{
								Key:      "tenant",
								Operator: v1.LabelSelectorOpExists,
							},
						},
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
					assert.Equal(t, "b", tenants[0].Logging.Name) &&
					assert.Equal(t, "c", tenants[1].Logging.Name) &&
					assert.Equal(t, []string{"bsd"}, tenants[0].Namespaces) &&
					assert.Equal(t, []string{"csd"}, tenants[1].Namespaces) &&
					assert.False(t, tenants[0].AllNamespace) &&
					assert.False(t, tenants[1].AllNamespace)
			},
		},
		{
			name: "allNamespace allowed for self referencing target",
			targets: []v1beta1.Target{
				{
					LoggingName: currentLoggingName,
				},
			},
			loggings: []*v1beta1.Logging{
				{
					ObjectMeta: v1.ObjectMeta{
						Name: currentLoggingName,
					},
				},
			},
			satisfy: func(tenants []fluentbit.Tenant) bool {
				return assert.Len(t, tenants, 1) &&
					assert.True(t, tenants[0].AllNamespace)
			},
		},
	}

	for _, td := range testData {
		td := td
		deferred := ensureCreatedAll(t, td.loggings)
		assert.Eventually(t, func() bool {
			tenants, err := fluentbit.FindTenants(context.TODO(), td.targets, currentLoggingName, mgr.GetClient(), mgr.GetLogger())
			if td.wantErr {
				assert.NoError(t, err)
			}
			return td.satisfy(tenants)
		}, time.Second, 100*time.Millisecond, td.name)
		deferred()
	}
}
