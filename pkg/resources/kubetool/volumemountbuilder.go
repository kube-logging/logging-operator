// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package kubetool

import (
	corev1 "k8s.io/api/core/v1"
)

// VolumeMountBuilder .
type VolumeMountBuilder struct {
	corev1.VolumeMount
}

// MountPropagationModeRef .
func MountPropagationModeRef(mountPropagationMode corev1.MountPropagationMode) *corev1.MountPropagationMode {
	return &mountPropagationMode
}

// NewVolumeMountBuilder .
func NewVolumeMountBuilder() *VolumeMountBuilder {
	return &VolumeMountBuilder{}
}

// WithName .
func (v *VolumeMountBuilder) WithName(name string) *VolumeMountBuilder {
	if v != nil {
		v.Name = name
	}
	return v
}

// WithMountPath .
func (v *VolumeMountBuilder) WithMountPath(path string) *VolumeMountBuilder {
	if v != nil {
		v.MountPath = path
	}
	return v
}

// WithSubPath .
func (v *VolumeMountBuilder) WithSubPath(path string) *VolumeMountBuilder {
	if v != nil {
		v.SubPath = path
	}
	return v
}

// WithSubPathExpr .
func (v *VolumeMountBuilder) WithSubPathExpr(expr string) *VolumeMountBuilder {
	if v != nil {
		v.SubPathExpr = expr
	}
	return v
}

// WithMountPropagation .
func (v *VolumeMountBuilder) WithMountPropagation(mountPropagation corev1.MountPropagationMode) *VolumeMountBuilder {
	if v != nil {
		v.MountPropagation = &mountPropagation
	}
	return v
}

// WithReadOnly .
func (v *VolumeMountBuilder) WithReadOnly(readonly bool) *VolumeMountBuilder {
	if v != nil {
		v.ReadOnly = readonly
	}
	return v
}
