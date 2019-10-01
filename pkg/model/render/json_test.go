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

package render

import (
	"bytes"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/banzaicloud/logging-operator/pkg/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/model/input"
	"github.com/banzaicloud/logging-operator/pkg/model/output"
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

func TestJsonRender(t *testing.T) {
	input, err := input.NewTailInputConfig("input.log").ToDirective(secret.NewSecretLoader(nil, "", nil))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	system := types.NewSystem(input, types.NewRouter())

	flow, err := types.NewFlow(
		"ns-test",
		map[string]string{
			"key1": "val1",
			"key2": "val2",
		})
	if err != nil {
		t.Fatal(err)
	}

	filter, err := filter.NewStdOutFilterConfig().ToDirective(secret.NewSecretLoader(nil, "", nil))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	nullOut, err := output.NewNullOutputConfig().ToDirective(secret.NewSecretLoader(nil, "", nil))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	flow.WithFilters(filter).
		WithOutputs(nullOut)

	err = system.RegisterFlow(flow)
	if err != nil {
		t.Fatal(err)
	}

	configuredSystem, err := system.Build()
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	jsonRender := JsonRender{out: b, indent: 2}
	err = jsonRender.Render(configuredSystem)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{
          "input": {
            "type": "tail",
            "directive": "source",
			"params": {
            	"path": "input.log"
			}
          },
          "router": {
			"type": "label_router",
            "directive": "match",
			"tag": "**",
            "routes": [
              {
                "directive": "route",
                "label": "@901f778f9602a78e8fd702c1973d8d8d",
                "labels": {
                  "key1": "val1",
                  "key2": "val2"
                },
                "namespace": "ns-test"
              }
            ]
          },
          "flows": [
            {
              "directive": "label",
              "tag": "@901f778f9602a78e8fd702c1973d8d8d",
              "filters": [
                {
                  "type": "stdout",
                  "directive": "filter",
                  "tag": "**"
                }
              ],
              "outputs": [
                {
                  "type": "null",
                  "directive": "match",
                  "tag": "**"
                }
              ]
            }
          ]
        }`
	if a, e := diff.TrimLinesInString(b.String()), diff.TrimLinesInString(expected); a != e {
		t.Errorf("Result not as expected:\n%v \nActual: %s", diff.LineDiff(a, e), b.String())
	}
}
