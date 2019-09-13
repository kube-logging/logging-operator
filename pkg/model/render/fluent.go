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
	"strings"

	"github.com/banzaicloud/logging-operator/pkg/model/types"
	"github.com/banzaicloud/logging-operator/pkg/util"
	"github.com/goph/emperror"
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
		meta := d.GetPluginMeta()
		if meta.Directive == "" {
			return fmt.Errorf("Directive must have a name %s", meta)
		}
		f.indented(indent, "<%s%s>", meta.Directive, tag(meta.Tag))
		if meta.Type != "" {
			f.indented(indent+f.Indent, "@type %s", meta.Type)
		}
		if meta.Id != "" {
			f.indented(indent+f.Indent, "@id %s", meta.Id)
		}
		if meta.Label != "" {
			f.indented(indent+f.Indent, "@label %s", meta.Label)
		}
		if meta.LogLevel != "" {
			f.indented(indent+f.Indent, "@log_level %s", meta.LogLevel)
		}
		if len(d.GetParams()) > 0 {
			for _, k := range util.OrderedStringMap(d.GetParams()).Keys() {
				f.indented(indent+f.Indent, "%s %s", k, d.GetParams()[k])
			}
		}
		if len(d.GetSections()) > 0 {
			err := f.RenderDirectives(d.GetSections(), indent+f.Indent)
			if err != nil {
				return emperror.Wrapf(err, "failed to render sections for %s", d.GetPluginMeta().Directive)
			}
		}
		f.indented(indent, "</%s>", meta.Directive)
	}
	return nil
}

func (f *FluentRender) indented(indent int, format string, values ...interface{}) {
	indentString := strings.Repeat(" ", indent)
	in := fmt.Sprintf(format, values...)
	for _, line := range strings.Split(in, "\n") {
		if line != "" {
			fmt.Fprint(f.Out, indentString+line+"\n")
		} else {
			fmt.Fprintln(f.Out, "")
		}
	}

}

func tag(tag string) string {
	if tag != "" {
		return " " + tag
	}
	return tag
}
