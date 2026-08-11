// Copyright © 2026 Kube logging authors
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
type _hugoRaw any //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"Raw"
// Configure custom or unexposed Fluentd filters via raw configuration. The configuration is parsed and rendered by the operator (parameter ordering and duplicate keys are not preserved). Disabled by default, set `logging.spec.enableRawFluentdFilter=true` on the Logging resource to enable it.
//
// > **Security warning:** enabling the flag grants **every** user who can create a `Flow` or `ClusterFlow` in the logging domain the ability to inject arbitrary configuration (and thereby run arbitrary code, for example via `exec`-style plugins) in the shared Fluentd aggregator, which processes all tenants' logs and mounts output credentials and TLS keys. Treat it as a break-glass setting and enable it only if you trust all Flow authors in the cluster. Setting the flag back to `false` does not retroactively remove already rendered raw configuration: delete the offending `Flow`, then rotate the aggregator's credentials and restart it.
//
// The configuration is limited to 64 KiB and 32 levels of nesting.
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
type _docRaw any //nolint:deadcode,unused

// +name:"Raw"
// +url:""
// +version:""
// +description:"Configure raw filter."
// +status:""
type _metaRaw any //nolint:deadcode,unused

var (
	sectionPattern = regexp.MustCompile(`^<([^\s>/]+)\s*([^>]*)>$`)
	closingPattern = regexp.MustCompile(`^</([^\s>/]+)>$`)
	paramPattern   = regexp.MustCompile(`^([^\s<]+)\s*(.*)$`)
)

const (
	// Rendered output grows with the square of the nesting depth (every level
	// indents all enclosed lines), so depth must stay bounded to keep a single
	// Flow from exhausting the operator's memory.
	maxRawConfigDepth = 32
	maxRawConfigBytes = 64 * 1024
)

// +kubebuilder:object:generate=true
type Raw struct {
	// Raw configuration for the filter.
	// +kubebuilder:validation:MaxLength=65536
	Config string `json:"config,omitempty"`
}

func (r *Raw) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	raw, err := parseRawConfig(r.Config)
	if err != nil {
		return nil, err
	}

	if raw.Type == "" {
		if len(raw.SubDirectives) == 1 && raw.SubDirectives[0].GetPluginMeta().Directive == "filter" {
			return nil, fmt.Errorf("raw filter config must not include the enclosing <filter> tags, provide only the filter body")
		}
		return nil, fmt.Errorf("raw filter config must specify @type")
	}

	raw.Id = id
	raw.Tag = "**"
	raw.Directive = "filter"

	return raw, nil
}

func parseRawConfig(config string) (*types.GenericDirective, error) {
	if len(config) > maxRawConfigBytes {
		return nil, fmt.Errorf("raw config is too large: %d bytes exceeds maximum of %d", len(config), maxRawConfigBytes)
	}

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

	return doParseRawConfig("filter", true, 0, nextLine)
}

func doParseRawConfig(sectionName string, topLevel bool, depth int, nextLine func() (string, bool, error)) (*types.GenericDirective, error) {
	if depth > maxRawConfigDepth {
		return nil, fmt.Errorf("raw config nesting is too deep: exceeds maximum of %d levels", maxRawConfigDepth)
	}

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
			if !topLevel {
				return nil, fmt.Errorf("unexpected end of raw config: missing closing tag </%s>", sectionName)
			}
			return directive, nil
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if matches := closingPattern.FindStringSubmatch(line); matches != nil {
			if topLevel || matches[1] != sectionName {
				return nil, fmt.Errorf("unexpected closing tag in raw config: %s", line)
			}
			break
		}

		if matches := sectionPattern.FindStringSubmatch(line); matches != nil {
			subSectionName := matches[1]
			subSectionTag := strings.TrimSpace(matches[2])
			subSectionDirective, err := doParseRawConfig(subSectionName, false, depth+1, nextLine)
			if err != nil {
				return nil, err
			}

			if subSectionTag != "" {
				subSectionDirective.Tag = subSectionTag
			}

			directive.SubDirectives = append(directive.SubDirectives, subSectionDirective)
			continue
		}

		if matches := paramPattern.FindStringSubmatch(line); matches != nil {
			paramName := matches[1]
			paramValue := strings.TrimSpace(matches[2])

			if paramName == "@id" {
				directive.Id = paramValue // top-level id will be overwritten by the operator; nested ids are preserved
				continue
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
