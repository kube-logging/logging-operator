// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package nodeagent

const BaseConfigName = "fluent-bit.conf"
const UpstreamConfigName = "upstream.conf"

var fluentBitConfigTemplate = `
[SERVICE]
    Flush        {{ .Flush }}
    Grace        {{ .Grace }}
    Daemon       Off
    Log_Level    {{ .LogLevel }}
    Parsers_File /fluent-bit/etc/parsers.conf
    Coro_Stack_Size    {{ .CoroStackSize }}
    {{- if .Monitor.Enabled }}
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_Port    {{ .Monitor.Port }}
    {{- end }}

[INPUT]
    Name         tail
    {{- range $key, $value := .Input.Values }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
    {{- range $id, $v := .Input.ParserN }}
    {{- if $v }}
    Parse_{{ $id}} {{$v}}
    {{- end }}
    {{- end }}
    {{- if .Input.MultilineParser }}
    multiline.parser {{- range $i, $v := .Input.MultilineParser }}{{ if $i }},{{ end}} {{ $v }}{{ end }}
    {{- end }}

{{- if not .DisableKubernetesFilter }}
[FILTER]
    Name        kubernetes
    {{- range $key, $value := .KubernetesFilter }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
{{- end}}

{{- if .AwsFilter }}

[FILTER]
    Name        aws
    {{- range $key, $value := .AwsFilter }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
{{- end}}

[OUTPUT]
    Name          forward
    Match         *
    {{- if .Upstream.Enabled }}
    Upstream /fluent-bit/conf_upstream/upstream.conf
    {{- else }}
    Host          {{ .TargetHost }}
    Port          {{ .TargetPort }}
    {{- end }}
    {{ if .TLS.Enabled }}
    tls           On
    tls.verify    Off
    tls.ca_file   /fluent-bit/tls/ca.crt
    tls.crt_file  /fluent-bit/tls/tls.crt
    tls.key_file  /fluent-bit/tls/tls.key
    {{- if .TLS.SharedKey }}
    Shared_Key    {{ .TLS.SharedKey }}
    {{- else }}
    Empty_Shared_Key true
    {{- end }}
    {{- end }}
    {{- if .Network.ConnectTimeoutSet }}
    net.connect_timeout {{.Network.ConnectTimeout}}
    {{- end }}
    {{- if .Network.KeepaliveSet}}
    net.keepalive {{if .Network.Keepalive }}on{{else}}off{{end}}
    {{- end }}
    {{- if .Network.KeepaliveIdleTimeoutSet }}
    net.keepalive_idle_timeout {{.Network.KeepaliveIdleTimeout}}
    {{- end }}
    {{- if .Network.KeepaliveMaxRecycleSet  }}
    net.keepalive_max_recycle {{.Network.KeepaliveMaxRecycle}}
    {{- end }}
    {{- if .ForwardOptions }}
    {{- range $key, $value := .ForwardOptions }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
    {{- end }}
`

var upstreamConfigTemplate = `
[UPSTREAM]
    Name {{ .Config.Name }}
{{- range $idx, $element:= .Config.Nodes}}
[NODE]
    Name {{.Name}}
    Host {{.Host}}
    Port {{.Port}}
{{- end}}
`
