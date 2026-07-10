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
// Configure custom or unexposed Fluentd filters via raw configuration. The configuration is parsed and rendered by the operator (parameter ordering and duplicate keys are not preserved).
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
          flatten true
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
  flatten true
  key_name ua_string
</filter>
{{</ highlight >}}

*/
type _docRaw interface{} //nolint:deadcode,unused

// +name:"Raw"
// +url:""
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

	if raw.Type == "" {
		return nil, fmt.Errorf("raw filter config must specify @type")
	}

	raw.Id = id
	raw.Tag = "**"
	raw.Directive = "filter"

	return raw, nil
}

func parseRawConfig(config string) (*types.GenericDirective, error) {
	scanner := bufio.NewScanner(strings.NewReader(config))
	// Allow reasonably large raw configs (default token limit is 64K).
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// nextLine should return:
	// line, eof, error
	// eof=true means the end of input
	nextLine := func() (string, bool, error) {
		if scanner.Scan() {
			return scanner.Text(), false, nil
		}

		if err := scanner.Err(); err != nil {
			return "", true, err
		}
		return "", true, nil
	}

	return doParseRawConfig("filter", nextLine)
}

func doParseRawConfig(sectionName string, nextLine func() (string, bool, error)) (*types.GenericDirective, error) {
	directive := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Directive: sectionName,
		},
		Params:        types.Params{},
		SubDirectives: []types.Directive{},
	}

	for {
		line, eof, err := nextLine()
		if err != nil {
			return nil, err
		}
		if eof {
			if sectionName != "filter" {
				return nil, fmt.Errorf("unexpected end of raw config: missing closing tag </%s>", sectionName)
			}
			return directive, nil
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if matches := sectionPattern.FindStringSubmatch(line); matches != nil {
			subSectionName := matches[1]
			subSectionTag := strings.TrimSpace(matches[2])
			subSectionDirective, err := doParseRawConfig(subSectionName, nextLine)
			if err != nil {
				return nil, err
			}

			if subSectionTag != "" {
				subSectionDirective.Tag = subSectionTag
			}

			directive.SubDirectives = append(directive.SubDirectives, subSectionDirective)
			continue
		}

		if line == "</"+sectionName+">" {
			break
		}

		if matches := paramPattern.FindStringSubmatch(line); matches != nil {
			paramName := matches[1]
			paramValue := strings.TrimSpace(matches[2])

			if paramName == "@id" {
				continue // ignore @id parameter, as it is set by the operator
			}

			if paramName == "@type" {
				directive.Type = paramValue
			} else {
				directive.Params[paramName] = paramValue
			}
			continue
		}

		return nil, fmt.Errorf("invalid line in raw config: %s", line)
	}
	return directive, nil
}
