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
	"fmt"
	"io"
	"sort"
	"strings"

	"emperror.dev/errors"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/maps/mapstrstr"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

type FluentRender struct {
	Out    io.Writer
	Indent int
}

func (f *FluentRender) Render(config types.FluentConfig) error {
	return f.RenderDirectives(config.GetDirectives(), 0)
}

func (f *FluentRender) RenderDirectives(directives []types.Directive, indent int) error {
	for _, d := range directives {
		if d == nil {
			continue
		}
		meta := d.GetPluginMeta()
		if meta.Directive == "" {
			return fmt.Errorf("directive must have a name %s", meta)
		}
		// Structural tokens can't be quoted, so a newline would break out.
		for _, t := range []struct{ kind, value string }{
			{"directive name", meta.Directive},
			{"@type", meta.Type},
			{"@id", meta.Id},
			{"@label", meta.Label},
			{"@log_level", meta.LogLevel},
			{"tag", meta.Tag},
		} {
			if err := validateFluentToken(t.kind, t.value); err != nil {
				return err
			}
		}
		f.indentedf(indent, "<%s%s>", meta.Directive, tag(meta.Tag))
		if meta.Type != "" {
			f.indentedf(indent+f.Indent, "@type %s", meta.Type)
		}
		if meta.Id != "" {
			f.indentedf(indent+f.Indent, "@id %s", meta.Id)
		}
		if meta.Label != "" {
			f.indentedf(indent+f.Indent, "@label %s", meta.Label)
		}
		if meta.LogLevel != "" {
			f.indentedf(indent+f.Indent, "@log_level %s", meta.LogLevel)
		}
		if params := d.GetParams(); len(params) > 0 {
			keys := mapstrstr.Keys(params)
			sort.Strings(keys)
			for _, k := range keys {
				if err := validateFluentToken("parameter name", k); err != nil {
					return err
				}
				f.indentedf(indent+f.Indent, "%s %s", k, escapeFluentValue(params[k]))
			}
		}
		if sections := d.GetSections(); len(sections) > 0 {
			if err := f.RenderDirectives(sections, indent+f.Indent); err != nil {
				return errors.WrapIff(err, "failed to render sections for %s", meta.Directive)
			}
		}
		f.indentedf(indent, "</%s>", meta.Directive)
	}
	return nil
}

func (f *FluentRender) indentedf(indent int, format string, values ...any) {
	indentString := strings.Repeat(" ", indent)
	in := fmt.Sprintf(format, values...)
	for line := range strings.SplitSeq(in, "\n") {
		if line != "" {
			fmt.Fprint(f.Out, indentString+line+"\n") //nolint: errcheck
		} else {
			fmt.Fprintln(f.Out, "") //nolint: errcheck
		}
	}
}

func tag(tag string) string {
	if tag != "" {
		return " " + tag
	}
	return tag
}

func validateFluentToken(kind, value string) error {
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("invalid %s %q: must not contain newline characters", kind, value)
	}
	return nil
}

var fluentEscaper = strings.NewReplacer(
	`\`, `\\`,
	`"`, `\"`,
	"\n", `\n`,
	"\r", `\r`,
	"\t", `\t`,
	`#`, `\#`,
)

// escapeFluentValue quotes values containing newlines, escaping '#' too so
// quoting can't enable `#{...}` Ruby interpolation. Plain values are unchanged.
func escapeFluentValue(value string) string {
	if !strings.ContainsAny(value, "\n\r") {
		return value
	}

	return `"` + fluentEscaper.Replace(value) + `"`
}
