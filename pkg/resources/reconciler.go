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

package resources

import (
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ComponentReconciler reconciler interface
type ComponentReconciler func() (*reconcile.Result, error)

// Resource redeclaration of function with return type kubernetes Object
type Resource func() (runtime.Object, k8sutil.DesiredState)

// ResourceVariation redeclaration of function with parameter and return type kubernetes Object
type ResourceVariation func(t string) runtime.Object

// ResourceWithLog redeclaration of function with logging parameter and return type kubernetes Object
type ResourceWithLog func(log logr.Logger) runtime.Object
