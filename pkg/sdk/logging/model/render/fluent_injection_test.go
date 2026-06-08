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

package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

func renderDirective(d types.Directive) (string, error) {
	b := bytes.Buffer{}
	renderer := render.FluentRender{Out: &b, Indent: 2}
	err := renderer.RenderDirectives([]types.Directive{d}, 0)
	return b.String(), err
}

func TestRenderDirective_Injection(t *testing.T) {
	tests := []struct {
		name        string
		directive   types.Directive
		wantErr     bool
		contains    []string
		notContains []string
	}{
		{
			name: "newline in value is escaped into a single quoted line",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{Directive: "record"},
				Params:     types.Params{"x": "y\n</record>\n</filter>\n<match **>\n  @type exec\n  command id\n</match>"},
			},
			contains:    []string{`x "y\n</record>\n</filter>\n<match **>\n  @type exec\n  command id\n</match>"`},
			notContains: []string{"\n<match **>\n", "\n  @type exec\n"},
		},
		{
			name: "ruby interpolation is neutralized when quoting",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{Directive: "record"},
				Params:     types.Params{"x": "a\n#{Socket.gethostname}"},
			},
			contains:    []string{`x "a\n\#{Socket.gethostname}"`},
			notContains: []string{`"a\n#{Socket.gethostname}"`},
		},
		{
			name: "quotes, backslashes and tabs are escaped",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{Directive: "record"},
				Params:     types.Params{"x": "a\nb\"c\\d\te"},
			},
			contains: []string{`x "a\nb\"c\\d\te"`},
		},
		{
			name: "value without newline is left unquoted",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{Directive: "record"},
				Params:     types.Params{"foo": "bar", "labels": `{"a":"b"}`},
			},
			contains: []string{"foo bar", `labels {"a":"b"}`},
		},
		{
			name: "newline in parameter name is rejected",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{Directive: "record"},
				Params:     types.Params{"bad\n</record>\n<match **>": "v"},
			},
			wantErr: true,
		},
		{
			name:      "newline in tag is rejected",
			directive: &types.GenericDirective{PluginMeta: types.PluginMeta{Directive: "match", Tag: "x>\n<match **"}},
			wantErr:   true,
		},
		{
			name:      "newline in id is rejected",
			directive: &types.GenericDirective{PluginMeta: types.PluginMeta{Directive: "match", Id: "x\n@type exec"}},
			wantErr:   true,
		},
		{
			name:      "newline in label is rejected",
			directive: &types.GenericDirective{PluginMeta: types.PluginMeta{Directive: "match", Label: "x\n@type exec"}},
			wantErr:   true,
		},
		{
			name:      "newline in log_level is rejected",
			directive: &types.GenericDirective{PluginMeta: types.PluginMeta{Directive: "match", LogLevel: "info\n@type exec"}},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := renderDirective(tt.directive)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil; output:\n%s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, want := range tt.contains {
				if !strings.Contains(out, want) {
					t.Errorf("expected output to contain %q, got:\n%s", want, out)
				}
			}
			for _, unwanted := range tt.notContains {
				if strings.Contains(out, unwanted) {
					t.Errorf("expected output not to contain %q, got:\n%s", unwanted, out)
				}
			}
		})
	}
}
