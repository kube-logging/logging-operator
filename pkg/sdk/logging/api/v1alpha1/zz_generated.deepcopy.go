//go:build !ignore_autogenerated

// Copyright © 2020 Banzai Cloud
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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFlow) DeepCopyInto(out *ClusterFlow) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFlow.
func (in *ClusterFlow) DeepCopy() *ClusterFlow {
	if in == nil {
		return nil
	}
	out := new(ClusterFlow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterFlow) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFlowList) DeepCopyInto(out *ClusterFlowList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterFlow, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFlowList.
func (in *ClusterFlowList) DeepCopy() *ClusterFlowList {
	if in == nil {
		return nil
	}
	out := new(ClusterFlowList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterFlowList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterOutput) DeepCopyInto(out *ClusterOutput) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterOutput.
func (in *ClusterOutput) DeepCopy() *ClusterOutput {
	if in == nil {
		return nil
	}
	out := new(ClusterOutput)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterOutput) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterOutputList) DeepCopyInto(out *ClusterOutputList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterOutput, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterOutputList.
func (in *ClusterOutputList) DeepCopy() *ClusterOutputList {
	if in == nil {
		return nil
	}
	out := new(ClusterOutputList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterOutputList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Flow) DeepCopyInto(out *Flow) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Flow.
func (in *Flow) DeepCopy() *Flow {
	if in == nil {
		return nil
	}
	out := new(Flow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Flow) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlowList) DeepCopyInto(out *FlowList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Flow, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowList.
func (in *FlowList) DeepCopy() *FlowList {
	if in == nil {
		return nil
	}
	out := new(FlowList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlowList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Logging) DeepCopyInto(out *Logging) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Logging.
func (in *Logging) DeepCopy() *Logging {
	if in == nil {
		return nil
	}
	out := new(Logging)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Logging) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingList) DeepCopyInto(out *LoggingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Logging, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingList.
func (in *LoggingList) DeepCopy() *LoggingList {
	if in == nil {
		return nil
	}
	out := new(LoggingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoggingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingSpec) DeepCopyInto(out *LoggingSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingSpec.
func (in *LoggingSpec) DeepCopy() *LoggingSpec {
	if in == nil {
		return nil
	}
	out := new(LoggingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingStatus) DeepCopyInto(out *LoggingStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingStatus.
func (in *LoggingStatus) DeepCopy() *LoggingStatus {
	if in == nil {
		return nil
	}
	out := new(LoggingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Output) DeepCopyInto(out *Output) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Output.
func (in *Output) DeepCopy() *Output {
	if in == nil {
		return nil
	}
	out := new(Output)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Output) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputList) DeepCopyInto(out *OutputList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputList.
func (in *OutputList) DeepCopy() *OutputList {
	if in == nil {
		return nil
	}
	out := new(OutputList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OutputList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputSpec) DeepCopyInto(out *OutputSpec) {
	*out = *in
	if in.S3OutputConfig != nil {
		in, out := &in.S3OutputConfig, &out.S3OutputConfig
		*out = new(output.S3OutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.AzureStorage != nil {
		in, out := &in.AzureStorage, &out.AzureStorage
		*out = new(output.AzureStorage)
		(*in).DeepCopyInto(*out)
	}
	if in.GCSOutput != nil {
		in, out := &in.GCSOutput, &out.GCSOutput
		*out = new(output.GCSOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.OSSOutput != nil {
		in, out := &in.OSSOutput, &out.OSSOutput
		*out = new(output.OSSOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.ElasticsearchOutput != nil {
		in, out := &in.ElasticsearchOutput, &out.ElasticsearchOutput
		*out = new(output.ElasticsearchOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.OpenSearchOutput != nil {
		in, out := &in.OpenSearchOutput, &out.OpenSearchOutput
		*out = new(output.OpenSearchOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.LogZOutput != nil {
		in, out := &in.LogZOutput, &out.LogZOutput
		*out = new(output.LogZOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.LokiOutput != nil {
		in, out := &in.LokiOutput, &out.LokiOutput
		*out = new(output.LokiOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.DatadogOutput != nil {
		in, out := &in.DatadogOutput, &out.DatadogOutput
		*out = new(output.DatadogOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.ForwardOutput != nil {
		in, out := &in.ForwardOutput, &out.ForwardOutput
		*out = new(output.ForwardOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.FileOutput != nil {
		in, out := &in.FileOutput, &out.FileOutput
		*out = new(output.FileOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.NullOutputConfig != nil {
		in, out := &in.NullOutputConfig, &out.NullOutputConfig
		*out = new(output.NullOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.KafkaOutputConfig != nil {
		in, out := &in.KafkaOutputConfig, &out.KafkaOutputConfig
		*out = new(output.KafkaOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.CloudWatchOutput != nil {
		in, out := &in.CloudWatchOutput, &out.CloudWatchOutput
		*out = new(output.CloudWatchOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.KinesisStreamOutputConfig != nil {
		in, out := &in.KinesisStreamOutputConfig, &out.KinesisStreamOutputConfig
		*out = new(output.KinesisStreamOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.LogDNAOutput != nil {
		in, out := &in.LogDNAOutput, &out.LogDNAOutput
		*out = new(output.LogDNAOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.NewRelicOutputConfig != nil {
		in, out := &in.NewRelicOutputConfig, &out.NewRelicOutputConfig
		*out = new(output.NewRelicOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.SplunkHecOutput != nil {
		in, out := &in.SplunkHecOutput, &out.SplunkHecOutput
		*out = new(output.SplunkHecOutput)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPOutput != nil {
		in, out := &in.HTTPOutput, &out.HTTPOutput
		*out = new(output.HTTPOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.AwsElasticsearchOutputConfig != nil {
		in, out := &in.AwsElasticsearchOutputConfig, &out.AwsElasticsearchOutputConfig
		*out = new(output.AwsElasticsearchOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.RedisOutputConfig != nil {
		in, out := &in.RedisOutputConfig, &out.RedisOutputConfig
		*out = new(output.RedisOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.SyslogOutputConfig != nil {
		in, out := &in.SyslogOutputConfig, &out.SyslogOutputConfig
		*out = new(output.SyslogOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.GelfOutputConfig != nil {
		in, out := &in.GelfOutputConfig, &out.GelfOutputConfig
		*out = new(output.GelfOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.SQSOutputConfig != nil {
		in, out := &in.SQSOutputConfig, &out.SQSOutputConfig
		*out = new(output.SQSOutputConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.RelabelOutputConfig != nil {
		in, out := &in.RelabelOutputConfig, &out.RelabelOutputConfig
		*out = new(output.RelabelOutputConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputSpec.
func (in *OutputSpec) DeepCopy() *OutputSpec {
	if in == nil {
		return nil
	}
	out := new(OutputSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputStatus) DeepCopyInto(out *OutputStatus) {
	*out = *in
	if in.Active != nil {
		in, out := &in.Active, &out.Active
		*out = new(bool)
		**out = **in
	}
	if in.Problems != nil {
		in, out := &in.Problems, &out.Problems
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputStatus.
func (in *OutputStatus) DeepCopy() *OutputStatus {
	if in == nil {
		return nil
	}
	out := new(OutputStatus)
	in.DeepCopyInto(out)
	return out
}
