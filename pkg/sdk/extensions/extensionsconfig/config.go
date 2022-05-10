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

package extensionsconfig

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// GlobalConfig is a configuration type for global configurations
type GlobalConfig struct {
	FluentBitPosFilePath   string
	FluentBitPosVolumeName string
	OperatorImage          string
	ContainerRuntime       string
}

// HostTailerConfig is a configuration type for HostTailer
type HostTailerConfig struct {
	FluentBitImage string
	TailerAffix    string
}

// VersionedFluentBitPathArgs returns fluent-bit config path/file args in particular format chosen by the image version
func (t HostTailerConfig) VersionedFluentBitPathArgs(filePath string) []string {
	return fluentBitConfigFilePath(t.FluentBitImage, filePath)
}

// EventTailerConfig is a configuration type for EventTailer
type EventTailerConfig struct {
	TailerAffix           string
	ConfigurationFileName string
	PositionVolumeName    string
}

// TailerWebhookConfig is a configuration type for TailerWebhook
type TailerWebhookConfig struct {
	FluentBitImage    string
	AnnotationKey     string
	ServerPath        string
	ServerPort        int
	CertDir           string
	DisableEnvVarName string
}

// VersionedFluentBitPathArgs returns fluent-bit config path/file args in particular format chosen by the image version
func (t TailerWebhookConfig) VersionedFluentBitPathArgs(filePath string) []string {
	return fluentBitConfigFilePath(t.FluentBitImage, filePath)
}

// Global configuration
var Global = GlobalConfig{
	FluentBitPosFilePath:   "/var/pos",
	FluentBitPosVolumeName: "positions",
	OperatorImage:          "033498657557.dkr.ecr.us-east-2.amazonaws.com/banzaicloud/logging-extensions:0.2.0",
}

// HostTailer configuration
var HostTailer = HostTailerConfig{
	FluentBitImage: "fluent/fluent-bit:1.9.3",
	TailerAffix:    "host-tailer",
}

// EventTailer configuration
var EventTailer = EventTailerConfig{
	TailerAffix:           "event-tailer",
	ConfigurationFileName: "config.json",
	PositionVolumeName:    "event-tailer-position",
}

// TailerWebhook configuration
var TailerWebhook = TailerWebhookConfig{
	FluentBitImage:    "fluent/fluent-bit:1.9.3",
	AnnotationKey:     "sidecar.logging-extensions.banzaicloud.io/tail",
	ServerPath:        "/tailer-webhook",
	ServerPort:        9443,
	CertDir:           "/tmp/k8s-webhook-server/serving-certs",
	DisableEnvVarName: "ENABLE_TAILER_WEBHOOK",
}

// FLuentBitFilePathBreakingChangeVersion holds the version where the fluent-bit command arguments are changed
const FLuentBitFilePathBreakingChangeVersion = "1.4.6"

// fluentBitConfigFilePath returns fluent-bit config args string array in the particular format depends by the version of image
func fluentBitConfigFilePath(image, filePath string) []string {
	if filePath == "" {
		return []string{}
	}
	var v0, v1 *semver.Version
	var err error
	fluentBitImageVersion := image[strings.LastIndex(image, ":")+1:]
	v0, err = semver.NewVersion(FLuentBitFilePathBreakingChangeVersion)
	if err != nil {
		goto defaultPath
	}
	v1, err = semver.NewVersion(fluentBitImageVersion)
	if err != nil {
		goto defaultPath
	}

	if v1.Equal(v0) || v1.GreaterThan(v0) {
		dir, file := filepath.Split(filePath)
		if dir == "" {
			goto defaultPath
		}
		return []string{
			"-p", fmt.Sprintf("path=%s", dir),
			"-p", fmt.Sprintf("file=%s", file),
		}
	}

defaultPath:
	return []string{"-p", fmt.Sprintf("path=%s", filePath)}
}
