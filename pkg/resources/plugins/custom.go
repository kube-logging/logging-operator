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

package plugins

// CustomOutput CRD name
const CustomOutput = "custom"

// CustomDefaultValues for custom output plugin
var CustomDefaultValues = map[string]string{
	"type":   "",
	"output": "",
}

// CustomTemplate for custom output plugin
const CustomTemplate = `
  <match {{ .pattern }}.** >
  @type {{ .type }}
  {{ .output }}
  </match>`
