// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package v1alpha1

import (
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/api/tailer"
	config "github.com/banzaicloud/logging-operator/pkg/sdk/extensionsconfig"
)

func (s SystemdTailer) defaults() SystemdTailer {
	result := s
	// setting defaults
	if result.Path == "" {
		result.Path = "/var/log/journal"
	}
	if result.MaxEntries == 0 {
		result.MaxEntries = 1000
	}
	return result
}

// Command returns the desired command for the current systemdtailer
func (s SystemdTailer) Command(Name string) []string {
	s = s.defaults()
	command := []string{
		"/fluent-bit/bin/fluent-bit", "-i", "systemd",
		"-p", fmt.Sprintf("path=%s", s.Path),
		"-p", fmt.Sprintf("db=/var/pos/%s.db", Name),
		"-p", fmt.Sprintf("max_entries=%d", s.MaxEntries),
	}
	if s.SystemdFilter != "" {
		command = append(command, "-p", fmt.Sprintf("systemd_filter=_SYSTEMD_UNIT=%s", s.SystemdFilter))
	}
	command = append(command,
		"-o", "file",
		"-p", "format=plain",
	)
	command = append(command, config.HostTailer.VersionedFluentBitPathArgs("/dev/stdout")...)
	return command
}

// GeneralDescriptor returns the tailer.General general Tailer struct
func (s SystemdTailer) GeneralDescriptor() tailer.General {
	s = s.defaults()
	return tailer.General{Name: s.Name, Path: s.Path, Disabled: s.Disabled, ContainerBase: s.ContainerBase}
}
