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

package eventtailer

import (
	"k8s.io/apimachinery/pkg/runtime"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/go-logr/logr"
	loggingextensionsv1alpha1 "github.com/kube-logging/logging-operator/pkg/sdk/extensions/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EventTailer .
type EventTailer struct {
	log logr.Logger
	*reconciler.GenericResourceReconciler
	customResource       loggingextensionsv1alpha1.EventTailer
	CommonSelectorLabels map[string]string `json:"selectorLabels,omitempty"`
}

// New .
func New(client client.Client, log logr.Logger, opts reconciler.ReconcilerOpts, customResource loggingextensionsv1alpha1.EventTailer) *EventTailer {
	return &EventTailer{
		log:                       log,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, log, opts),
		customResource:            customResource,
	}
}

func (e *EventTailer) Reconcile(object runtime.Object) (*reconcile.Result, error) {
	for _, res := range []reconciler.ResourceBuilder{
		e.ServiceAccount,
		e.ClusterRole,
		e.ClusterRoleBinding,
		e.ConfigMap,
		e.StatefulSet,
	} {
		o, state, err := res()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", res)
		}
		result, err := e.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

// RegisterWatches completes the implementation of ComponentReconciler
func (e *EventTailer) RegisterWatches(*builder.Builder) {
	// placeholder
}
