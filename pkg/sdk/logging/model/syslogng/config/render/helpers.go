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

package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/siliconbrain/go-seqs/seqs"
	"golang.org/x/exp/constraints"
)

type Renderer func(ctx RenderContext) error

type RenderContext struct {
	Out         io.Writer
	IndentWith  string
	IndentDepth int
}

func Error(err error) Renderer {
	return func(RenderContext) error {
		return err
	}
}

func String(s string) Renderer {
	return func(ctx RenderContext) error {
		return writeString(ctx.Out, s)
	}
}

func Formatted(format string, args ...any) Renderer {
	return func(ctx RenderContext) error {
		_, err := fmt.Fprintf(ctx.Out, format, args...)
		return err
	}
}

func Quoted(s string) Renderer {
	return Formatted("%q", s)
}

func AllOf(rs ...Renderer) Renderer {
	return func(ctx RenderContext) error {
		for _, r := range rs {
			if r == nil {
				continue
			}
			if err := r(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func AllFrom(rs seqs.Seq[Renderer]) Renderer {
	return func(ctx RenderContext) (err error) {
		rs.ForEachUntil(func(r Renderer) bool {
			if r != nil {
				err = r(ctx)
			}
			return err != nil
		})
		return
	}
}

func SpaceSeparated(rs ...Renderer) Renderer {
	return AllFrom(seqs.Intersperse(seqs.Filter(seqs.FromSlice(rs), func(r Renderer) bool { return r != nil }), Space))
}

func Line(r Renderer) Renderer {
	return AllOf(Indentation, r, NewLine)
}

func Indentation(ctx RenderContext) error {
	return writeString(ctx.Out, strings.Repeat(ctx.IndentWith, ctx.IndentDepth))
}

func NewLine(ctx RenderContext) error {
	return writeString(ctx.Out, "\n")
}

func Space(ctx RenderContext) error {
	return writeString(ctx.Out, " ")
}

func Empty(RenderContext) error {
	return nil
}

func If(cond bool, r Renderer) Renderer {
	if !cond {
		return nil
	}
	return r
}

func Indented(r Renderer) Renderer {
	return func(ctx RenderContext) error {
		ctx.IndentDepth += 1
		return r(ctx)
	}
}

func Literal[T LiteralTypes](v T) Renderer {
	switch v := any(v).(type) {
	case string:
		return Quoted(v)
	case bool:
		if v {
			return String("yes")
		}
		return String("no")
	default:
		return Formatted("%v", v)
	}
}

type LiteralTypes interface {
	bool | string | constraints.Float | constraints.Integer
}

func writeString(w io.Writer, s string) error {
	_, err := io.WriteString(w, s)
	return err
}
