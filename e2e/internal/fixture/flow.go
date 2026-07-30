// Copyright © 2026 Kube logging authors
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

package fixture

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	DefaultTimekey     = "1s"
	DefaultTimekeyWait = "0s"
)

func ReceiverURL(release, tag string) string {
	return fmt.Sprintf("http://%s-test-receiver:8080/%s", release, tag)
}

type BufferOption func(*output.Buffer)

func Buffer(tags string, opts ...BufferOption) *output.Buffer {
	b := &output.Buffer{
		Type:        "file",
		Tags:        &tags,
		Timekey:     DefaultTimekey,
		TimekeyWait: DefaultTimekeyWait,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func Timekey(v string) BufferOption {
	return func(b *output.Buffer) { b.Timekey = v }
}

func TimekeyWait(v string) BufferOption {
	return func(b *output.Buffer) { b.TimekeyWait = v }
}

func HTTPOutput(namespace, name, endpoint string, buffer *output.Buffer) *v1beta1.Output {
	return &v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: v1beta1.OutputSpec{
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    endpoint,
				ContentType: "application/json",
				Buffer:      buffer,
			},
		},
	}
}

func Flow(namespace, name string, selector map[string]string, outputRefs ...string) *v1beta1.Flow {
	return &v1beta1.Flow{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: v1beta1.FlowSpec{
			Match: []v1beta1.Match{
				{Select: &v1beta1.Select{Labels: selector}},
			},
			LocalOutputRefs: outputRefs,
		},
	}
}
