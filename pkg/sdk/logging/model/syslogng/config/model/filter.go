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

type FilterExpr interface {
	__FilterExpr_union()
}

type FilterExprAlts interface {
	// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/65#TOPIC-1829161
	// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/66#TOPIC-1829165
	FilterExprAnd | FilterExprMatch | FilterExprNot | FilterExprOr
}

func NewFilterExpr[Alt FilterExprAlts](alt Alt) FilterExpr {
	return FilterExprAlt[Alt]{
		Alt: alt,
	}
}

type FilterExprAlt[Alt FilterExprAlts] struct {
	Alt Alt
}

func (FilterExprAlt[Alt]) __FilterExpr_union() {}

type FilterExprAnd []FilterExpr

// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#TOPIC-1829171
type FilterExprMatch struct {
	Pattern string
	Scope   FilterExprMatchScope
	Type    string
	Flags   []string
}

type FilterExprNot struct {
	Expr FilterExpr
}

type FilterExprOr []FilterExpr

type FilterExprMatchScope interface {
	__FilterExprMatchScope_union()
}

func NewFilterExprMatchScope[Alt FilterExprMatchScopeAlts](alt Alt) FilterExprMatchScope {
	return FilterExprMatchScopeAlt[Alt]{
		Alt: alt,
	}
}

type FilterExprMatchScopeAlt[Alt FilterExprMatchScopeAlts] struct {
	Alt Alt
}

func (FilterExprMatchScopeAlt[Alt]) __FilterExprMatchScope_union() {}

type FilterExprMatchScopeAlts interface {
	FilterExprMatchScopeTemplate | FilterExprMatchScopeValue
}

type FilterExprMatchScopeTemplate string

type FilterExprMatchScopeValue string
