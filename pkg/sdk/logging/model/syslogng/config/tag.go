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

import (
	"reflect"
	"strings"
)

const syslogNGTagName = "syslog-ng"

func structFieldSettings(f reflect.StructField) syslogNGTagSettings {
	return parseSettingsFromTag(f.Tag.Get(syslogNGTagName))
}

func parseSettingsFromTag(tag string) (settings syslogNGTagSettings) {
	if tag == "" {
		return
	}
	settings = make(map[string]string)
	for _, setting := range strings.Split(tag, ",") {
		switch s := strings.SplitN(setting, "=", 2); len(s) {
		case 0: // ignore empty setting
		case 1:
			if s[0] != "" {
				settings[s[0]] = ""
			}
		case 2:
			settings[s[0]] = s[1]
		}
	}
	return
}

type syslogNGTagSettings map[string]string

func (s syslogNGTagSettings) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s syslogNGTagSettings) Ignore() bool {
	return s.Has("ignore")
}

func (s syslogNGTagSettings) Name() string {
	return s["name"]
}

func (s syslogNGTagSettings) Optional() bool {
	return s.Has("optional")
}

func jsonNameOf(f reflect.StructField) string {
	name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
	return name
}

func hasJSONOmitempty(f reflect.StructField) bool {
	return strings.Contains(f.Tag.Get("json"), ",omitempty")
}
