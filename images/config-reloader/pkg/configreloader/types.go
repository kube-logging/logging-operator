// Copyright © 2025 Kube logging authors
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

package configreloader

import (
	"fmt"
	"net/url"
)

type ConfigReloader struct {
	InitMode          *bool
	VolumeDirs        volumeDirsFlag
	VolumeDirsArchive volumeDirsArchiveFlag
	DirForUnarchive   *string
	Webhook           Webhook
}

type Webhook struct {
	Urls       urlsFlag
	Method     *string
	StatusCode *int
	Retries    *int
}

type volumeDirsFlag []string
type volumeDirsArchiveFlag []string
type urlsFlag []*url.URL

func (v *volumeDirsFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}
func (v *volumeDirsFlag) String() string {
	return fmt.Sprint(*v)
}

func (v *volumeDirsArchiveFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}
func (v *volumeDirsArchiveFlag) String() string {
	return fmt.Sprint(*v)
}

func (v *urlsFlag) Set(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	*v = append(*v, u)
	return nil
}

func (v *urlsFlag) String() string {
	return fmt.Sprint(*v)
}
