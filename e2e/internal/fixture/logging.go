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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// Logging is the smallest valid object and nothing more. ControlNamespace is
// the only field every suite sets; a spec built from anything else would be
// hiding a divergence rather than sharing a default. Logging is cluster-scoped,
// so controlNamespace fills Spec.ControlNamespace and not ObjectMeta.Namespace.
func Logging(controlNamespace, name string) *v1beta1.Logging {
	return &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       v1beta1.LoggingSpec{ControlNamespace: controlNamespace},
	}
}
