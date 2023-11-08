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
	"fmt"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/siliconbrain/go-seqs/seqs"
)

func filterDefStmt(name string, body render.Renderer) render.Renderer {
	return braceDefStmt("filter", name, body)
}

func filterExprStmt(expr model.FilterExpr) render.Renderer {
	return render.Line(render.AllOf(filterExpr(expr), render.String(";")))
}

func filterExpr(expr model.FilterExpr) render.Renderer {
	switch expr := expr.(type) {
	case model.FilterExprAlt[model.FilterExprAnd]:
		return render.AllOf(
			render.String("("),
			render.AllFrom(seqs.Intersperse(seqs.Map(seqs.FromSlice(expr.Alt), filterExpr), render.String(" and "))),
			render.String(")"),
		)
	case model.FilterExprAlt[model.FilterExprMatch]:
		args := []render.Renderer{
			render.Quoted(expr.Alt.Pattern),
		}
		if expr.Alt.Scope != nil {
			switch scope := expr.Alt.Scope.(type) {
			case model.FilterExprMatchScopeAlt[model.FilterExprMatchScopeValue]:
				args = append(args, optionExpr("value", render.Literal(string(scope.Alt))))
			case model.FilterExprMatchScopeAlt[model.FilterExprMatchScopeTemplate]:
				args = append(args, optionExpr("template", render.Literal(string(scope.Alt))))
			}
		}
		if typ := expr.Alt.Type; typ != "" {
			args = append(args, optionExpr("type", render.Literal(typ)))
		}
		if flags := expr.Alt.Flags; len(flags) > 0 {
			args = append(args, optionExpr("flags", seqs.ToSlice(seqs.Map(seqs.FromSlice(flags), render.Literal[string]))...))
		}
		return render.AllOf(
			render.String("match("),
			render.SpaceSeparated(args...),
			render.String(")"),
		)
	case model.FilterExprAlt[model.FilterExprNot]:
		return render.AllOf(render.String("(not "), filterExpr(expr.Alt.Expr), render.String(")"))
	case model.FilterExprAlt[model.FilterExprOr]:
		return render.AllOf(
			render.String("("),
			render.AllFrom(seqs.Intersperse(seqs.Map(seqs.FromSlice(expr.Alt), filterExpr), render.String(" or "))),
			render.String(")"),
		)
	default:
		return render.Error(fmt.Errorf("unsupported filter expression %T", expr))
	}
}
