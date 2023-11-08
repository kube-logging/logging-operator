// Copyright © 2019 Banzai Cloud
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

package filter_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestGeoIP(t *testing.T) {
	CONFIG := []byte(`
geoip2_database: /opt/geoip-city.mmdb
geoip_lookup_keys: remote_addr
records:
  - city: ${city.names.en["remote_addr"]}
    location_array: '''[${location.longitude["remote"]},${location.latitude["remote"]}]'''
    country: ${country.iso_code["remote_addr"]}
    country_name: ${country.names.en["remote_addr"]}
    postal_code:  ${postal.code["remote_addr"]}
`)
	expected := `
<filter **>
  @type geoip
  @id test
  geoip2_database /opt/geoip-city.mmdb
  geoip_lookup_keys remote_addr
  skip_adding_null_record true
  <record>
    city ${city.names.en["remote_addr"]}
    country ${country.iso_code["remote_addr"]}
    country_name ${country.names.en["remote_addr"]}
    location_array '[${location.longitude["remote"]},${location.latitude["remote"]}]'
    postal_code ${postal.code["remote_addr"]}
  </record>
</filter>
`
	parser := &filter.GeoIP{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestGeoIPKeepNull(t *testing.T) {
	CONFIG := []byte(`
geoip_lookup_keys: remote_addr
skip_adding_null_record: false
records:
  - city: ${city.names.en["remote_addr"]}
    location_array: '''[${location.longitude["remote"]},${location.latitude["remote"]}]'''
    country: ${country.iso_code["remote_addr"]}
    country_name: ${country.names.en["remote_addr"]}
    postal_code:  ${postal.code["remote_addr"]}
`)

	expected := `
<filter **>
  @type geoip
  @id test
  geoip_lookup_keys remote_addr
  skip_adding_null_record false
  <record>
    city ${city.names.en["remote_addr"]}
    country ${country.iso_code["remote_addr"]}
    country_name ${country.names.en["remote_addr"]}
    location_array '[${location.longitude["remote"]},${location.latitude["remote"]}]'
    postal_code ${postal.code["remote_addr"]}
  </record>
</filter>
`

	parser := &filter.GeoIP{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}
