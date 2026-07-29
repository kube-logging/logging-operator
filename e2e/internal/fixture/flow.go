// Copyright © 2025 Kube logging authors
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

// Buffer defaults, counted across the suites: Timekey "1s" at eight sites
// against "1m" at three and "10s" at one; TimekeyWait "0s" at nine against
// "10s" at three.
const (
	DefaultTimekey     = "1s"
	DefaultTimekeyWait = "0s"
)

// ReceiverURL is the endpoint the test receiver serves, built the same way at
// nine call sites.
func ReceiverURL(release, tag string) string {
	return fmt.Sprintf("http://%s-test-receiver:8080/%s", release, tag)
}

// BufferOption modifies a Buffer under construction.
type BufferOption func(*output.Buffer)

// Buffer builds the file buffer the HTTP outputs share.
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

// Timekey overrides the buffer flush interval.
func Timekey(v string) BufferOption {
	return func(b *output.Buffer) { b.Timekey = v }
}

// TimekeyWait overrides the buffer flush delay.
func TimekeyWait(v string) BufferOption {
	return func(b *output.Buffer) { b.TimekeyWait = v }
}

// HTTPOutput builds a namespaced Output posting JSON to the test receiver.
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

// Flow builds a namespaced Flow selecting pods by label and routing to the
// named local outputs.
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
