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

package docgen_test

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/andreyvit/diff"
	"github.com/banzaicloud/logging-operator/pkg/docgen"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger logr.Logger

func init() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = zapr.NewLogger(log)
}

func TestGenParse(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	var testData = []struct {
		docItem  docgen.DocItem
		expected string
	}{
		{
			docItem: docgen.DocItem{
				Name:       "sample-name",
				SourcePath: filepath.Join(currentDir, "testdata", "sample.go"),
				DestPath:   filepath.Join(currentDir, "../../build/_test/docgen"),
			},
			expected: heredoc.Doc(`
				### Sample
				| Variable Name | Type | Required | Default | Description |
				|---|---|---|---|---|
				| field1 | string | No | - |  |
			`),
		},
		{
			docItem: docgen.DocItem{
				Name:       "sample-name",
				SourcePath: filepath.Join(currentDir, "testdata", "sample_default.go"),
				DestPath:   filepath.Join(currentDir, "../../build/_test/docgen"),
				DefaultValueFromTagExtractor: func(tag string) string {
					return docgen.GetPrefixedValue(tag, `asd:\"default:(.*)\"`)
				},
			},
			expected: heredoc.Doc(`
				### SampleDefault
				| Variable Name | Type | Required | Default | Description |
				|---|---|---|---|---|
				| field1 | string | No | testval |  |
			`),
		},
	}

	for _, item := range testData {
		parser := docgen.GetDocumentParser(item.docItem, logger)
		err := parser.Generate()
		if err != nil {
			t.Fatalf("%+v", err)
		}

		bytes, err := ioutil.ReadFile(filepath.Join(item.docItem.DestPath, item.docItem.Name+".md"))
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if a, e := diff.TrimLinesInString(string(bytes)), diff.TrimLinesInString(item.expected); a != e {
			t.Errorf("Result does not match (-actual vs +expected):\n%v\nActual: %s", diff.LineDiff(a, e), string(bytes))
		}

		if err != nil {
			t.Fatalf("%+v", err)
		}
	}
}
