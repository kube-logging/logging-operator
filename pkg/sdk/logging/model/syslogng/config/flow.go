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
	"reflect"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	filter "github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/siliconbrain/go-seqs/seqs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func renderClusterFlow(sourceName string, f v1beta1.SyslogNGClusterFlow, secretLoaderFactory SecretLoaderFactory) render.Renderer {
	return logDefStmt(
		[]string{sourceName},
		seqs.ToSlice(seqs.Concat(
			seqs.FromValues(
				renderFlowMatch(f.Spec.Match),
			),
			seqs.MapWithIndex(seqs.FromSlice(f.Spec.Filters), func(idx int, flt v1beta1.SyslogNGFilter) render.Renderer {
				return renderFlowFilter(flt, &f, idx, secretLoaderFactory.SecretLoaderForNamespace(f.Namespace))
			}),
		)),
		seqs.ToSlice(seqs.Map(seqs.FromSlice(f.Spec.GlobalOutputRefs), func(ref string) string { return clusterOutputDestName(f.Namespace, ref) })),
	)
}

func renderFlow(controlNS string, sourceName string, f v1beta1.SyslogNGFlow, secretLoaderFactory SecretLoaderFactory) render.Renderer {
	return logDefStmt(
		[]string{sourceName},
		seqs.ToSlice(seqs.Concat(
			seqs.FromValues(
				filterDefStmt("", filterExprStmt(model.NewFilterExpr(model.FilterExprMatch{
					Pattern: f.Namespace,
					Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue("json.kubernetes.namespace_name")),
					Type:    "string",
				}))),
				renderFlowMatch(f.Spec.Match),
			),
			seqs.MapWithIndex(seqs.FromSlice(f.Spec.Filters), func(idx int, flt v1beta1.SyslogNGFilter) render.Renderer {
				return renderFlowFilter(flt, &f, idx, secretLoaderFactory.SecretLoaderForNamespace(f.Namespace))
			}),
		)),
		seqs.ToSlice(seqs.Concat(
			seqs.Map(seqs.FromSlice(f.Spec.GlobalOutputRefs), func(ref string) string { return clusterOutputDestName(f.Namespace, ref) }),
			seqs.Map(seqs.FromSlice(f.Spec.LocalOutputRefs), func(ref string) string { return outputDestName(f.Namespace, ref) }),
		)),
	)
}

func renderFlowMatch(m *v1beta1.SyslogNGMatch) render.Renderer {
	if m == nil {
		return nil
	}
	return filterDefStmt("", renderMatchExpr(filter.MatchExpr(*m)))
}

func renderFlowFilter(flt v1beta1.SyslogNGFilter, flow metav1.Object, index int, secretLoader secret.SecretLoader) render.Renderer {
	xformFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(reflect.ValueOf(flt))), isActiveTransform))
	switch len(xformFields) {
	case 0:
		return render.Error(fmt.Errorf("no transformation specified on filter %d of flow %s/%s", index, flow.GetNamespace(), flow.GetName()))
	case 1:
		xformField := xformFields[0]
		settings := structFieldSettings(xformField.Meta)
		switch xformKind := settings[xformKindKey]; xformKind {
		case "filter":
			val := derefAll(xformField.Value)
			if !val.CanConvert(matchExprType) {
				return render.Error(fmt.Errorf("value of type %s is not a valid filter expression", xformField.Value.Type()))
			}
			return filterDefStmt("", filterExprStmt(filterExprFromMatchExpr(val.Convert(matchExprType).Interface().(filter.MatchExpr))))
		case "parser":
			driverFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(xformField.Value)), isActiveParserDriver))
			switch len(driverFields) {
			case 0:
				return render.Error(fmt.Errorf(
					"no parser driver specified on parser %s of filter %d of flow %s/%s",
					xformField.KeyOrEmpty(), index, flow.GetNamespace(), flow.GetName(),
				))
			case 1:
				return parserDefStmt("", renderDriver(driverFields[0], secretLoader))
			default:
				return render.Error(fmt.Errorf(
					"multiple parser drivers (%v) specified on parser %s of filter %d of flow %s/%s",
					seqs.ToSlice(seqs.Map(seqs.FromSlice(driverFields), Field.KeyOrEmpty)),
					xformField.KeyOrEmpty(), index, flow.GetNamespace(), flow.GetName(),
				))
			}
		case "rewrite":
			driverFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(xformField.Value)), isActiveRewriteDriver))
			switch len(driverFields) {
			case 0:
				return render.Error(fmt.Errorf(
					"no rewrite driver specified on parser %s of filter %d of flow %s/%s",
					xformField.KeyOrEmpty(), index, flow.GetNamespace(), flow.GetName(),
				))
			case 1:
				return rewriteDefStmt("", renderDriver(driverFields[0], secretLoader))
			default:
				return render.Error(fmt.Errorf(
					"multiple rewrite drivers (%v) specified on parser %s of filter %d of flow %s/%s",
					seqs.ToSlice(seqs.Map(seqs.FromSlice(driverFields), Field.KeyOrEmpty)),
					xformField.KeyOrEmpty(), index, flow.GetNamespace(), flow.GetName(),
				))
			}
		default:
			return render.Error(fmt.Errorf("unsupported transformation kind %q", xformKind))
		}
	default:
		return render.Error(fmt.Errorf(
			"multiple transformations (%v) specified on filter %d of flow %s/%s",
			seqs.ToSlice(seqs.Map(seqs.FromSlice(xformFields), Field.KeyOrEmpty)),
			index, flow.GetNamespace(), flow.GetName(),
		))
	}
}

func renderMatchExpr(expr filter.MatchExpr) render.Renderer {
	return filterExprStmt(filterExprFromMatchExpr(expr))
}

func filterExprFromMatchExpr(expr filter.MatchExpr) model.FilterExpr {
	switch {
	case len(expr.And) > 0:
		return model.NewFilterExpr(model.FilterExprAnd(seqs.ToSlice(seqs.Map(seqs.FromSlice(expr.And), filterExprFromMatchExpr))))
	case expr.Not != nil:
		return model.NewFilterExpr(model.FilterExprNot{Expr: filterExprFromMatchExpr(filter.MatchExpr(*expr.Not))})
	case len(expr.Or) > 0:
		return model.NewFilterExpr(model.FilterExprOr(seqs.ToSlice(seqs.Map(seqs.FromSlice(expr.Or), filterExprFromMatchExpr))))
	case expr.Regexp != nil:
		m := model.FilterExprMatch{
			Pattern: expr.Regexp.Pattern,
			Type:    expr.Regexp.Type,
			Flags:   expr.Regexp.Flags,
		}
		switch {
		case expr.Regexp.Template != "":
			m.Scope = model.NewFilterExprMatchScope(model.FilterExprMatchScopeTemplate(expr.Regexp.Template))
		case expr.Regexp.Value != "":
			m.Scope = model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue(expr.Regexp.Value))
		}
		return model.NewFilterExpr(m)
	default:
		return nil
	}
}

func isActiveTransform(f Field) bool {
	return isTransform(f) && isActiveField(f)
}

func isTransform(f Field) bool {
	return structFieldSettings(f.Meta).Has(xformKindKey)
}

const xformKindKey = "xform-kind"
