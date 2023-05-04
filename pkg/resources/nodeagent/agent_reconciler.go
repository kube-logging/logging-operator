// Copyright © 2023 Kube logging authors
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

package nodeagent

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ResourceReconciler interface {
	ReconcileResource(desired runtime.Object, desiredState reconciler.DesiredState) (*reconcile.Result, error)
}

type GenericAgentReconciler struct {
	resourceReconciler ResourceReconciler
	dataProvider       AgentDataProvider
	logger             logr.Logger
}

func NewGenericAgentReconciler(
	resourceReconciler ResourceReconciler,
	dataProvider AgentDataProvider,
	logger logr.Logger) *GenericAgentReconciler {
	return &GenericAgentReconciler{
		resourceReconciler: resourceReconciler,
		dataProvider:       dataProvider,
		logger:             logger,
	}
}

// Reconcile is responsible to handle each object, created by ResourceBuilders according to the desired state.
// In case an object is nil, Reconcile will skip and move on
// In case a ResourceBuilder returns an error, Reconcile will report it then and move on
func (a *GenericAgentReconciler) Reconcile(resources []reconciler.ResourceBuilder) (reconcile.Result, error) {
	result := reconciler.CombinedResult{}
	for _, factory := range resources {
		o, state, err := factory()
		if err != nil {
			result.CombineErr(errors.WrapIf(err, "failed to create desired object"))
			continue
		}
		if o == nil {
			a.logger.Info(fmt.Sprintf("Resource not implemented. Resource %#v returns with nil object", factory))
			continue
		}
		resourceResult, err := a.resourceReconciler.ReconcileResource(o, state)
		if err != nil {
			result.CombineErr(errors.WrapWithDetails(err,
				"failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind()))
		}
		if resourceResult != nil {
			result.Combine(resourceResult, nil)
		}
	}
	return result.Result, result.Err
}

func (a *GenericAgentReconciler) ChildObjectMeta(resource client.Object) error {
	resource.SetNamespace(a.dataProvider.Namespace())
	resource.SetLabels(util.MergeLabels(
		// get logging agent specific labels
		a.dataProvider.ResourceLabels(),
		// add instance label with the name of the parent
		map[string]string{},
	))
	resource.SetOwnerReferences(a.dataProvider.OwnerRefs())
	return nil
}
