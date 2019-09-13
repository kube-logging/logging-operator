package output

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +kubebuilder:object:generate=true
// +docName:"Aliyun OSS"
type OSSOutput struct {
	// OSS endpoint to connect to'
	Endpoint string `json:"endpoint"`
	// Your bucket name
	Bucket string `json:"bucket"`
	// Your access key id
	// +docLink:"Secret,./secret.md"
	AccessKeyId *secret.Secret `json:"access_key_id"`
	// Your access secret key
	// +docLink:"Secret,./secret.md"
	AaccessKeySecret *secret.Secret `json:"aaccess_key_secret"`
	// Path prefix of the files on OSS (default: fluent/logs)
	Path string `json:"path,omitempty"`
	// Upload crc enabled (default: true)
	UploadCrcEnable bool `json:"upload_crc_enable,omitempty"`
	// Download crc enabled (default: true)
	DownloadCrcEnable bool `json:"download_crc_enable,omitempty"`
	// Timeout for open connections (default: 10)
	OpenTimeout int `json:"open_timeout,omitempty"`
	// Timeout for read response (default: 120)
	ReadTimeout int `json:"read_timeout,omitempty"`
	// OSS SDK log directory (default: /var/log/td-agent)
	OssSdkLogDir string `json:"oss_sdk_log_dir,omitempty"`
	// The format of OSS object keys (default: %{path}/%{time_slice}_%{index}_%{thread_id}.%{file_extension})
	KeyFormat string `json:"key_format,omitempty"`
	// Archive format on OSS: gzip, json, text, lzo, lzma2 (default: gzip)
	StoreAs string `json:"store_as,omitempty"`
	// desc 'Create OSS bucket if it does not exists (default: false)
	AutoCreateBucket bool `json:"auto_create_bucket,omitempty"`
	// Overwrite already existing path (default: false)
	Overwrite bool `json:"overwrite,omitempty"`
	// Check bucket if exists or not (default: true)
	CheckBucket bool `json:"check_bucket,omitempty"`
	// Check object before creation (default: true)
	CheckObject bool `json:"check_object,omitempty"`
	// The length of `%{hex_random}` placeholder(4-16) (default: 4)
	HexRandomLength int `json:"hex_random_length,omitempty"`
	// `sprintf` format for `%{index}` (default: %d)
	IndexFormat string `json:"index_format,omitempty"`
	// Given a threshold to treat events as delay, output warning logs if delayed events were put into OSS
	WarnForDelay string `json:"warn_for_delay,omitempty"`
	// +docLink:"Format,./format.md"
	Format *Format `json:"format,omitempty"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (o *OSSOutput) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	oss := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      "gcs",
			Directive: "match",
			Tag:       "**",
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(o); err != nil {
		return nil, err
	} else {
		oss.Params = params
	}
	if o.Buffer != nil {
		if buffer, err := o.Buffer.ToDirective(secretLoader); err != nil {
			return nil, err
		} else {
			oss.SubDirectives = append(oss.SubDirectives, buffer)
		}
	}
	if o.Format != nil {
		if format, err := o.Format.ToDirective(secretLoader); err != nil {
			return nil, err
		} else {
			oss.SubDirectives = append(oss.SubDirectives, format)
		}
	}
	return oss, nil
}
