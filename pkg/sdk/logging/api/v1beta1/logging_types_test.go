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

package v1beta1_test

import (
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

// These tests are written in BDD-style using Ginkgo framework. Refer to
// http://onsi.github.io/ginkgo to learn more.
var _ = Describe("FluentBitSpec", func() {
	Context("When creating the default FluentBitSpec", func() {
		It("should overwrite DNSPolicy and HostNetwork", func() {
			When("Kubelet_Host is not defined", func() {
				spec := &v1beta1.FluentbitSpec{
					FilterKubernetes: v1beta1.FilterKubernetes{
						UseKubelet: "On",
					},
				}

				err := v1beta1.FluentBitDefaults(spec)
				Expect(err).To(BeNil())
				Expect(spec.FilterKubernetes.UseKubelet).To(Equal("On"))
				Expect(spec.FilterKubernetes.KubeletHost).To(Equal(""))
				Expect(spec.DNSPolicy).To(Equal(v1.DNSClusterFirstWithHostNet))
				Expect(spec.HostNetwork).To(Equal(true))
			})
		})

		It("should use the given DNSPolicy and HostNetwork", func() {
			When("Kubelet_Host is defined", func() {
				spec := &v1beta1.FluentbitSpec{
					HostNetwork: false,
					DNSPolicy:   v1.DNSClusterFirst,
					FilterKubernetes: v1beta1.FilterKubernetes{
						UseKubelet:  "On",
						KubeletHost: "${HOST_IP}",
					},
				}

				err := v1beta1.FluentBitDefaults(spec)
				Expect(err).To(BeNil())
				Expect(spec.FilterKubernetes.UseKubelet).To(Equal("On"))
				Expect(spec.FilterKubernetes.KubeletHost).To(Equal("${HOST_IP}"))
				Expect(spec.DNSPolicy).To(Equal(v1.DNSClusterFirst))
				Expect(spec.HostNetwork).To(Equal(false))
			})
		})
	})
})
