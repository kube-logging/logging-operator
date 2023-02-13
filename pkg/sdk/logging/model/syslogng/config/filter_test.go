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

package config

import (
	"strings"
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/stretchr/testify/require"
)

func TestFilterExpr(t *testing.T) {
	tests := map[string]struct {
		expr    model.FilterExpr
		wantOut string
		wantErr any
	}{
		"empty expr": {
			expr:    nil,
			wantErr: true,
		},
		"empty match expr": {
			expr:    model.NewFilterExpr(model.FilterExprMatch{}),
			wantOut: `match("")`,
		},
		"match expr with template": {
			expr: model.NewFilterExpr(model.FilterExprMatch{
				Pattern: "^foo",
				Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeTemplate("${HOST}|${MESSAGE}")),
			}),
			wantOut: `match("^foo" template("${HOST}|${MESSAGE}"))`,
		},
		"match expr with all options": {
			expr: model.NewFilterExpr(model.FilterExprMatch{
				Pattern: "^foo",
				Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("MESSAGE")),
				Flags: []string{
					"utf-8",
					"global",
				},
				Type: "pcre",
			}),
			wantOut: `match("^foo" value("MESSAGE") type("pcre") flags("utf-8" "global"))`,
		},
		"and": {
			expr: model.NewFilterExpr(model.FilterExprAnd{
				model.NewFilterExpr(model.FilterExprMatch{
					Pattern: "^foo",
					Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("MESSAGE")),
				}),
				model.NewFilterExpr(model.FilterExprMatch{
					Pattern: "bar$",
					Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("MESSAGE")),
				}),
			}),
			wantOut: `(match("^foo" value("MESSAGE")) and match("bar$" value("MESSAGE")))`,
		},
		"not": {
			expr: model.NewFilterExpr(model.FilterExprNot{
				Expr: model.NewFilterExpr(model.FilterExprMatch{
					Pattern: "^foo",
					Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("MESSAGE")),
				}),
			}),
			wantOut: `(not match("^foo" value("MESSAGE")))`,
		},
		"or": {
			expr: model.NewFilterExpr(model.FilterExprOr{
				model.NewFilterExpr(model.FilterExprMatch{
					Pattern: "^foo",
					Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("MESSAGE")),
				}),
				model.NewFilterExpr(model.FilterExprMatch{
					Pattern: "bar$",
					Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("MESSAGE")),
				}),
			}),
			wantOut: `(match("^foo" value("MESSAGE")) or match("bar$" value("MESSAGE")))`,
		},
		"regexp": {
			expr: model.NewFilterExpr(model.FilterExprMatch{
				Pattern: "^foo",
				Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeTemplate("${HOST}|${MESSAGE}")),
			}),
			wantOut: `match("^foo" template("${HOST}|${MESSAGE}"))`,
		},
	}
	for name, testCase := range tests {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			b := strings.Builder{}
			err := filterExpr(testCase.expr)(render.RenderContext{
				Out:        &b,
				IndentWith: "    ",
			})
			CheckError(t, testCase.wantErr, err)
			if err == nil {
				require.Equal(t, testCase.wantOut, b.String())
			}
		})
	}
}

//func TestRegexpMatchExpr_RenderAsSyslogNGConfig(t *testing.T) {
//}
//
//func TestAndExpr_RenderAsSyslogNGConfig(t *testing.T) {
//	ctx := syslogng.Context{
//		Indent: "    ",
//	}
//	tests := map[string]struct {
//		expr    AndExpr
//		ctx     syslogng.Context
//		wantOut string
//		wantErr bool
//	}{
//		"empty expr": {
//			expr:    AndExpr{},
//			ctx:     ctx,
//			wantOut: `()`,
//		},
//		"singleton": {
//			expr: AndExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")))`,
//		},
//		"binary": {
//			expr: AndExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "bar$",
//						Value:   "MESSAGE",
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")) and match("bar$" value("MESSAGE")))`,
//		},
//		"multi": {
//			expr: AndExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "bar$",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "buz",
//						Value:   "MESSAGE",
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")) and match("bar$" value("MESSAGE")) and match("buz" value("MESSAGE")))`,
//		},
//		"nested": {
//			expr: AndExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Or: OrExpr{
//						MatchExpr{
//							Regexp: &RegexpMatchExpr{
//								Pattern: "bar$",
//								Value:   "MESSAGE",
//							},
//						},
//						MatchExpr{
//							Regexp: &RegexpMatchExpr{
//								Pattern: "buz$",
//								Value:   "MESSAGE",
//							},
//						},
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")) and (match("bar$" value("MESSAGE")) or match("buz$" value("MESSAGE"))))`,
//		},
//	}
//	for name, testCase := range tests {
//		testCase := testCase
//		t.Run(name, func(t *testing.T) {
//			b := strings.Builder{}
//			testCase.ctx.Out = &b
//			err := testCase.expr.RenderAsSyslogNGConfig(testCase.ctx)
//			if (err != nil) != testCase.wantErr {
//				t.Errorf("AndExpr.RenderAsSyslogNGConfig() error = %v, wantErr %v", err, testCase.wantErr)
//			}
//			require.Equal(t, testCase.wantOut, b.String())
//		})
//	}
//}
//
//func TestNotExpr_RenderAsSyslogNGConfig(t *testing.T) {
//	ctx := syslogng.Context{
//		Indent: "    ",
//	}
//	tests := map[string]struct {
//		expr    *NotExpr
//		ctx     syslogng.Context
//		wantOut string
//		wantErr bool
//	}{
//		"nil expr": {
//			expr:    (*NotExpr)(nil),
//			ctx:     ctx,
//			wantOut: ``,
//		},
//		"non-empty expr": {
//			expr: (*NotExpr)(
//				&MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//			),
//			ctx:     ctx,
//			wantOut: `(not match("^foo" value("MESSAGE")))`,
//		},
//		"nested": {
//			expr: (*NotExpr)(&MatchExpr{
//				And: AndExpr{
//					MatchExpr{
//						Regexp: &RegexpMatchExpr{
//							Pattern: "^foo",
//							Value:   "MESSAGE",
//						},
//					},
//					MatchExpr{
//						Or: OrExpr{
//							MatchExpr{
//								Regexp: &RegexpMatchExpr{
//									Pattern: "bar$",
//									Value:   "MESSAGE",
//								},
//							},
//							MatchExpr{
//								Regexp: &RegexpMatchExpr{
//									Pattern: "buz$",
//									Value:   "MESSAGE",
//								},
//							},
//						},
//					},
//				}}),
//			ctx:     ctx,
//			wantOut: `(not (match("^foo" value("MESSAGE")) and (match("bar$" value("MESSAGE")) or match("buz$" value("MESSAGE")))))`,
//		},
//	}
//	for name, testCase := range tests {
//		testCase := testCase
//		t.Run(name, func(t *testing.T) {
//			b := strings.Builder{}
//			testCase.ctx.Out = &b
//			err := testCase.expr.RenderAsSyslogNGConfig(testCase.ctx)
//			if (err != nil) != testCase.wantErr {
//				t.Errorf("NotExpr.RenderAsSyslogNGConfig() error = %v, wantErr %v", err, testCase.wantErr)
//			}
//			require.Equal(t, testCase.wantOut, b.String())
//		})
//	}
//}
//
//func TestOrExpr_RenderAsSyslogNGConfig(t *testing.T) {
//	ctx := syslogng.Context{
//		Indent: "    ",
//	}
//	tests := map[string]struct {
//		expr    OrExpr
//		ctx     syslogng.Context
//		wantOut string
//		wantErr bool
//	}{
//		"empty expr": {
//			expr:    OrExpr{},
//			ctx:     ctx,
//			wantOut: `()`,
//		},
//		"singleton": {
//			expr: OrExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")))`,
//		},
//		"binary": {
//			expr: OrExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "bar$",
//						Value:   "MESSAGE",
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")) or match("bar$" value("MESSAGE")))`,
//		},
//		"multi": {
//			expr: OrExpr{
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "^foo",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "bar$",
//						Value:   "MESSAGE",
//					},
//				},
//				MatchExpr{
//					Regexp: &RegexpMatchExpr{
//						Pattern: "buz",
//						Value:   "MESSAGE",
//					},
//				},
//			},
//			ctx:     ctx,
//			wantOut: `(match("^foo" value("MESSAGE")) or match("bar$" value("MESSAGE")) or match("buz" value("MESSAGE")))`,
//		},
//	}
//	for name, testCase := range tests {
//		testCase := testCase
//		t.Run(name, func(t *testing.T) {
//			b := strings.Builder{}
//			testCase.ctx.Out = &b
//			err := testCase.expr.RenderAsSyslogNGConfig(testCase.ctx)
//			if (err != nil) != testCase.wantErr {
//				t.Errorf("OrExpr.RenderAsSyslogNGConfig() error = %v, wantErr %v", err, testCase.wantErr)
//			}
//			require.Equal(t, testCase.wantOut, b.String())
//		})
//	}
//}
//
