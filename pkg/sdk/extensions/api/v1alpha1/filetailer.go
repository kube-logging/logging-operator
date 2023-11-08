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

package v1alpha1

import (
	"fmt"

	"github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/tailer"
	config "github.com/kube-logging/logging-operator/pkg/sdk/extensions/extensionsconfig"
)

func (f FileTailer) defaults() FileTailer {
	result := f
	// setting defaults
	if result.BufferMaxSize == "" {
		result.BufferMaxSize = "32k"
	}
	if result.BufferChunkSize == "" {
		result.BufferChunkSize = result.BufferMaxSize
	}

	if result.SkipLongLines == "" {
		result.SkipLongLines = "On"
	}

	return result
}

// Command returns the desired command for the current filetailer
func (f FileTailer) Command(Name string) []string {
	f = f.defaults()
	command := []string{
		"/fluent-bit/bin/fluent-bit", "-i", "tail",
		"-p", fmt.Sprintf("path=%s", f.Path),
		"-p", fmt.Sprintf("db=/var/pos/%s.db", Name),
		"-p", fmt.Sprintf("buffer_chunk_size=%s", f.BufferChunkSize),
		"-p", fmt.Sprintf("buffer_max_size=%s", f.BufferMaxSize),
		"-p", fmt.Sprintf("skip_long_lines=%s", f.SkipLongLines),
		"-p", fmt.Sprintf("read_from_head=%t", f.ReadFromHead),
		"-o", "file",
		"-p", "format=template",
		"-p", "template={log}",
	}
	command = append(command, config.HostTailer.VersionedFluentBitPathArgs("/dev/stdout")...)
	return command
}

// GeneralDescriptor returns the tailer.General general Tailer struct
func (f FileTailer) GeneralDescriptor() tailer.General {
	return tailer.General{Name: f.Name, Path: f.Path, Disabled: f.Disabled, ContainerBase: f.ContainerBase, Image: f.Image}
}
