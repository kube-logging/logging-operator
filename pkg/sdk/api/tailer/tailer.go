// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package tailer

import "github.com/banzaicloud/operator-tools/pkg/types"

// Tailer .
type Tailer interface {
	Command(Name string) []string
	GeneralDescriptor() General
}

// +kubebuilder:object:generate=true

// General descriptor for hosttailers
type General struct {
	Name          string               `json:"name"`
	Path          string               `json:"path,omitempty"`
	Disabled      bool                 `json:"disabled,omitempty"`
	ContainerBase *types.ContainerBase `json:"containerOverrides,omitempty"`
}
