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

package filter

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +docName:"Fluentd GeoIP filter"
// Fluentd Filter plugin to add information about geographical location of IP addresses with Maxmind GeoIP databases.
// More information at https://github.com/y-ken/fluent-plugin-geoip
//
// #### Example record configurations
// ```
// spec:
//  filters:
//    - tag_normaliser:
//        format: ${namespace_name}.${pod_name}.${container_name}
//    - parser:
//        remove_key_name_field: true
//        parse:
//          type: nginx
//    - geoip:
//        geoip_lookup_keys: remote_addr
//        records:
//          - city: ${city.names.en["remote_addr"]}
//            location_array: '''[${location.longitude["remote"]},${location.latitude["remote"]}]'''
//            country: ${country.iso_code["remote_addr"]}
//            country_name: ${country.names.en["remote_addr"]}
//            postal_code:  ${postal.code["remote_addr"]}
// ```
type _docGeoIP interface{}

// +name:"Geo IP"
// +url:"https://github.com/y-ken/fluent-plugin-geoip"
// +version:"more info"
// +description:"Fluentd GeoIP filter"
// +status:"GA"
type _metaGeoIP interface{}

// +kubebuilder:object:generate=true
type GeoIP struct {
	//Specify one or more geoip lookup field which has ip address (default: host)
	GeoipLookupKeys string `json:"geoip_lookup_keys,omitempty"`
	//Specify optional geoip database (using bundled GeoLiteCity databse by default)
	GeoipDatabase string `json:"geoip_database,omitempty"`
	//Specify optional geoip2 database (using bundled GeoLite2-City.mmdb by default)
	Geoip2Database string `json:"geoip_2_database,omitempty"`
	//Specify backend library (geoip2_c, geoip, geoip2_compat)
	BackendLibrary string `json:"backend_library,omitempty"`
	// To avoid get stacktrace error with `[null, null]` array for elasticsearch.
	SkipAddingNullRecord bool `json:"skip_adding_null_record,omitempty" plugin:"default:true"`
	// Records are represented as maps: `key: value`
	Records []Record `json:"records,omitempty"`
}

func (g *GeoIP) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "geoip"
	geoIP := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "filter",
			Tag:       "**",
			Id:        id + "_" + pluginType,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(g); err != nil {
		return nil, err
	} else {
		geoIP.Params = params
	}
	if len(g.Records) > 0 {
		for _, record := range g.Records {
			if meta, err := record.ToDirective(secretLoader, ""); err != nil {
				return nil, err
			} else {
				geoIP.SubDirectives = append(geoIP.SubDirectives, meta)
			}
		}
	}
	return geoIP, nil
}
