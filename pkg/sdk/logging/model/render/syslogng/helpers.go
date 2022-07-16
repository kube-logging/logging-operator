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

import (
	"fmt"
	"io"
	"strings"
)

type ConfigRendererFunc func(ctx Context) error

func (fn ConfigRendererFunc) RenderAsSyslogNGConfig(ctx Context) error {
	return fn(ctx)
}

func AllOf(components ...ConfigRenderer) ConfigRenderer {
	return ConfigRendererFunc(func(ctx Context) error {
		for _, c := range components {
			if c == nil {
				continue
			}
			if err := c.RenderAsSyslogNGConfig(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

type String string

func (s String) RenderAsSyslogNGConfig(ctx Context) error {
	_, err := io.WriteString(ctx.Out, string(s))
	return err
}

func Printf(format string, args ...interface{}) ConfigRenderer {
	return ConfigRendererFunc(func(ctx Context) error {
		_, err := fmt.Fprintf(ctx.Out, format, args...)
		return err
	})
}

func Indent(r ConfigRenderer) ConfigRenderer {
	return ConfigRendererFunc(func(ctx Context) error {
		return r.RenderAsSyslogNGConfig(ctx.WithDepth(ctx.Depth + 1))
	})
}

func RenderIf(cond bool, r ConfigRenderer) ConfigRenderer {
	return ConfigRendererFunc(func(ctx Context) error {
		if cond {
			return r.RenderAsSyslogNGConfig(ctx)
		}
		return nil
	})
}

func Indentation() ConfigRenderer {
	return ConfigRendererFunc(func(ctx Context) error {
		_, err := io.WriteString(ctx.Out, strings.Repeat(ctx.Indent, ctx.Depth))
		return err
	})
}

func NoIndent(r ConfigRenderer) ConfigRenderer {
	return ConfigRendererFunc(func(ctx Context) error {
		ctx.Indent = ""
		return r.RenderAsSyslogNGConfig(ctx)
	})
}
