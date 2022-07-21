package output

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/render/syslogng"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

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

func (o SyslogNGSyslogOutput) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	var CaFile string
	var err error
	if o.CaFile != nil {
		CaFile, err = ctx.SecretLoader.Load(o.CaFile)
		if err != nil {
			return err
		}
	}

	return syslogng.AllOf(
		syslogng.String("syslog("),
		syslogng.Printf("%q ", o.Host),
		syslogng.RenderIf(o.Transport != "", syslogng.Printf("transport(%q) ", o.Transport)),
		syslogng.RenderIf(o.Port != 0, syslogng.Printf("port(%q) ", o.Port)),
		syslogng.RenderIf(CaFile != "", syslogng.Printf("ca-file(%q) ", CaFile)),
		syslogng.RenderIf(o.DiskBuffer != nil, o.DiskBuffer),
		syslogng.String(");"),
	).RenderAsSyslogNGConfig(ctx)
}

type TLS struct {
	//TODO
}

type SyslogNGDiskBuffer struct {
	Reliable     *bool
	Compaction   *bool
	Dir          string
	DiskBufSize  float64
	MemBufLength int64
	MemBufSize   float64
	QOutSize     int64
}

func (o SyslogNGDiskBuffer) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return nil
}
