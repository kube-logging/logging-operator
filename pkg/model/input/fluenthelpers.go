package input

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

type Security struct {
	// Hostname
	SelfHostname string `json:"self_hostname"`
	// Shared key for authentication.
	SharedKey string `json:"shared_key"`
	// If true, use user based authentication.
	UserAuth bool `json:"user_auth,omitempty"`
	// Allow anonymous source. <client> sections are required if disabled.
	AllowAnonymousSource bool `json:"allow_anonymous_source,omitempty"`
}

func (s *Security) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	return types.NewFlatDirective(types.PluginMeta{
		Directive: "security",
	}, s, secretLoader)
}

type Transport struct {
	// Protocol Default: :tcp
	Protocol string `json:"protocol,omitempty"`
	// Version Default: 'TLSv1_2'
	Version string `json:"version,omitempty"`
	// Ciphers Default: "ALL:!aNULL:!eNULL:!SSLv2"
	Ciphers string `json:"ciphers,omitempty"`
	// Use secure connection when use tls) Default: false
	Insecure bool `json:"insecure,omitempty"`
	// Specify path to CA certificate file
	CaPath string `json:"ca_path,omitempty"`
	// Specify path to Certificate file
	CertPath string `json:"cert_path,omitempty"`
	// Specify path to private Key file
	PrivateKeyPath string `json:"private_key_path,omitempty"`
	// public CA private key passphrase contained path
	PrivateKeyPassphrase string `json:"private_key_passphrase,omitempty"`
	// When this is set Fluentd will check all incoming HTTPS requests
	// for a client certificate signed by the trusted CA, requests that
	// don't supply a valid client certificate will fail.
	ClientCertAuth bool `json:"client_cert_auth,omitempty"`
	// Specify private CA contained path
	CaCertPath string `json:"ca_cert_path,omitempty"`
	// private CA private key contained path
	CaPrivateKeyPath string `json:"ca_private_key_path,omitempty"`
	// private CA private key passphrase contained path
	CaPrivateKeyPassphrase string `json:"ca_private_key_passphrase,omitempty"`
}

func (t *Transport) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	return types.NewFlatDirective(types.PluginMeta{
		Directive: "transport",
		Tag:       "tls",
	}, t, secretLoader)
}
