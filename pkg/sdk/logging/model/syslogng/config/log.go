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

func logDefStmt(def model.LogDef) Renderer {
	return braceDefStmt("log", "", AllOf(
		AllFrom(seqs.Map(seqs.FromSlice(def.SourceNames), sourceRefStmt)),
		AllFrom(seqs.Map(seqs.FromSlice(def.OptionalElements), optionalLogElement)),
		AllFrom(seqs.Map(seqs.FromSlice(def.DestinationNames), destinationRefStmt)),
	))
}

func sourceRefStmt(name string) Renderer {
	return parenDefStmt("source", literal(name))
}

func destinationRefStmt(name string) Renderer {
	return parenDefStmt("destination", literal(name))
}

func optionalLogElement(e model.LogElement) Renderer {
	switch e := e.(type) {
	case model.LogElementAlt[model.FilterDef]:
		return filterDefStmt(e.Alt)
	case model.LogElementAlt[model.ParserDef]:
		return parserDefStmt(e.Alt)
	case model.LogElementAlt[model.RewriteDef]:
		return rewriteDefStmt(e.Alt)
	default:
		return Error(fmt.Errorf("unsupported optional log element type %T", e))
	}
}
