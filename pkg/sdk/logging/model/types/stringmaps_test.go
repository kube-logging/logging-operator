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

package types_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
	corev1 "k8s.io/api/core/v1"
)

func TestRequired(t *testing.T) {
	expectedError := "field \"field1\" is required"
	type Asd struct {
		Field1 string `json:"field1" plugin:"required"`
	}
	_, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
	if err == nil {
		t.Fatalf("required error is expected")
	} else {
		if err.Error() != expectedError {
			t.Fatalf("error message `%s` does not match expected `%s`", err.Error(), expectedError)
		}
	}
}

func TestRequiredMeansItCannotEvenBeEmpty(t *testing.T) {
	expectedError := "field \"field1\" is required"
	type Asd struct {
		Field1 string `json:"field1" plugin:"required"`
	}
	_, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{Field1: ""})
	if err == nil {
		t.Fatalf("required error is expected")
	} else {
		if err.Error() != expectedError {
			t.Fatalf("error message `%s` does not match expected `%s`", err.Error(), expectedError)
		}
	}
}

func TestJsonTagsWithDefaultsAndOmitempty(t *testing.T) {
	type Asd struct {
		Field1 string  `json:"field1"`
		Field2 string  `json:"field2,omitempty" plugin:"default:http://asdf and some space"`
		Field3 string  `json:"field3,omitempty"`
		Field4 *string `json:"field4,omitempty" plugin:"default:nonempty"`
		Field5 *string `json:"field5,omitempty" plugin:"default:nonempty"`
	}
	actual, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(
		Asd{
			Field1: "value",
			Field4: utils.StringPointer(""), // looks like empty, but it's not, we expect this to be set
			Field5: nil,                     // empty for real
		})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{
		"field1": "value",
		"field2": "http://asdf and some space",
		"field4": "",
		"field5": "nonempty",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match\nexpected:\n%+v\n\nactual:\n%+v", expected, actual)
	}
}

func TestSliceFields(t *testing.T) {
	type Asd struct {
		Field1 []string `json:"field1,omitempty"`
		Field2 []string `json:"field2" plugin:"default:item1"`
		Field3 []int    `json:"field3,omitempty"`
		Field4 []int    `json:"field4" plugin:"default:1"`
	}

	tests := []struct {
		source   Asd
		expected map[string]string
	}{
		{
			source: Asd{},
			expected: map[string]string{
				"field2": `["item1"]`,
				"field4": "[1]",
			},
		},
		{
			source: Asd{
				Field1: []string{"item1", "item2"},
				Field2: []string{"item3"},
				Field3: []int{1, 2},
				Field4: []int{3},
			},
			expected: map[string]string{
				"field1": `["item1","item2"]`,
				"field2": `["item3"]`,
				"field3": "[1,2]",
				"field4": "[3]",
			},
		},
	}

	for _, tt := range tests {
		actual, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(tt.source)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if !reflect.DeepEqual(tt.expected, actual) {
			t.Fatalf("failed to match expected %+v with %+v", tt.expected, actual)
		}
	}
}

func TestInvalidSliceDefault(t *testing.T) {
	expectedError := `can't unmarshal default value "[str]" into field "field1"`
	type Asd struct {
		Field1 []int `json:"field1" plugin:"default:str"`
	}
	_, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
	if err == nil {
		t.Fatalf("required error is expected")
	} else {
		if !strings.HasPrefix(err.Error(), expectedError) {
			t.Fatalf("error message `%s` does not have expected prefix `%s`", err.Error(), expectedError)
		}
	}
}

func TestConflictingTags(t *testing.T) {
	expectedError := "tags for field \"field2\" are conflicting: required and omitempty cannot be set simultaneously"
	type Asd struct {
		Field2 string `json:"field2,omitempty" plugin:"required"`
	}
	_, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
	if err == nil {
		t.Fatalf("required error is expected")
	} else {
		if err.Error() != expectedError {
			t.Fatalf("error message `%s` does not match expected `%s`", err.Error(), expectedError)
		}
	}
}

func TestIgnoreNestedStructs(t *testing.T) {
	type Nested struct {
		Field string `json:"asd"`
	}
	type Asd struct {
		Field2 string  `json:"field2"`
		Field3 *Nested `json:"nested"`
	}
	actual, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{Field2: "val"})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{
		"field2": "val",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match expected %+v with %+v", expected, actual)
	}
}

func TestEmptyStructStructs(t *testing.T) {
	type Asd struct {
	}
	actual, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match expected %+v with %+v", expected, actual)
	}
}

func TestConversion(t *testing.T) {
	type Asd struct {
		Field int `json:"field" plugin:"converter:magic"`
	}

	converter := func(f interface{}) (string, error) {
		if converted, ok := f.(int); ok {
			return strconv.Itoa(converted), nil
		}
		return "", errors.Errorf("unable to convert %+v to int", f)
	}

	testStruct := Asd{Field: 2}

	actual, err := types.NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).
		WithConverter("magic", converter).
		StringsMap(testStruct)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{
		"field": "2",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match expected %+v with %+v", expected, actual)
	}
}

type FakeLoader struct {
}

func (d *FakeLoader) Load(secret *secret.Secret) (string, error) {
	if secret.Value != "" {
		return secret.Value, nil
	}
	if secret.ValueFrom != nil && secret.ValueFrom.SecretKeyRef != nil {
		return fmt.Sprintf("%s:%s",
			secret.ValueFrom.SecretKeyRef.Name,
			secret.ValueFrom.SecretKeyRef.Key), nil
	}
	return "", errors.New("no value found")
}

func (d *FakeLoader) Mount(secret *secret.Secret) (string, error) {
	if secret.Value != "" {
		return secret.Value, nil
	}
	if secret.ValueFrom != nil && secret.ValueFrom.SecretKeyRef != nil {
		return fmt.Sprintf("mountedFrom:%s:%s",
			secret.ValueFrom.SecretKeyRef.Name,
			secret.ValueFrom.SecretKeyRef.Key), nil
	}
	return "", errors.New("no value found")
}

func TestSecretValue(t *testing.T) {
	type Asd struct {
		Field *secret.Secret `json:"field"`
	}

	testStruct := Asd{Field: &secret.Secret{Value: "asd"}}

	actual, err := types.NewStructToStringMapper(&FakeLoader{}).
		StringsMap(testStruct)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{
		"field": "asd",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match expected %+v with %+v", expected, actual)
	}
}

func TestSecretValueFrom(t *testing.T) {
	type Asd struct {
		Field *secret.Secret `json:"field"`
	}

	testStruct := Asd{
		Field: &secret.Secret{
			ValueFrom: &secret.ValueFrom{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "a",
					},
					Key: "b",
				},
			},
		},
	}

	actual, err := types.NewStructToStringMapper(&FakeLoader{}).
		StringsMap(testStruct)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{
		"field": "a:b",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match expected %+v with %+v", expected, actual)
	}
}

func TestSecretErrorWhenEmpty(t *testing.T) {
	type Asd struct {
		Field *secret.Secret `json:"field"`
	}

	testStruct := Asd{
		Field: &secret.Secret{},
	}

	_, err := types.NewStructToStringMapper(&FakeLoader{}).
		StringsMap(testStruct)
	if err == nil {
		t.Fatal("expected an error when secret contains no value or valuefrom")
	}

	expectedError := "failed to load secret for field \"field\": no value found"
	if err.Error() != expectedError {
		t.Fatalf("Expected `%s` got `%s`", expectedError, err.Error())
	}
}
