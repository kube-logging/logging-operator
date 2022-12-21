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

package fluentbit

const BaseConfigName = "fluent-bit.conf"
const UpstreamConfigName = "upstream.conf"

var fluentBitConfigTemplate = `
[SERVICE]
    Flush        {{ .Flush }}
    Grace        {{ .Grace }}
    Daemon       Off
    Log_Level    {{ .LogLevel }}
    Parsers_File parsers.conf
    Coro_Stack_Size    {{ .CoroStackSize }}
    {{- if .Monitor.Enabled }}
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_Port    {{ .Monitor.Port }}
    {{- end }}
    {{- range $key, $value := .BufferStorage }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
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

{{- range $modify := .FilterModify }}

[FILTER]
    Name modify
    Match *
    {{- range $condition := $modify.Conditions }}
    {{- $operation :=  $condition.Operation }}
    Condition {{ $operation.Op }} {{ $operation.Key }} {{ if $operation.Value }}{{ $operation.Value }}{{ end }}
    {{- end }}

    {{- range $rule := $modify.Rules }}
    {{- $operation :=  $rule.Operation }}
    {{ $operation.Op }} {{ $operation.Key }} {{ if $operation.Value }}{{ $operation.Value }}{{ end }}
    {{- end }}
{{- end}}

{{- with .FluentForwardOutput }}
[OUTPUT]
    Name          forward
    Match         *
    {{- if .Upstream.Enabled }}
    Upstream      upstream.conf
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
    {{- if .Network.ConnectTimeoutLogErrorSet }}
    net.connect_timeout_log_error {{.Network.ConnectTimeoutLogError}}
    {{- end }}
    {{- if .Network.DNSMode }}
    net.dns.mode {{.Network.DNSMode}}
    {{- end }}
    {{- if .Network.DNSPreferIPV4Set }}
    net.dns.prefer_ipv4 {{.Network.DNSPreferIPV4}}
    {{- end }}
    {{- if .Network.DNSResolver }}
    net.dns.resolver {{.Network.DNSResolver}}
    {{- end }}
    {{- if .Network.KeepaliveSet}}
    net.keepalive {{if .Network.Keepalive }}on{{else}}off{{end}}
    {{- end }}
    {{- if .Network.KeepaliveIdleTimeoutSet }}
    net.keepalive_idle_timeout {{.Network.KeepaliveIdleTimeout}}
    {{- end }}
    {{- if .Network.KeepaliveMaxRecycleSet }}
    net.keepalive_max_recycle {{.Network.KeepaliveMaxRecycle}}
    {{- end }}
    {{- if .Network.SourceAddress }}
    net.source_address {{.Network.SourceAddress}}
    {{- end }}
    {{- with .Options }}
    {{- range $key, $value := . }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
    {{- end }}
{{- end }}

{{- with .SyslogNGOutput }}
[OUTPUT]
    Name tcp
    Match *
    Host {{ .Host }}
    Port {{ .Port }}
    Format json_lines
    {{- with .JSONDateKey }}
    json_date_key {{ . }}
    {{- end }}
    {{- with .JSONDateFormat }}
    json_date_format {{ . }}
    {{- end }}
    {{- with .Workers }}
    Workers {{ . }}
    {{- end }}
    {{- if .Network.ConnectTimeoutSet }}
    net.connect_timeout {{.Network.ConnectTimeout}}
    {{- end }}
    {{- if .Network.ConnectTimeoutLogErrorSet }}
    net.connect_timeout_log_error {{.Network.ConnectTimeoutLogError}}
    {{- end }}
    {{- if .Network.DNSMode }}
    net.dns.mode {{.Network.DNSMode}}
    {{- end }}
    {{- if .Network.DNSPreferIPV4Set }}
    net.dns.prefer_ipv4 {{.Network.DNSPreferIPV4}}
    {{- end }}
    {{- if .Network.DNSResolver }}
    net.dns.resolver {{.Network.DNSResolver}}
    {{- end }}
    {{- if .Network.KeepaliveSet}}
    net.keepalive {{if .Network.Keepalive }}on{{else}}off{{end}}
    {{- end }}
    {{- if .Network.KeepaliveIdleTimeoutSet }}
    net.keepalive_idle_timeout {{.Network.KeepaliveIdleTimeout}}
    {{- end }}
    {{- if .Network.KeepaliveMaxRecycleSet }}
    net.keepalive_max_recycle {{.Network.KeepaliveMaxRecycle}}
    {{- end }}
    {{- if .Network.SourceAddress }}
    net.source_address {{.Network.SourceAddress}}
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
