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
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
)

type JSONParser struct {
	__meta        struct{} `syslog-ng:"name=json-parser"` //lint:ignore U1000 field used for adding tag to the type
	ExtractPrefix string   `syslog-ng:"name=extract-prefix,optional"`
	Marker        string   `syslog-ng:"name=marker,optional"`
	Prefix        string   `syslog-ng:"name=prefix,optional"`
	Template      string   `syslog-ng:"name=template,optional"`
	KeyDelimiter  string   `syslog-ng:"name=key-delimiter,optional"`
}

func parserDefStmt(name string, body render.Renderer) render.Renderer {
	return braceDefStmt("parser", name, body)
}

func isActiveParserDriver(f Field) bool {
	return hasParserDriverTag(f) && isActiveField(f)
}

func hasParserDriverTag(f Field) bool {
	return structFieldSettings(f.Meta).Has("parser-drv")
}

type DateParser struct {
	__meta   struct{} `syslog-ng:"name=date-parser"` //lint:ignore U1000 field used for adding tag to the type
	Format   string   `syslog-ng:"name=format,optional"`
	Template string   `syslog-ng:"name=template,optional"`
}
