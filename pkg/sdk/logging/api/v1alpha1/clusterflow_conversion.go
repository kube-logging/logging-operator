// Copyright © 2019 Banzai Cloud
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

package v1alpha1

import (
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (o *ClusterFlow) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.ClusterFlow)

	Log.Info("ConvertTo", "source", o.TypeMeta, "destination", dst.TypeMeta)

	dst.ObjectMeta = o.ObjectMeta
	dst.Spec = o.Spec
	dst.Status = o.Status

	return nil
}

func (o *ClusterFlow) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.ClusterFlow)

	Log.Info("ConvertFrom", "source", src.TypeMeta, "destination", o.TypeMeta)

	o.ObjectMeta = src.ObjectMeta
	o.Spec = src.Spec
	o.Status = src.Status

	return nil
}
