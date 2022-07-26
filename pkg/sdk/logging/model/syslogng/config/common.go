// Copyright © 2020 Banzai Cloud
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
	"github.com/siliconbrain/go-seqs/seqs"
)

func versionStmt(version string) Renderer {
	return Line(Formatted("@version: %s", version))
}

func includeStmt(file string) Renderer {
	return Line(Formatted("@include %q", file))
}

func globalOptionsDefStmt(def globalOptionsDef) Renderer {
	return braceDefStmt("options", "", AllOf(
		If(def.StatsLevel != nil, parenDefStmt("stats-level", literal(*def.StatsLevel))),
		If(def.StatsFreq != nil, parenDefStmt("stats-freq", literal(*def.StatsFreq))),
	))
}

type globalOptionsDef struct {
	StatsLevel *int
	StatsFreq  *int
}

func flagsOption(flags []string) Renderer {
	return optionExpr("flags", flags...)
}

func optionExpr[V any](key string, values ...V) Renderer {
	return AllOf(String(key), String("("), SpaceSeparated(seqs.ToSlice(seqs.Map(seqs.Map(seqs.FromSlice(values), eraseType[V]), literal))...), String(")"))
}

func braceDefStmt(kind string, name string, body Renderer) Renderer {
	return AllOf(
		Line(AllOf(String(kind), Space, If(name != "", AllOf(String(name), Space)), String("{"))),
		Indented(body),
		Line(String("};")),
	)
}

func parenDefStmt(kind string, args ...Renderer) Renderer {
	return Line(AllOf(String(kind), String("("), SpaceSeparated(args...), String(");")))
}

func literal(v any) Renderer {
	switch v := v.(type) {
	case Renderer:
		return v
	case string:
		return Quoted(v)
	case bool:
		if v {
			return String("yes")
		}
		return String("no")
	default:
		return Formatted("%v", v)
	}
}

func eraseType[V any](v V) any {
	return v
}
