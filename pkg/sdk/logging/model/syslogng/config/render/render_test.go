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

package render_test

import (
	"strings"
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/stretchr/testify/assert"
)

func TestStringList(t *testing.T) {
	var expectedStringList string = `"field1", "field2", "field3"`
	var stringList = []string{"field1", "field2", "field3"}
	var renderer = render.StringList(stringList)

	options := config.OutputConfigCheckOptions{
		IndentWith: "    ",
	}
	actualStringList := &strings.Builder{}
	err := renderer(render.RenderContext{
		Out:        actualStringList,
		IndentWith: options.IndentWith,
	})
	config.CheckError(t, options.ExpectedError, err)
	assert.Equal(t, expectedStringList, actualStringList.String())

}
