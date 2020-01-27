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
	"crypto/md5"
	"fmt"
	"io"
	"sort"
)

type FluentConfig interface {
	GetDirectives() []Directive
}

type System struct {
	Input  Input   `json:"input"`
	Router *Router `json:"router"`
	Flows  []*Flow `json:"flows"`
}

func (s *System) GetDirectives() []Directive {
	directives := []Directive{
		s.Input,
		s.Router,
	}
	for _, flow := range s.Flows {
		directives = append(directives, flow)
	}
	return directives
}

type Flow struct {
	PluginMeta

	// Chain of Filters that will process the event. Can be zero or more.
	Filters []Filter `json:"filters,omitempty"`
	// List of Outputs that will emit the event, at least one output is required.
	Outputs []Output `json:"outputs"`

	// Optional set of kubernetes labels
	Labels map[string]string `json:"-"`
	// Optional namespace
	Namespace string `json:"-"`

	// Fluentd label
	FlowLabel string `json:"-"`
}

func (f *Flow) GetPluginMeta() *PluginMeta {
	return &f.PluginMeta
}

func (f *Flow) GetParams() map[string]string {
	return nil
}

func (f *Flow) GetSections() []Directive {
	sections := []Directive{}
	for _, filter := range f.Filters {
		sections = append(sections, filter)
	}
	if len(f.Outputs) > 1 {
		// We have to convert to General directive
		sections = append(sections, NewCopyDirective(f.Outputs))
	} else {
		for _, output := range f.Outputs {
			sections = append(sections, output)
		}
	}

	return sections
}

func (f *Flow) WithFilters(filter ...Filter) *Flow {
	f.Filters = append(f.Filters, filter...)
	return f
}

func (f *Flow) WithOutputs(output ...Output) *Flow {
	f.Outputs = append(f.Outputs, output...)
	return f
}

func NewFlow(namespaces []string, labels map[string]string) (*Flow, error) {
	flowLabel, err := calculateFlowLabel(namespaces, labels)
	if err != nil {
		return nil, err
	}
	return &Flow{
		PluginMeta: PluginMeta{
			Directive: "label",
			Tag:       flowLabel,
		},
		FlowLabel: flowLabel,
		Labels:    labels,
		Namespace: namespace,
	}, nil
}

func calculateFlowLabel(namespaces []string, labels map[string]string) (string, error) {
	b := md5.New()
	sort.Strings(namespaces)
	for _, n := range namespaces {
		if _, err := io.WriteString(b, n); err != nil {
			return "", err
		}
	}
	// Make sure the generated label is consistent
	keys := []string{}
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if _, err := io.WriteString(b, k); err != nil {
			return "", err
		}
		if _, err := io.WriteString(b, labels[k]); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("@%x", b.Sum(nil)), nil
}
