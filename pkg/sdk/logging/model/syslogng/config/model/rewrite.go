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

package model

type RewriteDef struct {
	Name  string
	Rules []RewriteRule
}

type RewriteRule interface {
	__RewriteRule_union()
	Name() string
}

type RewriteRuleAlts interface {
	RenameRule | SetRule | SubstituteRule | UnsetRule
	Name() string
}

func NewRewriteRule[Alt RewriteRuleAlts](alt Alt) RewriteRule {
	return RewriteRuleAlt[Alt]{
		Alt: alt,
	}
}

type RewriteRuleAlt[Alt RewriteRuleAlts] struct {
	Alt Alt
}

func (RewriteRuleAlt[Alt]) __RewriteRule_union() {}

func (alt RewriteRuleAlt[Alt]) Name() string {
	return alt.Alt.Name()
}

type RenameRule struct {
	OldFieldName string
	NewFieldName string
	Condition    *RewriteCondition
}

func (RenameRule) Name() string {
	return "rename"
}

type SetRule struct {
	FieldName string
	Value     string
	Condition *RewriteCondition
}

func (SetRule) Name() string {
	return "set"
}

type SubstituteRule struct {
	FieldName   string
	Pattern     string
	Replacement string
	Type        string
	Flags       []string
	Condition   *RewriteCondition
}

func (SubstituteRule) Name() string {
	return "subst"
}

type UnsetRule struct {
	FieldName string
	Condition *RewriteCondition
}

func (UnsetRule) Name() string {
	return "unset"
}

type RewriteCondition struct {
	Expr FilterExpr
}
