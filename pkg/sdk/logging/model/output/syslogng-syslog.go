package output

import (
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +kubebuilder:object:generate=true
// Documentation: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829124
type SyslogNGSyslogOutput struct {
	Host           string
	Port           int
	Transport      string
	CaDir          *secret.Secret
	CaFile         *secret.Secret
	CloseOnInput   *bool
	Flags          []string
	FlushLines     int
	SoKeepalive    *bool
	Suppress       int
	Template       string
	TemplateEscape *bool
	TLS            *TLS
	TSFormat       string
	DiskBuffer     *SyslogNGDiskBuffer
}

// +kubebuilder:object:generate=true
type TLS struct {
	//TODO
}

// +kubebuilder:object:generate=true
type SyslogNGDiskBuffer struct {
	DiskBufSize  float64
	Reliable     bool
	Compaction   *bool
	Dir          string
	MemBufLength int64
	MemBufSize   float64
	QOutSize     int64
}
