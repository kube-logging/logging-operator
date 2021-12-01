// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package v1alpha1

import (
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/api/tailer"
	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensionsconfig"
)

// Command returns the desired command for the current filetailer
func (f FileTailer) Command(Name string) []string {
	command := []string{
		"/fluent-bit/bin/fluent-bit", "-i", "tail",
		"-p", fmt.Sprintf("path=%s", f.Path),
		"-p", fmt.Sprintf("db=/var/pos/%s.db", Name),
		"-o", "file",
		"-p", "format=template",
		"-p", "template={log}",
	}
	command = append(command, config.HostTailer.VersionedFluentBitPathArgs("/dev/stdout")...)
	return command
}

// GeneralDescriptor returns the tailer.General general Tailer struct
func (f FileTailer) GeneralDescriptor() tailer.General {
	return tailer.General{Name: f.Name, Path: f.Path, Disabled: f.Disabled, ContainerBase: f.ContainerBase}
}
