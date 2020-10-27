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

package model

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func NewValidationReconciler(ctx context.Context, repo client.StatusClient, resources LoggingResources) func() (*reconcile.Result, error) {
	return func() (*reconcile.Result, error) {
		var patchRequests []patchRequest
		registerForPatching := func(obj runtime.Object) {
			patchRequests = append(patchRequests, patchRequest{
				Obj:   obj,
				Patch: client.MergeFrom(obj.DeepCopyObject()),
			})
		}

		for i := range resources.ClusterOutputs {
			output := &resources.ClusterOutputs[i]
			registerForPatching(output)

			output.Status.Active = false
		}

		for i := range resources.Outputs {
			output := &resources.Outputs[i]
			registerForPatching(output)

			output.Status.Active = false
		}

		for i := range resources.ClusterFlows {
			flow := &resources.ClusterFlows[i]
			registerForPatching(flow)

			flow.Status.Active = false
			flow.Status.Problems = nil

			if len(flow.Spec.GlobalOutputRefs) == 0 && len(flow.Spec.OutputRefs) > 0 {
				flow.Status.Problems = append(flow.Status.Problems, "\"outputRefs\" field is deprecated, use \"globalOutputRefs\" instead")
			}

			for _, ref := range flow.Spec.GlobalOutputRefs {
				if output := resources.ClusterOutputs.FindByName(ref); output != nil {
					flow.Status.Active = true
					output.Status.Active = true
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling global output reference: %s", ref))
				}
			}
		}

		for i := range resources.Flows {
			flow := &resources.Flows[i]
			registerForPatching(flow)

			flow.Status.Active = false
			flow.Status.Problems = nil

			if len(flow.Spec.LocalOutputRefs)+len(flow.Spec.GlobalOutputRefs) == 0 && len(flow.Spec.OutputRefs) > 0 {
				flow.Status.Problems = append(flow.Status.Problems, "\"outputRefs\" field is deprecated, use \"globalOutputRefs\" and \"localOutputRefs\" instead")
			}

			for _, ref := range flow.Spec.GlobalOutputRefs {
				if output := resources.ClusterOutputs.FindByName(ref); output != nil {
					flow.Status.Active = true
					output.Status.Active = true
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling global output reference: %s", ref))
				}
			}

			for _, ref := range flow.Spec.LocalOutputRefs {
				if output := resources.Outputs.FindByNamespacedName(flow.Namespace, ref); output != nil {
					flow.Status.Active = true
					output.Status.Active = true
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling local output reference: %s", ref))
				}
			}
		}

		var errs error
		for _, req := range patchRequests {
			if req.IsEmptyPatch() {
				continue
			}

			if err := repo.Status().Patch(ctx, req.Obj, req.Patch); err != nil {
				errs = errors.Append(errs, err)
			}
		}

		return nil, errs
	}
}

type patchRequest struct {
	Obj   runtime.Object
	Patch client.Patch
}

func (r patchRequest) IsEmptyPatch() bool {
	data, err := r.Patch.Data(r.Obj)
	return err == nil && string(data) == "{}"
}
