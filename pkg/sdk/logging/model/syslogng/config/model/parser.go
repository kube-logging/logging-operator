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

type ParserDef struct {
	Name    string
	Parsers []Parser
}

type Parser interface {
	__Parser_union()
	Name() string
}

type ParserAlts interface {
	JSONParser
	Name() string
}

func NewParser[Alt ParserAlts](alt Alt) Parser {
	return ParserAlt[Alt]{
		Alt: alt,
	}
}

type ParserAlt[Alt ParserAlts] struct {
	Alt Alt
}

func (ParserAlt[Alt]) __Parser_union() {}

func (alt ParserAlt[Alt]) Name() string {
	return alt.Alt.Name()
}

type JSONParser struct {
	ExtractPrefix string
	Marker        string
	Prefix        string
	Template      string
}

func (JSONParser) Name() string {
	return "json-parser"
}
