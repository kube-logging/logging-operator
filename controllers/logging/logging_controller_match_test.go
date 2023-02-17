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

package controllers_test

import (
	"fmt"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/andreyvit/diff"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFlowMatch(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := testLogging()
	output := testOutput()

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: output.Namespace,
		},
		Spec: v1beta1.FlowSpec{
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			LocalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(diff.TrimLinesInString(string(secret.Data[fluentd.AppConfigKey]))).Should(gomega.ContainSubstring(diff.TrimLinesInString(heredoc.Docf(`
		<match>
		  labels c:d
		  namespaces %s
		  negate false
		</match>
	`, flow.Namespace))))
}

func TestClusterFlowMatch(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := testLogging()
	output := testClusterOutput()

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			GlobalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(diff.TrimLinesInString(string(secret.Data[fluentd.AppConfigKey]))).Should(gomega.ContainSubstring(diff.TrimLinesInString(heredoc.Docf(`
		<match>
		  labels c:d
		  negate false
		</match>
	`))))
}

func TestClusterFlowMatchWithNamespaces(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := testLogging()
	output := testClusterOutput()

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						Labels: map[string]string{
							"c": "d",
						},
						Namespaces: []string{"a", "b"},
					},
				},
			},
			GlobalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(diff.TrimLinesInString(string(secret.Data[fluentd.AppConfigKey]))).Should(gomega.ContainSubstring(diff.TrimLinesInString(heredoc.Docf(`
		<match>
		  labels c:d
		  namespaces a,b
		  negate false
		</match>
	`))))
}

func TestInvalidFlowIfMatchAndSelectorBothSet(t *testing.T) {
	errors := make(chan error)
	defer beforeEachWithError(t, errors)()

	logging := testLogging()
	output := testOutput()

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: output.Namespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			LocalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to build model: match and selectors cannot be defined simultaneously for flow %s",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected, errors)
}

func TestInvalidFlowIfSelectorAndExcludeBothSet(t *testing.T) {
	errors := make(chan error)
	defer beforeEachWithError(t, errors)()

	logging := testLogging()
	output := testOutput()

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: output.Namespace,
		},
		Spec: v1beta1.FlowSpec{
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{
						Labels: map[string]string{
							"c": "d",
						},
					},
					Exclude: &v1beta1.Exclude{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			LocalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to build model: select and exclude cannot be set simultaneously for flow %s",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected, errors)
}

func TestInvalidClusterFlowIfSelectorAndExcludeBothSet(t *testing.T) {
	errors := make(chan error)
	defer beforeEachWithError(t, errors)()

	logging := testLogging()
	output := testClusterOutput()

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						Labels: map[string]string{
							"c": "d",
						},
					},
					ClusterExclude: &v1beta1.ClusterExclude{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			GlobalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to build model: select and exclude cannot be set simultaneously for clusterflow %s",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected, errors)
}

func TestInvalidClusterFlowIfMatchAndSelectorBothSet(t *testing.T) {
	errors := make(chan error)
	defer beforeEachWithError(t, errors)()

	logging := testLogging()
	output := testClusterOutput()

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			GlobalOutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to build model: match and selectors cannot be defined simultaneously for clusterflow %s",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected, errors)
}
