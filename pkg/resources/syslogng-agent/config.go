// Copyright © 2023 Kube logging authors
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

package syslogng_agent

type syslogNGConfig struct {
	TargetHost string
	TargetPort int32
	TLS        struct {
		Enabled   bool
		SharedKey string
	}
}

var syslogNGConfigTemplate = `
@version: 4.0
@include "scl.conf"

# Define the source for log messages
source s_file {
  file("/var/log/messages"
       flags(no-parse)
       keep-timestamp(yes)
       log-fetch-limit(10000)
       log-iw-size(100000)
       log-msg-size(65535)
       pad-size(2048)
       follow-freq(1)
  );
};

# Define the destination for log messages with TLS encryption
destination d_tcp {
  network("{{ .TargetHost }}" port({{ .TargetPort }}) 
	{{ if .TLS.Enabled }}
	transport("tls")
    tls(ca-dir("/etc/syslog-ng/ca.d") cert-file("/etc/syslog-ng/cert.pem")
    key-file("/etc/syslog-ng/key.pem") peer-verify(optional-untrusted) 
	)
	{{- end }}
  );
};

# Define the log path for this configuration
log {
  source(s_file);
  destination(d_tcp);
};
`
