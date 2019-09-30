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

package fluentd

var fluentdDefaultTemplate = `
# include other config files
@include /fluentd/etc/input.conf
@include /fluentd/app-config/*
@include /fluentd/etc/devnull.conf
`
var fluentdInputTemplate = `
# Enable RPC endpoint (this allows to trigger config reload without restart)
<system>
  rpc_endpoint 127.0.0.1:24444
</system>
# Prometheus monitoring
{{ if .Monitor.Enabled }}
<source>
    @type prometheus
    port {{ .Monitor.Port }}
    metrics_path {{ .Monitor.Path }}
</source>
<source>
    @type prometheus_monitor
</source>
<source>
    @type prometheus_output_monitor
</source>
{{ end }}

# Prevent fluentd from handling records containing its own logs. Otherwise
# it can lead to an infinite loop, when error in sending one message generates
# another message which also fails to be sent and so on.
<match **.fluentd**>
    @type null
</match>

<match **.fluent-bit**>
    @type null
</match>

`
var fluentdOutputTemplate = `
<match **>
    @type null
</match>
`
