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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Hub marks these types as conversion hub.
func (r *Logging) Hub()       {}
func (r *Output) Hub()        {}
func (r *ClusterOutput) Hub() {}
func (r *Flow) Hub()          {}
func (r *ClusterFlow) Hub()   {}

func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, apiType := range []runtime.Object{&Logging{}, &Output{}, &ClusterOutput{}, &Flow{}, &ClusterFlow{}} {
		// register webhook using controller-runtime because of interface checks
		if err := ctrl.NewWebhookManagedBy(mgr).
			For(apiType).
			Complete(); err != nil {
			return err
		}
	}
	return nil
}
