// Copyright © 2022 Cisco Systems, Inc. and/or its affiliates
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

import "github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"

func braceDefStmt(kind string, name string, body render.Renderer) render.Renderer {
	return render.AllOf(
		render.Line(render.SpaceSeparated(render.String(kind), render.If(name != "", render.Quoted(name)), render.String("{"))),
		render.Indented(body),
		render.Line(render.String("};")),
	)
}

func parenDefStmt(kind string, args ...render.Renderer) render.Renderer {
	return render.Line(render.AllOf(optionExpr(kind, args...), render.String(";")))
}

func optionExpr(key string, args ...render.Renderer) render.Renderer {
	return render.AllOf(render.String(key), render.String("("), render.SpaceSeparated(args...), render.String(")"))
}
