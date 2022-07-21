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

func parserDefStmt(def model.ParserDef) Renderer {
	return braceDefStmt("parser", def.Name, AllFrom(seqs.Map(seqs.FromSlice(def.Parsers), parserStmt)))
}

func parserStmt(parser model.Parser) Renderer {
	var args []Renderer
	switch parser := parser.(type) {
	case model.ParserAlt[model.JSONParser]:
		if parser.Alt.ExtractPrefix != "" {
			args = append(args, optionExpr("extract-prefix", parser.Alt.ExtractPrefix))
		}
		if parser.Alt.Marker != "" {
			args = append(args, optionExpr("marker", parser.Alt.Marker))
		}
		if parser.Alt.Prefix != "" {
			args = append(args, optionExpr("prefix", parser.Alt.Prefix))
		}
		if parser.Alt.Template != "" {
			args = append(args, optionExpr("template", parser.Alt.Template))
		}
	default:
		return Error(fmt.Errorf("unsupported parser type %q", parser.Name()))
	}
	return parenDefStmt(parser.Name(), args...)
}
