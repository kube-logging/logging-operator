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

package syslogng

import "io"

type ConfigRenderer interface {
	RenderAsSyslogNGConfig(ctx Context) error
}

type Context struct {
	Out io.Writer

	Depth  int
	Indent string

	ControlNamespace string
}

func (ctx Context) WithControlNamespace(ns string) Context {
	ctx.ControlNamespace = ns
	return ctx
}

func (ctx Context) WithDepth(depth int) Context {
	ctx.Depth = depth
	return ctx
}
