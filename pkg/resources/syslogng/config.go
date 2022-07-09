// Copyright © 2022 Banzai Cloud
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

package syslogng

var SyslogNGConfigCheckTemplate = `
@version: 3.37
@include "scl.conf"
options {
  # enable or
    disable directory creation for destination files
  create_dirs(yes);

  # keep
    hostnames from source host
  keep_hostname(yes);

  # use ISO8601 timestamps

    \ ts_format(iso);
};
`

var SyslogNGDefaultTemplate = `
@version: 3.27
@include "scl.conf"
options {
  # enable or
    disable directory creation for destination files
  create_dirs(yes);

  # keep
    hostnames from source host
  keep_hostname(yes);

  # use ISO8601 timestamps

    \ ts_format(iso);
};
`
var SyslogNGInputTemplate = `
log {
	source {
		network(transport(tcp) flags(no-parse));
	};
	destination { 
		file("/var/log/--syslog"); 
	};
};
`
var SyslogNGOutputTemplate = `
#Output Config
`

var SyslogNGLog = `
#log config
`
