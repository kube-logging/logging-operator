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

package types

import (
	"encoding/json"
	"strconv"
	"strings"

	util "github.com/banzaicloud/operator-tools/pkg/utils"
)

// OutputPlugin plugin: https://github.com/banzaicloud/fluent-plugin-label-router
type Router struct {
	PluginMeta
	Routes []Directive `json:"routes"`
	Params map[string]string
}

func (r *Router) GetPluginMeta() *PluginMeta {
	return &r.PluginMeta
}

func (r *Router) GetParams() map[string]string {
	return r.Params
}

func (r *Router) GetSections() []Directive {
	return r.Routes
}

type FlowMatch struct {
	// Optional set of kubernetes labels
	Labels map[string]string `json:"labels,omitempty"`
	// Optional namespace
	Namespaces []string `json:"namespaces,omitempty"`
	// ContainerNames
	ContainerNames []string `json:"container_names,omitempty"`
	// Hosts
	Hosts []string `json:"hosts,omitempty"`
	// Negate
	Negate bool `json:"negate,omitempty"`
}

func (f FlowMatch) GetPluginMeta() *PluginMeta {
	return &PluginMeta{
		Directive: "match",
	}
}
func (f FlowMatch) GetParams() map[string]string {
	params := map[string]string{
		"negate": strconv.FormatBool(f.Negate),
	}
	if len(f.Namespaces) > 0 {
		params["namespaces"] = strings.Join(f.Namespaces, ",")
	}
	if len(f.ContainerNames) > 0 {
		params["container_names"] = strings.Join(f.ContainerNames, ",")
	}
	if len(f.Hosts) > 0 {
		params["hosts"] = strings.Join(f.Hosts, ",")
	}
	if len(f.Labels) > 0 {
		var sb []string
		for _, key := range util.OrderedStringMap(f.Labels).Keys() {
			sb = append(sb, key+":"+f.Labels[key])
		}
		params["labels"] = strings.Join(sb, ",")
	}
	return params
}

func (f FlowMatch) GetSections() []Directive {
	return nil
}

type FlowRoute struct {
	PluginMeta
	Params  map[string]string
	Matches []Directive
}

func (f *FlowRoute) GetPluginMeta() *PluginMeta {
	return &f.PluginMeta
}

func (f *FlowRoute) GetParams() map[string]string {
	return f.Params
}

func (f *FlowRoute) GetSections() []Directive {
	return f.Matches
}

func (r *Router) AddRoute(flow *Flow) *Router {
	metrics_labels, _ := json.Marshal(map[string]string{"id": flow.FlowID})
	route := &FlowRoute{
		PluginMeta: PluginMeta{
			Directive: "route",
			Label:     flow.FlowLabel,
		},
		Params: map[string]string{
			"metrics_labels": string(metrics_labels),
		},
	}
	for _, f := range flow.Matches {
		route.Matches = append(route.Matches, f)
	}
	r.Routes = append(r.Routes, route)
	return r
}

func NewRouter(id string) *Router {
	pluginType := "label_router"
	return &Router{
		PluginMeta: PluginMeta{
			Type:      "label_router",
			Directive: "match",
			Tag:       "**",
			Id:        id + "_" + pluginType,
		},
		Params: map[string]string{
			"metrics": "true",
		},
	}
}
