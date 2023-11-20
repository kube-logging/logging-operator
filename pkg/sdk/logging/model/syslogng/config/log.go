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
	"reflect"

	"github.com/siliconbrain/go-seqs/seqs"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
)

type refScope string

const (
	Local  refScope = "local"
	Global refScope = "global"
)

type destination struct {
	logging          string
	renderedDestName string
	namespace        string
	name             string
	scope            refScope
	metricsProbes    []filter.MetricsProbe
}

func logDefStmt(sourceRefs []string, transforms []render.Renderer, destRefs []destination) render.Renderer {
	return braceDefStmt("log", "", render.AllOf(
		render.AllFrom(seqs.Map(seqs.FromSlice(sourceRefs), sourceRefStmt)),
		render.AllOf(transforms...),
		render.AllFrom(seqs.Map(seqs.FromSlice(destRefs), destinationLogPath)),
	))
}

func sourceRefStmt(name string) render.Renderer {
	return parenDefStmt("source", render.Literal(name))
}

func filterRefStmt(name string) render.Renderer {
	return parenDefStmt("filter", render.Literal(name))
}

func destinationLogPath(dest destination) render.Renderer {
	if len(dest.metricsProbes) == 0 {
		return braceDefStmt("log", "", render.AllOf(
			parenDefStmt("destination", render.Literal(dest.renderedDestName))),
		)
	}
	metricsProbesRenderer := make([]render.Renderer, len(dest.metricsProbes))
	for _, m := range dest.metricsProbes {
		if m.Labels == nil {
			m.Labels = make(filter.ArrowMap)
		}
		m.Labels["output_name"] = dest.name
		m.Labels["output_namespace"] = dest.namespace
		m.Labels["output_scope"] = string(dest.scope)
		m.Labels["logging"] = dest.logging
		metricsProbesRenderer = append(metricsProbesRenderer, renderDriver(Field{
			Value: reflect.ValueOf(m),
		}, nil))
	}

	return braceDefStmt("log", "", render.AllOf(
		parserDefStmt("", render.AllOf(metricsProbesRenderer...)),
		parenDefStmt("destination", render.Literal(dest.renderedDestName))),
	)
}
