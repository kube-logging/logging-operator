// Copyright © 2022 Banzai Cloud
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

package output

// +kubebuilder:object:generate=true
// Documentation: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829124
type FileOutput struct {
	Path       string      `json:"path"`
	CreateDirs bool        `json:"create_dirs,omitempty"`
	DirGroup   string      `json:"dir_group,omitempty"`
	DirOwner   string      `json:"dir_owner,omitempty"`
	DiskBuffer *DiskBuffer `json:"disk_buffer,omitempty"`
}
