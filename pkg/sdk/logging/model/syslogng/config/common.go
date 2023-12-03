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

package config

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/siliconbrain/go-seqs/seqs"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
)

func renderAny(value any, secretLoader secret.SecretLoader) []render.Renderer {
	return renderValue(reflect.ValueOf(value), secretLoader)
}

func renderValue(value reflect.Value, secretLoader secret.SecretLoader) []render.Renderer {
	if !value.IsValid() {
		return nil
	}

	// handle concrete types first
	switch val := value.Interface().(type) {
	case secret.Secret:
		sec, err := secretLoader.Load(&val)
		if err != nil {
			return []render.Renderer{render.Error(err)}
		}
		return []render.Renderer{render.Literal(sec)}
	}

	if value.CanConvert(matchExprType) {
		matchExpr := value.Convert(matchExprType).Interface().(filter.MatchExpr)
		return []render.Renderer{
			filterExpr(filterExprFromMatchExpr(matchExpr)),
		}
	} else if value.Type() == arrowMapType {
		arrowMap := value.Interface().(filter.ArrowMap)
		return []render.Renderer{render.ArrowMap(arrowMap, render.Literal[string], render.Literal[string])}
	} else if value.Type() == rawArrowMapType {
		rawArrowMap := value.Interface().(filter.RawArrowMap)
		return []render.Renderer{render.ArrowMap(rawArrowMap, render.String, render.String)}
	} else if value.Type() == stringListType {
		stringList := value.Interface().(output.StringList)
		return []render.Renderer{render.StringList(stringList.List)}
	} else if value.Type() == rawStringType {
		rawString := value.Interface().(output.RawString)
		return []render.Renderer{render.String(string(rawString))}
	}

	switch value.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.Pointer:
		return renderValue(derefAll(value), secretLoader)
	case reflect.Bool:
		return []render.Renderer{render.Literal(value.Bool())}
	case reflect.String:
		return []render.Renderer{render.Literal(value.String())}
	case reflect.Float32:
		return []render.Renderer{render.Literal(float32(value.Float()))}
	case reflect.Float64:
		return []render.Renderer{render.Literal(value.Float())}
	case reflect.Int:
		return []render.Renderer{render.Literal(int(value.Int()))}
	case reflect.Int16:
		return []render.Renderer{render.Literal(int16(value.Int()))}
	case reflect.Int32:
		return []render.Renderer{render.Literal(int32(value.Int()))}
	case reflect.Int64:
		return []render.Renderer{render.Literal(value.Int())}
	case reflect.Int8:
		return []render.Renderer{render.Literal(int8(value.Int()))}
	case reflect.Uint:
		return []render.Renderer{render.Literal(uint(value.Uint()))}
	case reflect.Uint16:
		return []render.Renderer{render.Literal(uint16(value.Uint()))}
	case reflect.Uint32:
		return []render.Renderer{render.Literal(uint32(value.Uint()))}
	case reflect.Uint64:
		return []render.Renderer{render.Literal(value.Uint())}
	case reflect.Uint8:
		return []render.Renderer{render.Literal(uint8(value.Uint()))}
	case reflect.Array, reflect.Slice:
		if value.Len() == 0 {
			return nil
		}
		res := make([]render.Renderer, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			res = append(res, renderValue(value.Index(i), secretLoader)...)
		}
		return res
	case reflect.Map:
		var res []render.Renderer
		if l := value.Len(); l > 0 {
			res = make([]render.Renderer, 0, value.Len())
		}
		for _, keyVal := range value.MapKeys() {
			switch keyVal.Kind() {
			case reflect.String:
				res = append(res, optionExpr(keyVal.String(), renderValue(value.MapIndex(keyVal), secretLoader)...))
			default:
				res = append(res, render.Error(fmt.Errorf("cannot render map entry with key type %s", keyVal.Type())))
			}
		}
		return res
	case reflect.Struct:
		fs := fieldsOf(value)

		var posArgs []posArg
		var nonPos []render.Renderer
		for _, f := range fs {
			renderField(f, secretLoader, &nonPos, &posArgs)
		}
		sort.Slice(posArgs, func(a, b int) bool { return posArgs[a].pos < posArgs[b].pos })
		var res []render.Renderer
		for _, a := range posArgs {
			res = append(res, a.rnds...)
		}
		res = append(res, nonPos...)
		return res
	}
	return []render.Renderer{render.Error(fmt.Errorf("cannot render value of type %s", value.Type()))}
}

type posArg struct {
	pos  int
	rnds []render.Renderer
}

func renderField(f Field, secretLoader secret.SecretLoader, nonPos *[]render.Renderer, posArgs *[]posArg) {
	if f.Meta.Anonymous {
		for _, ff := range fieldsOf(f.Value) {
			renderField(ff, secretLoader, nonPos, posArgs)
		}
		return
	}
	settings := structFieldSettings(f.Meta)
	if posStr := settings["pos"]; posStr != "" {
		// this field is a positional argument
		pos64, err := strconv.ParseInt(posStr, 10, 8)
		if err != nil {
			*nonPos = append(*nonPos, render.Error(fmt.Errorf("invalid position specifier %q in tag: %w", posStr, err)))
			return
		}
		*posArgs = append(*posArgs, posArg{
			pos:  int(pos64),
			rnds: renderValue(f.Value, secretLoader),
		})
		return
	}
	// this field is an option (keyword argument)
	key, err := fieldKey(f, settings)
	if err != nil {
		*nonPos = append(*nonPos, render.Error(err))
		return
	}
	*nonPos = append(*nonPos, optionExpr(key, renderValue(f.Value, secretLoader)...))
}

func renderDriver(f Field, secretLoader secret.SecretLoader) render.Renderer {
	name, err := fieldKey(f, nil)
	if err != nil {
		return render.Error(err)
	}
	switch f.Value.Kind() {
	case reflect.Array, reflect.Slice:
		var stmts []render.Renderer
		l := f.Value.Len()
		if l > 0 {
			stmts = make([]render.Renderer, l)
		}
		for i := 0; i < l; i++ {
			stmts[i] = parenDefStmt(name, render.SpaceSeparated(renderValue(f.Value.Index(i), secretLoader)...))
		}
		return render.AllFrom(seqs.Map(seqs.FromSlice(stmts), render.Line))
	default:
		return parenDefStmt(name, render.SpaceSeparated(renderValue(f.Value, secretLoader)...))
	}
}

var capitalSubstrings = regexp.MustCompile("[A-Z]+")

func goNameToSyslogName(s string) string {
	return strings.TrimRight(capitalSubstrings.ReplaceAllStringFunc(s, func(s string) string {
		return strings.ToLower(s) + "-"
	}), "-")
}

var matchExprType = reflect.TypeOf(filter.MatchExpr{})
var arrowMapType = reflect.TypeOf(filter.ArrowMap{})
var rawArrowMapType = reflect.TypeOf(filter.RawArrowMap{})
var stringListType = reflect.TypeOf(output.StringList{})
var rawStringType = reflect.TypeOf(*new(output.RawString))

func derefAll[T Derefable[T]](v T) T {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}

type Derefable[T any] interface {
	Kind() reflect.Kind
	Elem() T
}

func isActiveField(f Field) bool {
	switch f.Meta.Type.Kind() {
	case reflect.Array, reflect.Interface, reflect.Pointer, reflect.Slice:
		return !f.Value.IsZero()
	default:
		return true
	}
}
