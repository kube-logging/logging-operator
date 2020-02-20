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
	"errors"
)

type Builder struct {
	input  Input
	flows  []*Flow
	router *Router
}

func NewSystem(input Input, router *Router) *Builder {
	return &Builder{
		input:  input,
		router: router,
	}
}

func (s *Builder) RegisterFlow(f *Flow) error {
	for _, e := range s.flows {
		if e.FlowLabel == f.FlowLabel {
			return errors.New("Flow already exists")
		}
	}
	s.flows = append(s.flows, f)
	s.router.AddRoute(f)
	return nil
}

func (s *Builder) Build() (*System, error) {
	return &System{
		Input:  s.input,
		Router: s.router,
		Flows:  s.flows,
	}, nil
}
