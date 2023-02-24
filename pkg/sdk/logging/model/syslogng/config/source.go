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

type NetworkSourceDriver struct {
	__meta         struct{} `syslog-ng:"name=network"` //lint:ignore U1000 field used for adding tag to the type
	Flags          []string `syslog-ng:"name=flags,optional"`
	IP             string   `syslog-ng:"name=ip,optional"`
	Port           uint16   `syslog-ng:"name=port,optional"`
	Transport      string   `syslog-ng:"name=transport,optional"`
	MaxConnections int      `syslog-ng:"name=max-connections,optional"`
	LogIWSize      int      `syslog-ng:"name=log-iw-size,optional"`
}

func sourceDefStmt(name string, body render.Renderer) render.Renderer {
	return braceDefStmt("source", name, body)
}
