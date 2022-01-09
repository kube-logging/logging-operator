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

package volumepath

import (
	"regexp"
	"sort"
	"strings"
)

// List is an inherited receiver of []string
type List []string

// ApplyFn function signature to apply function to every elements
type ApplyFn func([]string, int) *string

// New is a constructor for List
func New() *List {
	return &List{}
}

// Init is an initializer function for List
func Init(strings []string) *List {
	if strings == nil {
		return nil
	}
	sl := List(strings)
	return &sl
}

// Reference makes reference to List
func Reference(List List) *List {
	return &List
}

// StringReference make reference to a String
func StringReference(str string) *string {
	return &str
}

// Apply applies the given ApplyFn function on every element of the list
func (l *List) Apply(fn ApplyFn) *List {
	if fn == nil || l == nil {
		return l
	}
	result := List{}
	for index := range *l {
		fnRes := fn(*l, index)
		if fnRes != nil {
			result = append(result, *fnRes)
		}
	}
	return &result
}

// Uniq is a special filter function which produces uniq elements
func (l *List) Uniq() *List {
	if l == nil {
		return nil
	}

	uniqIndex := make(map[string]bool, len(*l))

	filterFn := func(strings []string, index int) *string {
		str := strings[index]
		if uniqIndex[str] {
			return nil
		}
		uniqIndex[str] = true
		return &str
	}

	result := l.Apply(filterFn)

	return result
}

// First returns the first item of the List
func (l *List) First() *string {
	if l != nil && len(*l) > 0 {
		return &(*l)[0]
	}
	return nil
}

// Last returns the last item of the List
func (l *List) Last() *string {
	if l == nil || len(*l) == 0 {
		return nil
	}
	return &(*l)[len(*l)-1]
}

// Strings is a type casting helper function
func (l *List) Strings() []string {
	return []string(*l)
}

// TopLevelPathList returns the top level path which is represented in list for every item
func (l *List) TopLevelPathList() *List {
	modify := ApplyFn(
		func(paths []string, index int) *string {
			path := paths[index]

			// filter all path by matching subpath
			matches := l.Apply(ApplyFn(
				func(strs []string, idx int) *string {
					str := strs[idx]
					if strings.HasPrefix(path, str) {
						return &str
					}
					return nil
				}))

			// sort results by length
			sort.Slice(*matches, func(a, b int) bool {
				return len((*matches)[a]) < len((*matches)[b])
			})

			// take the shortest
			return matches.First()
		})

	return l.Apply(modify)
}

// RemoveInvalidPath filters out file pathes do not match validatorFn
func (l *List) RemoveInvalidPath(validatorFn ApplyFn) *List {
	if l == nil {
		return nil
	}
	if validatorFn == nil {
		re := regexp.MustCompile("/.+")
		validatorFn = ApplyFn(
			func(strs []string, idx int) *string {
				if !re.MatchString(strs[idx]) {
					return nil
				}
				return &strs[idx]
			},
		)
	}
	result := l.Apply(validatorFn)
	return result
}

// ConvertFilePath .
func ConvertFilePath(path string) string {
	separator := "-"
	re := regexp.MustCompile("/([^/]+)")
	matches := re.FindAllStringSubmatch(path, -1)

	dirs := []string{}
	for _, dir := range matches {
		dirs = append(dirs, dir[1])
	}

	return EscapeDNS1123(strings.Join(dirs, separator))
}

// EscapeDNS1123 .
func EscapeDNS1123(name string) string {
	name = strings.ToLower(name)
	re := regexp.MustCompile("[^-a-z0-9]")
	return re.ReplaceAllString(name, "-")
}
