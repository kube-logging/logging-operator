// Copyright © 2020 Banzai Cloud
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

package output

import (
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +kubebuilder:object:generate=true
// Documentation: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829124
type SyslogNGSyslogOutput struct {
	Host           string              `json:"host"`
	Port           int                 `json:"port,omitempty"`
	Transport      string              `json:"transport,omitempty"`
	CaDir          *secret.Secret      `json:"ca_dir,omitempty"`
	CaFile         *secret.Secret      `json:"ca_file,omitempty"`
	CloseOnInput   *bool               `json:"close_on_input,omitempty"`
	Flags          []string            `json:"flags,omitempty"`
	FlushLines     int                 `json:"flush_lines,omitempty"`
	SoKeepalive    *bool               `json:"so_keepalive,omitempty"`
	Suppress       int                 `json:"suppress,omitempty"`
	Template       string              `json:"template,omitempty"`
	TemplateEscape *bool               `json:"template_escape,omitempty"`
	TLS            *TLS                `json:"tls,omitempty"`
	TSFormat       string              `json:"ts_format,omitempty"`
	DiskBuffer     *SyslogNGDiskBuffer `json:"disk_buffer,omitempty"`
}

// +kubebuilder:object:generate=true
type TLS struct {
	//TODO
}

// +kubebuilder:object:generate=true
type SyslogNGDiskBuffer struct {
	DiskBufSize  int64  `json:"disk_buf_size"`
	Reliable     bool   `json:"reliable"`
	Compaction   *bool  `json:"compaction,omitempty"`
	Dir          string `json:"dir,omitempty"`
	MemBufLength *int64 `json:"mem_buf_length,omitempty"`
	MemBufSize   *int64 `json:"mem_buf_size,omitempty"`
	QOutSize     *int64 `json:"q_out_size,omitempty"`
}
