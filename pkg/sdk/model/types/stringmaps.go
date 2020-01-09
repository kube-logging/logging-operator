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
	"fmt"
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
)

type Converter func(interface{}) (string, error)

type StructToStringMapper struct {
	TagName         string
	PluginTagName   string
	ConversionHooks map[string]Converter
	SecretLoader    secret.SecretLoader
}

type NullSecretLoader struct {
}

func NewStructToStringMapper(secretLoader secret.SecretLoader) *StructToStringMapper {
	return &StructToStringMapper{
		TagName:         "json",
		PluginTagName:   "plugin",
		ConversionHooks: make(map[string]Converter),
		SecretLoader:    secretLoader,
	}
}

func (s *StructToStringMapper) WithConverter(name string, c Converter) *StructToStringMapper {
	s.ConversionHooks[name] = c
	return s
}

func (s *StructToStringMapper) StringsMap(in interface{}) (map[string]string, error) {
	out := make(map[string]string)
	err := s.fillMap(strctVal(in), out)
	return out, err
}

func (s *StructToStringMapper) fillMap(value reflect.Value, out map[string]string) error {
	if out == nil {
		return nil
	}

	fields := s.structFields(value)

	var multierror error
	for _, field := range fields {
		name := field.Name
		val := value.FieldByName(name)
		var finalVal string

		tagName, tagOpts := parseTagWithName(field.Tag.Get(s.TagName))
		if tagName != "" {
			name = tagName
		}

		pluginTagOpts := parseTag(field.Tag.Get(s.PluginTagName))
		required := pluginTagOpts.Has("required")

		if tagOpts.Has("omitempty") {
			if pluginTagOpts.Has("required") {
				multierror = errors.Combine(multierror, errors.Errorf(
					"tags for field %s are conflicting: required and omitempty cannot be set simultaneously", name))
				continue
			}
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()
			if reflect.DeepEqual(current, zero) {
				if ok, def := pluginTagOpts.ValueForPrefix("default:"); ok {
					out[name] = def
				}
				continue
			}
		}

		var v reflect.Value
		if ok, converterName := pluginTagOpts.ValueForPrefix("converter:"); ok {
			if hook, ok := s.ConversionHooks[converterName]; ok {
				convertedValue, err := hook(val.Interface())
				if err != nil {
					multierror = errors.Combine(err, errors.Errorf(
						"failed to convert field `%s` with converter %s", name, converterName))
				} else {
					v = reflect.ValueOf(convertedValue)
				}
			} else {
				multierror = errors.Combine(multierror, errors.Errorf(
					"unable to convert field `%s` as the specified converter `%s` is not registered", name, converterName))
				continue
			}
		} else {
			v = reflect.ValueOf(val.Interface())
		}

		if s.SecretLoader != nil {
			if v.Kind() == reflect.Ptr && !v.IsNil() {
				if secretItem, ok := val.Interface().(*secret.Secret); ok {
					loadedSecret, err := s.SecretLoader.Load(secretItem)
					if err != nil {
						multierror = errors.Combine(multierror, errors.WrapIff(err, "failed to load secret for field %s", name))
					}
					if err == nil {
						out[name] = loadedSecret
					}
					continue
				}
			}
		}

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// if the field is of string type and not empty, use it's value over the default
		switch v.Kind() {
		case reflect.String, reflect.Int, reflect.Bool:
			stringVal := fmt.Sprintf("%v", v)
			if stringVal != "" {
				finalVal = stringVal
			} else {
				// check if default has been set and use it
				if ok, def := pluginTagOpts.ValueForPrefix("default:"); ok {
					finalVal = def
				}
			}
			// can't let to return an empty string when it's required
			if finalVal == "" && required {
				multierror = errors.Combine(multierror, errors.Errorf("field %s is required", name))
			} else {
				out[name] = finalVal
			}
		case reflect.Slice:
			if stringSlice, ok := v.Interface().([]string); ok {
				if len(stringSlice) > 0 {
					b, err := json.Marshal(stringSlice)
					if err != nil {
						multierror = errors.Combine(multierror, errors.Errorf("can't marshal field: %q value: %q as json", name, stringSlice), err)
					}
					out[name] = string(b)
				} else {
					if ok, def := pluginTagOpts.ValueForPrefix("default:"); ok {
						b, err := json.Marshal(strings.Split(def, ","))
						if err != nil {
							multierror = errors.Combine(multierror, errors.Errorf("can't marshal field: %q value: %q as json", name, def), err)
						}
						out[name] = string(b)
					}
				}
			}
		case reflect.Map:
			if mapStringString, ok := v.Interface().(map[string]string); ok {
				if len(mapStringString) > 0 {
					b, err := json.Marshal(mapStringString)
					if err != nil {
						multierror = errors.Combine(multierror, errors.Errorf("can't marshal field: %q value: %q as json", name, mapStringString), err)
					}
					out[name] = string(b)
				} else {
					if ok, def := pluginTagOpts.ValueForPrefix("default:"); ok {
						validate := map[string]string{}
						if err := json.Unmarshal([]byte(def), &validate); err != nil {
							multierror = errors.Combine(multierror, errors.Errorf("can't marshal field: %q value: %q as json", name, def), err)
						}
						out[name] = def
					}
				}
			}
		}
	}
	return multierror
}

func strctVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	return v
}

func (s *StructToStringMapper) structFields(value reflect.Value) []reflect.StructField {
	t := value.Type()

	var f []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get(s.TagName); tag == "-" {
			continue
		}

		f = append(f, field)
	}

	return f
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func parseTagWithName(tag string) (string, tagOptions) {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"

	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

// tagOptions contains a slice of tag options
type tagOptions []string

// Has returns true if the given option is available in tagOptions
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

// Has returns true if the given option is available in tagOptions
func (t tagOptions) ValueForPrefix(opt string) (bool, string) {
	for _, tagOpt := range t {
		if strings.HasPrefix(tagOpt, opt) {
			return true, strings.Replace(tagOpt, opt, "", 1)
		}
	}
	return false, ""
}

// parseTag returns all the options in the tag
func parseTag(tag string) tagOptions {
	return tagOptions(strings.Split(tag, ","))
}
