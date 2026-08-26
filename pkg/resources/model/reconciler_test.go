// Copyright © 2024 Kube logging authors
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

package model

import (
	"context"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

type stubSecretLoaderFactory struct {
	client client.Client
}

func (f stubSecretLoaderFactory) OutputSecretLoaderForNamespace(namespace string) secret.SecretLoader {
	return secret.NewSecretLoader(f.client, namespace, "", &secret.MountSecrets{})
}

// Regression: a ClusterOutput referenced only by Logging.spec.defaultFlow.globalOutputRefs
// is rendered into the fluentd config by FlowForDefaultFlow, but the validation reconciler
// reported it as inactive because it only looked at Flow and ClusterFlow references.
func TestValidationReconciler_DefaultFlowRef(t *testing.T) {
	tests := []struct {
		name       string
		outputSpec v1beta1.OutputSpec
		refs       []string
		want       bool
	}{
		{
			name:       "valid output becomes active",
			outputSpec: v1beta1.OutputSpec{NullOutputConfig: output.NewNullOutputConfig()},
			refs:       []string{"general"},
			want:       true,
		},
		{
			// Same semantics as the flow and clusterflow loops: no target configured is a problem
			name:       "output with problems stays inactive",
			outputSpec: v1beta1.OutputSpec{},
			refs:       []string{"general"},
			want:       false,
		},
		{
			name:       "unreferenced output stays inactive",
			outputSpec: v1beta1.OutputSpec{NullOutputConfig: output.NewNullOutputConfig()},
			refs:       []string{"other"},
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterOutput := v1beta1.ClusterOutput{
				ObjectMeta: metav1.ObjectMeta{Name: "general", Namespace: "test"},
				Spec:       v1beta1.ClusterOutputSpec{OutputSpec: tt.outputSpec},
			}
			logging := loggingFor()
			logging.Spec.DefaultFlowSpec = &v1beta1.DefaultFlowSpec{GlobalOutputRefs: tt.refs}

			scheme := runtime.NewScheme()
			if err := v1beta1.AddToScheme(scheme); err != nil {
				t.Fatalf("add to scheme: %v", err)
			}
			cl := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(&clusterOutput, &logging).
				WithStatusSubresource(&clusterOutput, &logging).
				Build()

			res := LoggingResources{
				Logging: logging,
				Fluentd: FluentdLoggingResources{ClusterOutputs: ClusterOutputs{clusterOutput}},
			}

			reconcile := NewValidationReconciler(cl, res, stubSecretLoaderFactory{client: cl}, logr.Discard())
			if _, err := reconcile(context.Background()); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			active := res.Fluentd.ClusterOutputs[0].Status.Active
			if got := active != nil && *active; got != tt.want {
				t.Errorf("Status.Active = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasIntersection(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no intersection empty",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: false,
		},
		{
			name: "no intersection nonempty",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"d", "e"},
			},
			want: false,
		},
		{
			name: "has intersection",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"b"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIntersection(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("hasIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
