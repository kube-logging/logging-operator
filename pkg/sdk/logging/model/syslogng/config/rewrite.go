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

func rewriteDefStmt(def model.RewriteDef) Renderer {
	return braceDefStmt("rewrite", def.Name, AllFrom(seqs.Map(seqs.FromSlice(def.Rules), rewriteRuleStmt)))
}

func rewriteRuleStmt(rule model.RewriteRule) Renderer {
	var args []Renderer
	switch rule := rule.(type) {
	case model.RewriteRuleAlt[model.RenameRule]:
		args = append(args, Quoted(rule.Alt.OldFieldName), Quoted(rule.Alt.NewFieldName))
		if rule.Alt.Condition != nil {
			args = append(args, rewriteCondition(*rule.Alt.Condition))
		}
	case model.RewriteRuleAlt[model.SetRule]:
		args = append(args, Quoted(rule.Alt.Value), optionExpr("value", rule.Alt.FieldName))
		if rule.Alt.Condition != nil {
			args = append(args, rewriteCondition(*rule.Alt.Condition))
		}
	case model.RewriteRuleAlt[model.SubstituteRule]:
		args = append(args, Quoted(rule.Alt.Pattern), Quoted(rule.Alt.Replacement))
		if rule.Alt.Type != "" {
			args = append(args, optionExpr("type", rule.Alt.Type))
		}
		if flags := rule.Alt.Flags; len(flags) > 0 {
			args = append(args, flagsOption(flags))
		}
		if rule.Alt.Condition != nil {
			args = append(args, rewriteCondition(*rule.Alt.Condition))
		}
	case model.RewriteRuleAlt[model.UnsetRule]:
		args = append(args, optionExpr("value", rule.Alt.FieldName))
		if rule.Alt.Condition != nil {
			args = append(args, rewriteCondition(*rule.Alt.Condition))
		}
	default:
		return Error(fmt.Errorf("unsupported rewrite rule type %q", rule.Name()))
	}
	return parenDefStmt(rule.Name(), args...)
}

func rewriteCondition(cond model.RewriteCondition) Renderer {
	return AllOf(String("condition("), filterExpr(cond.Expr), String(")"))
}
