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

var fluentBitConfigTemplate = `
[SERVICE]
    Flush        1
    Daemon       Off
    Log_Level    info
    Parsers_File parsers.conf
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

[FILTER]
    Name        kubernetes
    {{- range $key, $value := .Filter }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}

[OUTPUT]
    Name          forward
    Match         *
    Host          {{ .TargetHost }}
    Port          {{ .TargetPort }}
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
    Retry_Limit   False
`
