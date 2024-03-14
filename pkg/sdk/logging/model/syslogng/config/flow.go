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
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/siliconbrain/go-seqs/seqs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	filter "github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
)

func validateClusterOutputs(clusterOutputRefs map[string]types.NamespacedName, flow string, globalOutputRefs []string) error {
	return seqs.SeededReduce(seqs.FromSlice(globalOutputRefs), nil, func(err error, ref string) error {
		if _, ok := clusterOutputRefs[ref]; !ok {
			return errors.Append(err, errors.Errorf("cluster output reference %s for flow %s cannot be found", ref, flow))
		}
		return err
	})
}

func renderClusterFlow(logging string, clusterOutputRefs map[string]types.NamespacedName, sourceName string, f v1beta1.SyslogNGClusterFlow, secretLoaderFactory SecretLoaderFactory) render.Renderer {
	baseName := fmt.Sprintf("clusterflow_%s_%s", f.Namespace, f.Name)
	matchName := fmt.Sprintf("%s_match", baseName)
	filterDefs := seqs.MapWithIndex(seqs.FromSlice(f.Spec.Filters), func(idx int, flt v1beta1.SyslogNGFilter) render.Renderer {
		return renderFlowFilter(flt, &f, idx, baseName, secretLoaderFactory.SecretLoaderForNamespace(f.Namespace))
	})
	return render.AllOf(
		renderFlowMatch(matchName, f.Spec.Match),
		render.AllFrom(filterDefs),
		logDefStmt(
			[]string{sourceName},
			seqs.ToSlice(seqs.Concat(
				seqs.FromValues(
					render.If(!f.Spec.Match.IsEmpty(), filterRefStmt(matchName)),
				),
				seqs.MapWithIndex(seqs.FromSlice(f.Spec.Filters), func(idx int, flt v1beta1.SyslogNGFilter) render.Renderer {
					return parenDefStmt(filterKind(flt), render.Literal(filterID(flt, idx, baseName)))
				}),
			)),
			seqs.ToSlice(seqs.Map(seqs.FromSlice(f.Spec.GlobalOutputRefs), func(ref string) destination {
				return destination{
					logging:          logging,
					renderedDestName: clusterOutputDestName(clusterOutputRefs[ref].Namespace, ref),
					namespace:        clusterOutputRefs[ref].Namespace,
					name:             ref,
					scope:            Global,
					metricsProbes:    f.Spec.OutputMetrics,
				}
			})),
		),
	)
}

func renderFlow(logging string, clusterOutputRefs map[string]types.NamespacedName, sourceName string, keyDelim string, f v1beta1.SyslogNGFlow, secretLoaderFactory SecretLoaderFactory) render.Renderer {
	baseName := fmt.Sprintf("flow_%s_%s", f.Namespace, f.Name)
	matchName := fmt.Sprintf("%s_match", baseName)
	nsFilterName := fmt.Sprintf("%s_ns_filter", baseName)
	filterDefs := render.AllFrom(seqs.MapWithIndex(seqs.FromSlice(f.Spec.Filters), func(idx int, flt v1beta1.SyslogNGFilter) render.Renderer {
		return renderFlowFilter(flt, &f, idx, baseName, secretLoaderFactory.SecretLoaderForNamespace(f.Namespace))
	}))
	return render.AllOf(
		filterDefStmt(nsFilterName, filterExprStmt(model.NewFilterExpr(model.FilterExprMatch{
			Pattern: f.Namespace,
			Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue(strings.Join([]string{"json", "kubernetes", "namespace_name"}, keyDelim))),
			Type:    "string",
		}))),
		renderFlowMatch(matchName, f.Spec.Match),
		filterDefs,
		logDefStmt(
			[]string{sourceName},
			seqs.ToSlice(seqs.Concat(
				seqs.FromValues(
					filterRefStmt(nsFilterName),
					render.If(!f.Spec.Match.IsEmpty(), filterRefStmt(matchName)),
				),
				seqs.MapWithIndex(seqs.FromSlice(f.Spec.Filters), func(idx int, flt v1beta1.SyslogNGFilter) render.Renderer {
					return parenDefStmt(filterKind(flt), render.Literal(filterID(flt, idx, baseName)))
				}),
			)),
			seqs.ToSlice(seqs.Concat(
				seqs.Map(seqs.FromSlice(f.Spec.GlobalOutputRefs), func(ref string) destination {
					return destination{
						logging:          logging,
						renderedDestName: clusterOutputDestName(clusterOutputRefs[ref].Namespace, ref),
						namespace:        clusterOutputRefs[ref].Namespace,
						name:             ref,
						scope:            Global,
						metricsProbes:    f.Spec.OutputMetrics,
					}
				}),
				seqs.Map(seqs.FromSlice(f.Spec.LocalOutputRefs), func(ref string) destination {
					return destination{
						logging:          logging,
						renderedDestName: outputDestName(f.Namespace, ref),
						namespace:        f.Namespace,
						name:             ref,
						scope:            Local,
						metricsProbes:    f.Spec.OutputMetrics,
					}
				}),
			)),
		),
	)
}

func renderFlowMatch(name string, m *v1beta1.SyslogNGMatch) render.Renderer {
	if m.IsEmpty() {
		return nil
	}
	return filterDefStmt(name, renderMatchExpr(filter.MatchExpr(*m)))
}

func renderFlowFilter(flt v1beta1.SyslogNGFilter, flow metav1.Object, index int, baseName string, secretLoader secret.SecretLoader) render.Renderer {
	filterID := filterID(flt, index, baseName)

	xformFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(reflect.ValueOf(flt))), isActiveTransform))
	switch len(xformFields) {
	case 0:
		return render.Error(fmt.Errorf("no transformation specified on filter %s of flow %s/%s", filterID, flow.GetNamespace(), flow.GetName()))
	case 1:
		xformField := xformFields[0]
		settings := structFieldSettings(xformField.Meta)
		switch xformKind := settings[xformKindKey]; xformKind {
		case "filter":
			val := derefAll(xformField.Value)
			if !val.CanConvert(matchExprType) {
				return render.Error(fmt.Errorf("value of type %s is not a valid filter expression", xformField.Value.Type()))
			}
			return filterDefStmt(filterID, filterExprStmt(filterExprFromMatchExpr(val.Convert(matchExprType).Interface().(filter.MatchExpr))))
		case "parser":
			driverFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(xformField.Value)), isActiveParserDriver))
			switch len(driverFields) {
			case 0:
				return render.Error(fmt.Errorf(
					"no parser driver specified on parser %s of filter %s of flow %s/%s",
					xformField.KeyOrEmpty(), filterID, flow.GetNamespace(), flow.GetName(),
				))
			case 1:
				return parserDefStmt(filterID, renderDriver(driverFields[0], secretLoader))
			default:
				return render.Error(fmt.Errorf(
					"multiple parser drivers (%v) specified on parser %s of filter %s of flow %s/%s",
					seqs.ToSlice(seqs.Map(seqs.FromSlice(driverFields), Field.KeyOrEmpty)),
					xformField.KeyOrEmpty(), filterID, flow.GetNamespace(), flow.GetName(),
				))
			}
		case "rewrite":
			switch xformField.Value.Kind() {
			case reflect.Array, reflect.Slice:
				var stmts []render.Renderer
				l := xformField.Value.Len()
				if l > 0 {
					stmts = make([]render.Renderer, l)
				}
				for i := 0; i < l; i++ {
					stmts[i] = renderRewriteDriver(xformField.Value.Index(i), xformField.KeyOrEmpty(), filterID, flow, secretLoader)
				}
				return rewriteDefStmt(filterID, render.AllOf(stmts...))
			default:
				return rewriteDefStmt(filterID, renderRewriteDriver(xformField.Value, xformField.KeyOrEmpty(), filterID, flow, secretLoader))
			}
		default:
			return render.Error(fmt.Errorf("unsupported transformation kind %q", xformKind))
		}
	default:
		return render.Error(fmt.Errorf(
			"multiple transformations (%v) specified on filter %s of flow %s/%s",
			seqs.ToSlice(seqs.Map(seqs.FromSlice(xformFields), Field.KeyOrEmpty)),
			filterID, flow.GetNamespace(), flow.GetName(),
		))
	}
}

func renderRewriteDriver(value reflect.Value, key string, filter string, flow metav1.Object, secretLoader secret.SecretLoader) render.Renderer {
	driverFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(value)), isActiveRewriteDriver))
	switch len(driverFields) {
	case 0:
		return render.Error(fmt.Errorf(
			"no rewrite driver specified on rewrite %s of filter %s of flow %s/%s",
			key, filter, flow.GetNamespace(), flow.GetName(),
		))
	case 1:
		return renderDriver(driverFields[0], secretLoader)
	default:
		return render.Error(fmt.Errorf(
			"multiple rewrite drivers (%v) specified on rewrite %s of filter %s of flow %s/%s",
			seqs.ToSlice(seqs.Map(seqs.FromSlice(driverFields), Field.KeyOrEmpty)),
			key, filter, flow.GetNamespace(), flow.GetName(),
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

func filterID(filter v1beta1.SyslogNGFilter, index int, baseName string) string {
	filterID := filter.ID
	if filterID == "" {
		filterID = strconv.FormatInt(int64(index), 10)
	}
	return baseName + "_filters_" + filterID
}

func filterKind(filter v1beta1.SyslogNGFilter) string {
	xformFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(reflect.ValueOf(filter))), isActiveTransform))
	if len(xformFields) != 1 {
		return ""
	}
	xformField := xformFields[0]
	settings := structFieldSettings(xformField.Meta)
	return settings[xformKindKey]
}

func isActiveTransform(f Field) bool {
	return isTransform(f) && isActiveField(f)
}

func isTransform(f Field) bool {
	return structFieldSettings(f.Meta).Has(xformKindKey)
}

const xformKindKey = "xform-kind"
