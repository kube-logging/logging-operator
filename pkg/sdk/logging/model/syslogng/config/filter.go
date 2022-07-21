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

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
	"github.com/siliconbrain/go-seqs/seqs"
)

func filterDefStmt(def model.FilterDef) Renderer {
	return braceDefStmt("filter", def.Name, filterExprStmt(def.Expr))
}

func filterExprStmt(expr model.FilterExpr) Renderer {
	return Line(AllOf(filterExpr(expr), String(";")))
}

func filterExpr(expr model.FilterExpr) Renderer {
	switch expr := expr.(type) {
	case model.FilterExprAlt[model.FilterExprAnd]:
		return AllOf(
			String("("),
			AllFrom(seqs.Intersperse(seqs.Map(seqs.FromSlice(expr.Alt), filterExpr), String(" and "))),
			String(")"),
		)
	case model.FilterExprAlt[model.FilterExprMatch]:
		args := []Renderer{
			Quoted(expr.Alt.Pattern),
		}
		if expr.Alt.Scope != nil {
			switch scope := expr.Alt.Scope.(type) {
			case model.FilterExprMatchScopeAlt[model.FilterExprMatchScopeValue]:
				args = append(args, optionExpr("value", string(scope.Alt)))
			case model.FilterExprMatchScopeAlt[model.FilterExprMatchScopeTemplate]:
				args = append(args, optionExpr("template", string(scope.Alt)))
			}
		}
		if typ := expr.Alt.Type; typ != "" {
			args = append(args, optionExpr("type", typ))
		}
		if flags := expr.Alt.Flags; len(flags) > 0 {
			args = append(args, flagsOption(flags))
		}
		return AllOf(
			String("match("),
			SpaceSeparated(args...),
			String(")"),
		)
	case model.FilterExprAlt[model.FilterExprNot]:
		return AllOf(String("(not "), filterExpr(expr.Alt.Expr), String(")"))
	case model.FilterExprAlt[model.FilterExprOr]:
		return AllOf(
			String("("),
			AllFrom(seqs.Intersperse(seqs.Map(seqs.FromSlice(expr.Alt), filterExpr), String(" or "))),
			String(")"),
		)
	default:
		return Error(fmt.Errorf("unsupported filter expression %T", expr))
	}
}
