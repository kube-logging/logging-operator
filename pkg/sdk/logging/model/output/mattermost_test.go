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

package output_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestMattermost(t *testing.T) {
	CONFIG := []byte(`
webhook_url: 
  value: https://mattermost:8080/webhooks/j4nrh34brj43br34jrb
channel_id: rfrf4r4r4ad
message_color: "#D9B9C9"
message_title: Test message
enable_tls: false
`)

	expected := `
  <match **>
    @type mattermost
    @id test
    channel_id rfrf4r4r4ad
	enable_tls false
    message_color #D9B9C9
    message_title Test message
    webhook_url https://mattermost:8080/webhooks/j4nrh34brj43br34jrb
  </match>
`

	mattermost := &output.MattermostOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, mattermost))
	test := render.NewOutputPluginTest(t, mattermost)
	test.DiffResult(expected)
}
