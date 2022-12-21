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

package mirror

import (
	"reflect"
)

// StructRange returns a range iterator for a struct value.
// It returns nil if v's Kind is not Struct.
//
// Call Next to advance the iterator, and Field/Value to access each entry.
// Next returns false when the iterator is exhausted.
// StructRange follows the same iteration semantics as a range statement.
//
// Example:
//
//	 s := struct{
//	 	...
//	 }{
//	 	...
//	 }
//		iter := unstructured.StructRange(s)
//		for iter.Next() {
//			f := iter.Field()
//			v := iter.Value()
//			...
//		}
func StructRange(v interface{}) *StructIter {
	return NewStructIter(reflect.ValueOf(v))
}

// NewStructIter returns a new instance of StructIter for the specified value.
// It returns nil if v's Kind is not Struct.
func NewStructIter(val reflect.Value) *StructIter {
	if val.Kind() != reflect.Struct {
		return nil
	}

	return &StructIter{
		val: val,
		typ: val.Type(),
		n:   val.NumField(),
		i:   -1,
	}
}

// StructIter is an iterator for ranging over a struct's fields with their values
type StructIter struct {
	val reflect.Value
	typ reflect.Type
	n   int
	i   int
}

// Next advances the iterator and reports whether there is another entry.
// It returns false when the iterator is exhausted; subsequent calls will panic.
func (it *StructIter) Next() bool {
	if it == nil {
		panic("nil iterator")
	}
	if it.i >= it.n {
		panic("StructIter.Next called on exhausted iterator")
	}
	it.i++
	return it.i < it.n
}

// Field returns the struct field of the iterator's current entry.
func (it StructIter) Field() reflect.StructField {
	if it.i < 0 {
		panic("StructIter.Field called before Next")
	}
	if it.i >= it.n {
		panic("StructIter.Field called on exhausted iterator")
	}
	return it.typ.Field(it.i)
}

// Value returns the value of the iterator's current entry.
func (it StructIter) Value() reflect.Value {
	if it.i < 0 {
		panic("StructIter.Value called before Next")
	}
	if it.i >= it.n {
		panic("StructIter.Value called on exhausted iterator")
	}
	return it.val.Field(it.i)
}
