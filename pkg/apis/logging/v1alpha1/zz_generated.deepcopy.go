// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPlugin) DeepCopyInto(out *FPlugin) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make([]Parameter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPlugin.
func (in *FPlugin) DeepCopy() *FPlugin {
	if in == nil {
		return nil
	}
	out := new(FPlugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Fluentbit) DeepCopyInto(out *Fluentbit) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Fluentbit.
func (in *Fluentbit) DeepCopy() *Fluentbit {
	if in == nil {
		return nil
	}
	out := new(Fluentbit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Fluentbit) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentbitList) DeepCopyInto(out *FluentbitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Fluentbit, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentbitList.
func (in *FluentbitList) DeepCopy() *FluentbitList {
	if in == nil {
		return nil
	}
	out := new(FluentbitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FluentbitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentbitSpec) DeepCopyInto(out *FluentbitSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.Image = in.Image
	out.TLS = in.TLS
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentbitSpec.
func (in *FluentbitSpec) DeepCopy() *FluentbitSpec {
	if in == nil {
		return nil
	}
	out := new(FluentbitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentbitStatus) DeepCopyInto(out *FluentbitStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentbitStatus.
func (in *FluentbitStatus) DeepCopy() *FluentbitStatus {
	if in == nil {
		return nil
	}
	out := new(FluentbitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentbitTLS) DeepCopyInto(out *FluentbitTLS) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentbitTLS.
func (in *FluentbitTLS) DeepCopy() *FluentbitTLS {
	if in == nil {
		return nil
	}
	out := new(FluentbitTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Fluentd) DeepCopyInto(out *Fluentd) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Fluentd.
func (in *Fluentd) DeepCopy() *Fluentd {
	if in == nil {
		return nil
	}
	out := new(Fluentd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Fluentd) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdList) DeepCopyInto(out *FluentdList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Fluentd, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdList.
func (in *FluentdList) DeepCopy() *FluentdList {
	if in == nil {
		return nil
	}
	out := new(FluentdList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FluentdList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdSpec) DeepCopyInto(out *FluentdSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.TLS = in.TLS
	out.Image = in.Image
	in.FluentdPvcSpec.DeepCopyInto(&out.FluentdPvcSpec)
	out.VolumeModImage = in.VolumeModImage
	out.ConfigReloaderImage = in.ConfigReloaderImage
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdSpec.
func (in *FluentdSpec) DeepCopy() *FluentdSpec {
	if in == nil {
		return nil
	}
	out := new(FluentdSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdStatus) DeepCopyInto(out *FluentdStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdStatus.
func (in *FluentdStatus) DeepCopy() *FluentdStatus {
	if in == nil {
		return nil
	}
	out := new(FluentdStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdTLS) DeepCopyInto(out *FluentdTLS) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdTLS.
func (in *FluentdTLS) DeepCopy() *FluentdTLS {
	if in == nil {
		return nil
	}
	out := new(FluentdTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageSpec) DeepCopyInto(out *ImageSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageSpec.
func (in *ImageSpec) DeepCopy() *ImageSpec {
	if in == nil {
		return nil
	}
	out := new(ImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Input) DeepCopyInto(out *Input) {
	*out = *in
	if in.Label != nil {
		in, out := &in.Label, &out.Label
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Input.
func (in *Input) DeepCopy() *Input {
	if in == nil {
		return nil
	}
	out := new(Input)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesSecret) DeepCopyInto(out *KubernetesSecret) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesSecret.
func (in *KubernetesSecret) DeepCopy() *KubernetesSecret {
	if in == nil {
		return nil
	}
	out := new(KubernetesSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parameter) DeepCopyInto(out *Parameter) {
	*out = *in
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(ValueFrom)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parameter.
func (in *Parameter) DeepCopy() *Parameter {
	if in == nil {
		return nil
	}
	out := new(Parameter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Plugin) DeepCopyInto(out *Plugin) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Plugin.
func (in *Plugin) DeepCopy() *Plugin {
	if in == nil {
		return nil
	}
	out := new(Plugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Plugin) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginList) DeepCopyInto(out *PluginList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Plugin, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginList.
func (in *PluginList) DeepCopy() *PluginList {
	if in == nil {
		return nil
	}
	out := new(PluginList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PluginList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginSpec) DeepCopyInto(out *PluginSpec) {
	*out = *in
	in.Input.DeepCopyInto(&out.Input)
	if in.Filter != nil {
		in, out := &in.Filter, &out.Filter
		*out = make([]FPlugin, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Output != nil {
		in, out := &in.Output, &out.Output
		*out = make([]FPlugin, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginSpec.
func (in *PluginSpec) DeepCopy() *PluginSpec {
	if in == nil {
		return nil
	}
	out := new(PluginSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginStatus) DeepCopyInto(out *PluginStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginStatus.
func (in *PluginStatus) DeepCopy() *PluginStatus {
	if in == nil {
		return nil
	}
	out := new(PluginStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueFrom) DeepCopyInto(out *ValueFrom) {
	*out = *in
	out.SecretKeyRef = in.SecretKeyRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueFrom.
func (in *ValueFrom) DeepCopy() *ValueFrom {
	if in == nil {
		return nil
	}
	out := new(ValueFrom)
	in.DeepCopyInto(out)
	return out
}
