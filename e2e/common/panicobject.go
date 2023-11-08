// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package common

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type PanicObject struct{}

var _ runtime.Object = (*PanicObject)(nil)

func (*PanicObject) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (*PanicObject) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

var _ metav1.Object = (*PanicObject)(nil)

func (*PanicObject) GetNamespace() string {
	panic("not implemented")
}
func (*PanicObject) SetNamespace(namespace string) {
	panic("not implemented")
}
func (*PanicObject) GetName() string {
	panic("not implemented")
}
func (*PanicObject) SetName(name string) {
	panic("not implemented")
}
func (*PanicObject) GetGenerateName() string {
	panic("not implemented")
}
func (*PanicObject) SetGenerateName(name string) {
	panic("not implemented")
}
func (*PanicObject) GetUID() types.UID {
	panic("not implemented")
}
func (*PanicObject) SetUID(uid types.UID) {
	panic("not implemented")
}
func (*PanicObject) GetResourceVersion() string {
	panic("not implemented")
}
func (*PanicObject) SetResourceVersion(version string) {
	panic("not implemented")
}
func (*PanicObject) GetGeneration() int64 {
	panic("not implemented")
}
func (*PanicObject) SetGeneration(generation int64) {
	panic("not implemented")
}
func (*PanicObject) GetSelfLink() string {
	panic("not implemented")
}
func (*PanicObject) SetSelfLink(selfLink string) {
	panic("not implemented")
}
func (*PanicObject) GetCreationTimestamp() metav1.Time {
	panic("not implemented")
}
func (*PanicObject) SetCreationTimestamp(timestamp metav1.Time) {
	panic("not implemented")
}
func (*PanicObject) GetDeletionTimestamp() *metav1.Time {
	panic("not implemented")
}
func (*PanicObject) SetDeletionTimestamp(timestamp *metav1.Time) {
	panic("not implemented")
}
func (*PanicObject) GetDeletionGracePeriodSeconds() *int64 {
	panic("not implemented")
}
func (*PanicObject) SetDeletionGracePeriodSeconds(*int64) {
	panic("not implemented")
}
func (*PanicObject) GetLabels() map[string]string {
	panic("not implemented")
}
func (*PanicObject) SetLabels(labels map[string]string) {
	panic("not implemented")
}
func (*PanicObject) GetAnnotations() map[string]string {
	panic("not implemented")
}
func (*PanicObject) SetAnnotations(annotations map[string]string) {
	panic("not implemented")
}
func (*PanicObject) GetFinalizers() []string {
	panic("not implemented")
}
func (*PanicObject) SetFinalizers(finalizers []string) {
	panic("not implemented")
}
func (*PanicObject) GetOwnerReferences() []metav1.OwnerReference {
	panic("not implemented")
}
func (*PanicObject) SetOwnerReferences([]metav1.OwnerReference) {
	panic("not implemented")
}
func (*PanicObject) GetZZZ_DeprecatedClusterName() string {
	panic("not implemented")
}
func (*PanicObject) SetZZZ_DeprecatedClusterName(clusterName string) {
	panic("not implemented")
}
func (*PanicObject) GetManagedFields() []metav1.ManagedFieldsEntry {
	panic("not implemented")
}
func (*PanicObject) SetManagedFields(managedFields []metav1.ManagedFieldsEntry) {
	panic("not implemented")
}
