// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package kubetool

import (
	corev1 "k8s.io/api/core/v1"
)

// VolumeBuilder .
type VolumeBuilder struct {
	corev1.Volume
}

// NewVolumeBuilder .
func NewVolumeBuilder() *VolumeBuilder {
	return &VolumeBuilder{}
}

// WithName .
func (v *VolumeBuilder) WithName(name string) *VolumeBuilder {
	if v != nil {
		v.Name = name
	}
	return v
}

// WithVolumeSource .
func (v *VolumeBuilder) WithVolumeSource(volumeSource corev1.VolumeSource) *VolumeBuilder {
	if v != nil {
		v.VolumeSource = volumeSource
	}
	return v
}

// WithEmptyDir .
func (v *VolumeBuilder) WithEmptyDir(emptyDir corev1.EmptyDirVolumeSource) *VolumeBuilder {
	if v != nil {
		v.VolumeSource.EmptyDir = &emptyDir
	}
	return v
}

// WithHostPath .
func (v *VolumeBuilder) WithHostPath(hostPath corev1.HostPathVolumeSource) *VolumeBuilder {
	if v != nil {
		v.VolumeSource.HostPath = &hostPath
	}
	return v
}

// WithHostPathFromPath .
func (v *VolumeBuilder) WithHostPathFromPath(path string) *VolumeBuilder {
	if v != nil {
		v.WithHostPath(corev1.HostPathVolumeSource{Path: path})
	}
	return v
}
