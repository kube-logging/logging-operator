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

package filter

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

// +name:"Raw"
// +weight:"200"
type _hugoRaw interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"Raw"
// Configure custom or unexposed Fluentd filters via raw configuration. This allows you to specify any configuration that is not supported by the operator. The configuration should be in the format of a Fluentd filter configuration.
/*
## Example `Raw` filter configurations

### Configure a custom filter via raw configuration

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - raw:
        config: |
          @type my_filter
          <my_section>
            foo bar
        	tags ["web", "api", "db"]
          </my_section>
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
  @type my_filter
  @id test
  <my_section>
    foo bar
    tags ["web", "api", "db"]
  </my_section>
</filter>
{{</ highlight >}}

### Configure an unexposed filter via raw configuration

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - raw:
        config: |
          @type ua_parser
          flatten
          key_name ua_string
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
  @type ua_parser
  @id test
  flatten
  key_name ua_string
</filter>
{{</ highlight >}}

*/
type _docRaw interface{} //nolint:deadcode,unused

// +name:"Raw"
// +url:"TODO"
// +version:""
// +description:"Configure raw filter."
// +status:""
type _metaRaw interface{} //nolint:deadcode,unused

var (
	sectionPattern = regexp.MustCompile(`^<([a-zA-Z0-9_]+)\s*(.+?)?>$`)
	paramPattern   = regexp.MustCompile(`^([a-zA-Z0-9_@]+)\s*(.*)$`)
)

// +kubebuilder:object:generate=true
type Raw struct {
	// Raw configuration for the filter.
	Config string `json:"config,omitempty"`
}

func (r *Raw) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	raw, err := parseRawConfig(r.Config)
	if err != nil {
		return nil, err
	}

	raw.Id = id
	raw.Tag = "**"
	raw.Directive = "filter"

	return raw, nil
}

func parseRawConfig(config string) (*types.GenericDirective, error) {
	scanner := bufio.NewScanner(strings.NewReader(config))

	// nextLine should return:
	// line, eof
	// eof=true means the end of input
	nextLine := func() (string, bool) {
		if scanner.Scan() {
			return scanner.Text(), false
		}
		return "", true
	}

	return doParseRawConfig("filter", nextLine)
}

func doParseRawConfig(sectionName string, nextLine func() (string, bool)) (*types.GenericDirective, error) {
	directive := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Directive: sectionName,
		},
		Params:        types.Params{},
		SubDirectives: []types.Directive{},
	}

	for {
		line, eof := nextLine()
		if eof {
			return directive, nil
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if matches := sectionPattern.FindStringSubmatch(line); matches != nil {
			subSectionName := matches[1]
			subSectionDirective, err := doParseRawConfig(subSectionName, nextLine)
			if err != nil {
				return nil, err
			}
			directive.SubDirectives = append(directive.SubDirectives, subSectionDirective)
			continue
		}

		if line == "</"+sectionName+">" {
			break
		}

		if matches := paramPattern.FindStringSubmatch(line); matches != nil {
			paramName := matches[1]
			paramValue := matches[2]

			if paramName == "@id" {
				continue // ignore @id parameter, as it is set by the operator
			}

			if paramName == "@type" {
				directive.PluginMeta.Type = paramValue
			} else {
				directive.Params[paramName] = paramValue
			}
			continue
		}

		return nil, fmt.Errorf("invalid line in raw config: %s", line)
	}
	return directive, nil
}
