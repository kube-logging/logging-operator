// +build !ignore_autogenerated

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

package filter

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AndSection) DeepCopyInto(out *AndSection) {
	*out = *in
	if in.Regexp != nil {
		in, out := &in.Regexp, &out.Regexp
		*out = make([]RegexpSection, len(*in))
		copy(*out, *in)
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]ExcludeSection, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AndSection.
func (in *AndSection) DeepCopy() *AndSection {
	if in == nil {
		return nil
	}
	out := new(AndSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ColorStripper) DeepCopyInto(out *ColorStripper) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ColorStripper.
func (in *ColorStripper) DeepCopy() *ColorStripper {
	if in == nil {
		return nil
	}
	out := new(ColorStripper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Concat) DeepCopyInto(out *Concat) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Concat.
func (in *Concat) DeepCopy() *Concat {
	if in == nil {
		return nil
	}
	out := new(Concat)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DedotFilterConfig) DeepCopyInto(out *DedotFilterConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DedotFilterConfig.
func (in *DedotFilterConfig) DeepCopy() *DedotFilterConfig {
	if in == nil {
		return nil
	}
	out := new(DedotFilterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DetectExceptions) DeepCopyInto(out *DetectExceptions) {
	*out = *in
	if in.Languages != nil {
		in, out := &in.Languages, &out.Languages
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DetectExceptions.
func (in *DetectExceptions) DeepCopy() *DetectExceptions {
	if in == nil {
		return nil
	}
	out := new(DetectExceptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExcludeSection) DeepCopyInto(out *ExcludeSection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExcludeSection.
func (in *ExcludeSection) DeepCopy() *ExcludeSection {
	if in == nil {
		return nil
	}
	out := new(ExcludeSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeoIP) DeepCopyInto(out *GeoIP) {
	*out = *in
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]Record, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(Record, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeoIP.
func (in *GeoIP) DeepCopy() *GeoIP {
	if in == nil {
		return nil
	}
	out := new(GeoIP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrepConfig) DeepCopyInto(out *GrepConfig) {
	*out = *in
	if in.Regexp != nil {
		in, out := &in.Regexp, &out.Regexp
		*out = make([]RegexpSection, len(*in))
		copy(*out, *in)
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]ExcludeSection, len(*in))
		copy(*out, *in)
	}
	if in.Or != nil {
		in, out := &in.Or, &out.Or
		*out = make([]OrSection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.And != nil {
		in, out := &in.And, &out.And
		*out = make([]AndSection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrepConfig.
func (in *GrepConfig) DeepCopy() *GrepConfig {
	if in == nil {
		return nil
	}
	out := new(GrepConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricSection) DeepCopyInto(out *MetricSection) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(Label, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricSection.
func (in *MetricSection) DeepCopy() *MetricSection {
	if in == nil {
		return nil
	}
	out := new(MetricSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OrSection) DeepCopyInto(out *OrSection) {
	*out = *in
	if in.Regexp != nil {
		in, out := &in.Regexp, &out.Regexp
		*out = make([]RegexpSection, len(*in))
		copy(*out, *in)
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]ExcludeSection, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OrSection.
func (in *OrSection) DeepCopy() *OrSection {
	if in == nil {
		return nil
	}
	out := new(OrSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParseSection) DeepCopyInto(out *ParseSection) {
	*out = *in
	if in.Multiline != nil {
		in, out := &in.Multiline, &out.Multiline
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Patterns != nil {
		in, out := &in.Patterns, &out.Patterns
		*out = make([]SingleParseSection, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParseSection.
func (in *ParseSection) DeepCopy() *ParseSection {
	if in == nil {
		return nil
	}
	out := new(ParseSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParserConfig) DeepCopyInto(out *ParserConfig) {
	*out = *in
	in.Parse.DeepCopyInto(&out.Parse)
	if in.Parsers != nil {
		in, out := &in.Parsers, &out.Parsers
		*out = make([]ParseSection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParserConfig.
func (in *ParserConfig) DeepCopy() *ParserConfig {
	if in == nil {
		return nil
	}
	out := new(ParserConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusConfig) DeepCopyInto(out *PrometheusConfig) {
	*out = *in
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make([]MetricSection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(Label, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusConfig.
func (in *PrometheusConfig) DeepCopy() *PrometheusConfig {
	if in == nil {
		return nil
	}
	out := new(PrometheusConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordModifier) DeepCopyInto(out *RecordModifier) {
	*out = *in
	if in.Replaces != nil {
		in, out := &in.Replaces, &out.Replaces
		*out = make([]Replace, len(*in))
		copy(*out, *in)
	}
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]Record, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(Record, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordModifier.
func (in *RecordModifier) DeepCopy() *RecordModifier {
	if in == nil {
		return nil
	}
	out := new(RecordModifier)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordTransformer) DeepCopyInto(out *RecordTransformer) {
	*out = *in
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]Record, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(Record, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordTransformer.
func (in *RecordTransformer) DeepCopy() *RecordTransformer {
	if in == nil {
		return nil
	}
	out := new(RecordTransformer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegexpSection) DeepCopyInto(out *RegexpSection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegexpSection.
func (in *RegexpSection) DeepCopy() *RegexpSection {
	if in == nil {
		return nil
	}
	out := new(RegexpSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replace) DeepCopyInto(out *Replace) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replace.
func (in *Replace) DeepCopy() *Replace {
	if in == nil {
		return nil
	}
	out := new(Replace)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SingleParseSection) DeepCopyInto(out *SingleParseSection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SingleParseSection.
func (in *SingleParseSection) DeepCopy() *SingleParseSection {
	if in == nil {
		return nil
	}
	out := new(SingleParseSection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StdOutFilterConfig) DeepCopyInto(out *StdOutFilterConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StdOutFilterConfig.
func (in *StdOutFilterConfig) DeepCopy() *StdOutFilterConfig {
	if in == nil {
		return nil
	}
	out := new(StdOutFilterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Throttle) DeepCopyInto(out *Throttle) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Throttle.
func (in *Throttle) DeepCopy() *Throttle {
	if in == nil {
		return nil
	}
	out := new(Throttle)
	in.DeepCopyInto(out)
	return out
}
