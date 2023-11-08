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

package plugins

import (
	"reflect"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

type DirectiveConverter interface {
	ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error)
}

func CreateOutput(outputSpec v1beta1.OutputSpec, outputName string, secretLoader secret.SecretLoader) (types.Directive, error) {
	v := reflect.ValueOf(outputSpec)
	var converters []DirectiveConverter
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
			if converter, ok := v.Field(i).Interface().(DirectiveConverter); ok {
				converters = append(converters, converter)
			}
		}
	}
	switch len(converters) {
	case 0:
		return nil, errors.New("no plugin config available for output")
	case 1:
		return converters[0].ToDirective(secretLoader, outputName)
	default:
		return nil, errors.Errorf("more then one plugin config is not allowed for an output")
	}
}

func CreateFilter(filter v1beta1.Filter, id string, secretLoader secret.SecretLoader) (types.Directive, error) {
	v := reflect.ValueOf(filter)
	var converters []DirectiveConverter
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
			if converter, ok := v.Field(i).Interface().(DirectiveConverter); ok {
				converters = append(converters, converter)
			}
		}
	}
	switch len(converters) {
	case 0:
		return nil, errors.New("no plugin config available for filter")
	case 1:
		return converters[0].ToDirective(secretLoader, id)
	default:
		return nil, errors.Errorf("more then one plugin config is not allowed for a filter")
	}
}
