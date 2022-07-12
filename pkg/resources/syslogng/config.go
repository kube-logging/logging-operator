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
#SyslogNGConfigCheckTemplate
@version: 3.37

options { stats-level(3); };


source main-input { network(transport(tcp) flags(no-parse) port(2000)); };
destination file { file("/tmp/valami" template("$(format-json json.* )\n")); };

log {


        source(main-input);
        parser { json-parser(prefix("json.")); };
        destination(file);
};
`

var SyslogNGDefaultTemplate = `
#SyslogNGDefaultTemplate
@version: 3.37

options { stats-level(3); };


source main-input { network(transport(tcp) flags(no-parse) port(2000)); };
destination file { file("/tmp/valami" template("$(format-json json.* )\n")); };

log {


        source(main-input);
        parser { json-parser(prefix("json.")); };
        destination(file);
};
`

var SyslogNGLog = `
##SyslogNGLog
@version: 3.37

options { stats-level(3); };


source main-input { network(transport(tcp) flags(no-parse) port(2000)); };
destination file { file("/tmp/valami" template("$(format-json json.* )\n")); };

log {


        source(main-input);
        parser { json-parser(prefix("json.")); };
        destination(file);
};
`
var SyslogNGInputTemplate = `
#SyslogNGInputTemplate
`
var SyslogNGOutputTemplate = `
#SyslogNGOutputTemplate
`
