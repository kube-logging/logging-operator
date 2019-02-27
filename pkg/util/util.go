/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

// MergeLabels merges two map[string]string map
func MergeLabels(l map[string]string, l2 map[string]string) map[string]string {
	for lKey, lValue := range l2 {
		l[lKey] = lValue
	}
	return l
}

// IntPointer converts int32 to *int32
func IntPointer(i int32) *int32 {
	return &i
}
