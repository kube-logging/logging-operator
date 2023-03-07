// Copyright © 2022 Cisco Systems, Inc. and/or its affiliates
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

	"github.com/siliconbrain/go-seqs/seqs"
)

type Field struct {
	Value reflect.Value
	Meta  reflect.StructField
}

func (f Field) KeyOrEmpty() string {
	key, _ := fieldKey(f, nil)
	return key
}

func fieldsOf(val reflect.Value) []Field {
	val = derefAll(val)
	if val.Kind() != reflect.Struct {
		return nil
	}

	exportedStructFields := seqs.SeqFunc(func(fn func(reflect.StructField) bool) {
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			if f := typ.Field(i); f.IsExported() {
				if fn(f) {
					return
				}
			}
		}
	})
	return seqs.ToSlice(seqs.Filter(
		seqs.Map(
			seqs.Filter(exportedStructFields, func(f reflect.StructField) bool {
				return !structFieldSettings(f).Ignore()
			}),
			func(meta reflect.StructField) Field {
				return Field{
					Meta:  meta,
					Value: val.FieldByIndex(meta.Index),
				}
			},
		),
		func(f Field) bool {
			return !(structFieldSettings(f.Meta).Optional() || hasJSONOmitempty(f.Meta)) || !f.Value.IsZero()
		},
	))
}

func fieldKey(f Field, settings syslogNGTagSettings) (key string, err error) {
	if settings == nil {
		settings = structFieldSettings(f.Meta)
	}
	key = settings.Name()
	if key == "" {
		dynamicMetaFieldSettings := structFieldSettings(metaField(f.Value.Type()))
		key = dynamicMetaFieldSettings.Name()
	}
	if key == "" {
		staticMetaFieldSettings := structFieldSettings(metaField(f.Meta.Type))
		key = staticMetaFieldSettings.Name()
	}
	if key == "" {
		key = jsonNameOf(f.Meta)
	}
	if key == "" {
		key = goNameToSyslogName(f.Meta.Name)
	}
	if key == "" {
		key = goNameToSyslogName(f.Meta.Type.Name())
	}
	if key == "" {
		err = fmt.Errorf("option key could not be determined for field %#v", f)
	}
	return
}

func metaField(t reflect.Type) (f reflect.StructField) {
	t = derefAll(t)
	if t.Kind() == reflect.Struct {
		f, _ = t.FieldByName("__meta")
	}
	return
}
