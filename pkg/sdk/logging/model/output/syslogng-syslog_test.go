package output

import (
	"strings"
	"testing"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/render/syslogng"
	"github.com/stretchr/testify/require"
)

func TestSyslogOutputConfig_RenderAsSyslogNGConfig(t *testing.T) {

	ctx := syslogng.Context{
		Indent:    "    ",
		Namespace: "default",
	}
	tests := map[string]struct {
		cfg     SyslogNGSyslogOutput
		ctx     syslogng.Context
		wantOut string
		wantErr bool
	}{
		"empty expr": {
			cfg: SyslogNGSyslogOutput{
				Host:      "test",
				Transport: "tcp",
			},
			ctx:     ctx,
			wantOut: `syslog("test" transport("tcp") );`,
		},
	}
	for name, testCase := range tests {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			b := strings.Builder{}
			testCase.ctx.Out = &b
			err := testCase.cfg.RenderAsSyslogNGConfig(testCase.ctx)
			if (err != nil) != testCase.wantErr {
				t.Errorf("MatchConfig.RenderAsSyslogNGConfig() error = %v, wantErr %v", err, testCase.wantErr)
			}
			require.Equal(t, testCase.wantOut, b.String())
		})
	}
}
