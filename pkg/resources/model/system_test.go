// Copyright © 2026 Kube logging authors
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

package model

import (
	"reflect"
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/common"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/input"
)

func TestForwardInputConfig(t *testing.T) {
	tlsEnabled := v1beta1.FluentdTLS{Enabled: true, SecretName: "fluentd-tls", SharedKey: "shared"}
	tlsSecurity := &common.Security{SelfHostname: "fluentd", SharedKey: "shared"}

	tests := []struct {
		name          string
		spec          *v1beta1.FluentdSpec
		wantTransport *common.Transport
		wantSecurity  *common.Security
		wantPort      string
	}{
		{
			name:     "nil spec",
			spec:     nil,
			wantPort: "",
		},
		{
			name: "tls disabled keeps user transport as is",
			spec: &v1beta1.FluentdSpec{
				ForwardInputConfig: &input.ForwardInputConfig{
					Port:      "24241",
					Transport: &common.Transport{Ciphers: "ECDHE-RSA-AES256-GCM-SHA384"},
				},
			},
			wantTransport: &common.Transport{Ciphers: "ECDHE-RSA-AES256-GCM-SHA384"},
			wantPort:      "24241",
		},
		{
			name: "tls enabled without forward input config sets paths and default version",
			spec: &v1beta1.FluentdSpec{TLS: tlsEnabled},
			wantTransport: &common.Transport{
				Version:        "TLSv1_2",
				CaPath:         "/fluentd/tls/ca.crt",
				CertPath:       "/fluentd/tls/tls.crt",
				PrivateKeyPath: "/fluentd/tls/tls.key",
				ClientCertAuth: true,
			},
			wantSecurity: tlsSecurity,
		},
		{
			name: "tls enabled with user transport without version defaults version",
			spec: &v1beta1.FluentdSpec{
				TLS: tlsEnabled,
				ForwardInputConfig: &input.ForwardInputConfig{
					Transport: &common.Transport{Ciphers: "ECDHE-RSA-AES256-GCM-SHA384"},
				},
			},
			wantTransport: &common.Transport{
				Version:        "TLSv1_2",
				Ciphers:        "ECDHE-RSA-AES256-GCM-SHA384",
				CaPath:         "/fluentd/tls/ca.crt",
				CertPath:       "/fluentd/tls/tls.crt",
				PrivateKeyPath: "/fluentd/tls/tls.key",
				ClientCertAuth: true,
			},
			wantSecurity: tlsSecurity,
		},
		{
			name: "tls enabled keeps user ciphers, version and protocol but forces cert paths",
			spec: &v1beta1.FluentdSpec{
				TLS: tlsEnabled,
				ForwardInputConfig: &input.ForwardInputConfig{
					Port: "24241",
					Transport: &common.Transport{
						Protocol:       "tcp",
						Version:        "TLSv1_3",
						Ciphers:        "ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384",
						CaPath:         "/somewhere/else/ca.crt",
						CertPath:       "/somewhere/else/tls.crt",
						PrivateKeyPath: "/somewhere/else/tls.key",
					},
				},
			},
			wantTransport: &common.Transport{
				Protocol:       "tcp",
				Version:        "TLSv1_3",
				Ciphers:        "ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384",
				CaPath:         "/fluentd/tls/ca.crt",
				CertPath:       "/fluentd/tls/tls.crt",
				PrivateKeyPath: "/fluentd/tls/tls.key",
				ClientCertAuth: true,
			},
			wantSecurity: tlsSecurity,
			wantPort:     "24241",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before *v1beta1.FluentdSpec
			if tt.spec != nil {
				before = tt.spec.DeepCopy()
			}

			got := forwardInputConfig(tt.spec)

			if !reflect.DeepEqual(got.Transport, tt.wantTransport) {
				t.Errorf("transport mismatch\n got: %+v\nwant: %+v", got.Transport, tt.wantTransport)
			}
			if !reflect.DeepEqual(got.Security, tt.wantSecurity) {
				t.Errorf("security mismatch\n got: %+v\nwant: %+v", got.Security, tt.wantSecurity)
			}
			if got.Port != tt.wantPort {
				t.Errorf("port mismatch: got %q, want %q", got.Port, tt.wantPort)
			}
			if tt.spec != nil && !reflect.DeepEqual(tt.spec, before) {
				t.Errorf("fluentd spec was mutated\n got: %+v\nwant: %+v", tt.spec, before)
			}
		})
	}
}
