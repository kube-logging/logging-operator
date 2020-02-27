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

	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInvalidFlowIfMatchAndSelectorBothSet(t *testing.T) {
	defer beforeEach(t)()

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
			OutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to create model: match and selectors cannot be defined simultaneously for flow %s",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected)
}

func TestInvalidFlowIfSelectorAndExcludeBothSet(t *testing.T) {
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
					Exclude: &v1beta1.Exclude{
						Labels: map[string]string{
							"c": "d",
						},
					},
				},
			},
			OutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to create model: failed to process match for %s: select and exclude cannot be set simultaneously",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected)
}

func TestInvalidClusterFlowIfSelectorAndExcludeBothSet(t *testing.T) {
	defer beforeEach(t)()

	logging := testLogging()
	output := testOutput()

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
			OutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to create model: failed to process match for %s: select and exclude cannot be set simultaneously",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected)
}

func TestInvalidClusterFlowIfMatchAndSelectorBothSet(t *testing.T) {
	defer beforeEach(t)()

	logging := testLogging()
	output := testOutput()

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
			OutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("failed to create model: match and selectors cannot be defined simultaneously for clusterflow %s",
		utils.ObjectKeyFromObjectMeta(flow).String(),
	)

	expectError(t, expected)
}
