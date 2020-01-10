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
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

func TestRequired(t *testing.T) {
	expectedError := "field field1 is required"
	type Asd struct {
		Field1 string `json:"field1" plugin:"required"`
	}
	_, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
	if err == nil {
		t.Fatalf("required error is expected")
	} else {
		if err.Error() != expectedError {
			t.Fatalf("error message `%s` does not match expected `%s`", err.Error(), expectedError)
		}
	}
}

func TestRequiredMeansItCannotEvenBeEmpty(t *testing.T) {
	expectedError := "field field1 is required"
	type Asd struct {
		Field1 string `json:"field1" plugin:"required"`
	}
	_, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{Field1: ""})
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
		Field1 string `json:"field1"`
		Field2 string `json:"field2,omitempty" plugin:"default:http://asdf and some space"`
		Field3 string `json:"field3,omitempty"`
	}
	actual, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{Field1: "value"})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expected := map[string]string{
		"field1": "value",
		"field2": "http://asdf and some space",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("failed to match expected %+v with %+v", expected, actual)
	}
}

func TestConflictingTags(t *testing.T) {
	expectedError := "tags for field field2 are conflicting: required and omitempty cannot be set simultaneously"
	type Asd struct {
		Field2 string `json:"field2,omitempty" plugin:"required"`
	}
	_, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
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
	actual, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{Field2: "val"})
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
	actual, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).StringsMap(Asd{})
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

	actual, err := NewStructToStringMapper(secret.NewSecretLoader(nil, "", "", nil)).
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

	actual, err := NewStructToStringMapper(&FakeLoader{}).
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
				SecretKeyRef: &secret.KubernetesSecret{
					Name: "a",
					Key:  "b",
				},
			},
		},
	}

	actual, err := NewStructToStringMapper(&FakeLoader{}).
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

	_, err := NewStructToStringMapper(&FakeLoader{}).
		StringsMap(testStruct)
	if err == nil {
		t.Fatal("expected an error when secret contains no value or valuefrom")
	}

	expectedError := "failed to load secret for field field; no value found"
	if err.Error() != expectedError {
		t.Fatalf("Expected `%s` got `%s`", expectedError, err.Error())
	}
}
