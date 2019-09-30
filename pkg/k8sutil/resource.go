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

package k8sutil

import (
	"context"
	"reflect"

	"emperror.dev/errors"
	"github.com/banzaicloud/k8s-objectmatcher/patch"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	StateAbsent  DesiredState = "Absent"
	StatePresent DesiredState = "Present"
)

type DesiredState string

// GenericResourceReconciler generic resource reconciler
type GenericResourceReconciler struct {
	Log    logr.Logger
	Client runtimeClient.Client
}

// NewReconciler returns GenericResourceReconciler
func NewReconciler(client runtimeClient.Client, log logr.Logger) *GenericResourceReconciler {
	return &GenericResourceReconciler{
		Log:    log,
		Client: client,
	}
}

// CreateResource creates a resource if it doesn't exist
func (r *GenericResourceReconciler) CreateResource(desired runtime.Object) error {
	_, _, err := r.createIfNotExists(desired)
	return err
}

// ReconcileResource reconciles various kubernetes types
func (r *GenericResourceReconciler) ReconcileResource(desired runtime.Object, desiredState DesiredState) error {
	log := r.Log.WithValues("type", reflect.TypeOf(desired))

	switch desiredState {
	case StateAbsent:
		_, err := r.delete(desired)
		if err != nil {
			return errors.Wrapf(err, "failed to delete resource %+v", desired)
		}

	case StatePresent:
		created, current, err := r.createIfNotExists(desired)
		if err == nil && created {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to create resource %+v", desired)
		}
		key, err := runtimeClient.ObjectKeyFromObject(current)
		if err != nil {
			return errors.Wrapf(err, "meta accessor failed %+v", current)
		}
		if err == nil {
			patchResult, err := patch.DefaultPatchMaker.Calculate(current, desired)
			if err != nil {
				log.Error(err, "could not match objects",
					"kind", desired.GetObjectKind().GroupVersionKind(), "name", key.Name)
			} else if patchResult.IsEmpty() {
				log.V(1).Info("resource is in sync",
					"kind", desired.GetObjectKind().GroupVersionKind(), "name", key.Name)
				return nil
			} else {
				log.V(1).Info("resource diffs",
					"patch", string(patchResult.Patch),
					"current", string(patchResult.Current),
					"modified", string(patchResult.Modified),
					"original", string(patchResult.Original))
			}

			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(desired); err != nil {
				log.Error(err, "Failed to set last applied annotation", "desired", desired)
			}

			metaAccessor := meta.NewAccessor()

			currentResourceVersion, err := metaAccessor.ResourceVersion(current)
			if err != nil {
				return errors.Wrap(err, "failed to access resourceVersion from metadata")
			}
			metaAccessor.SetResourceVersion(desired, currentResourceVersion)

			var name string
			if name, err = metaAccessor.Name(current); err != nil {
				return errors.Wrap(err, "failed to access Name from metadata")
			}

			log.V(1).Info("Updating resource",
				"gvk", desired.GetObjectKind().GroupVersionKind(), "name", name)
			if err := r.Client.Update(context.TODO(), desired); err != nil {
				return emperror.WrapWith(err, "updating resource failed",
					"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
			}
			log.Info("resource updated", "resource", desired.GetObjectKind().GroupVersionKind())
		}

	}
	return nil
}

func (r *GenericResourceReconciler) createIfNotExists(desired runtime.Object) (bool, runtime.Object, error) {
	log := r.Log.WithValues("type", reflect.TypeOf(desired))

	var current = desired.DeepCopyObject()
	switch current.(type) {
	case *v1.ServiceMonitor:
		var crd apiextensions.CustomResourceDefinition
		o := runtimeClient.ObjectKey{
			Name: v1.SchemeGroupVersion.WithResource(v1.ServiceMonitorName).GroupResource().String(),
		}
		err := r.Client.Get(context.TODO(), o, &crd)
		if err != nil {
			return false, current, err
		}

	}
	key, err := runtimeClient.ObjectKeyFromObject(current)
	if err != nil {
		return false, nil, emperror.With(err)
	}
	err = r.Client.Get(context.TODO(), key, current)
	if err != nil && !apierrors.IsNotFound(err) {
		return false, nil, emperror.WrapWith(err, "getting resource failed",
			"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
	}
	if apierrors.IsNotFound(err) {
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(desired); err != nil {
			log.Error(err, "Failed to set last applied annotation", "desired", desired)
		}
		if err := r.Client.Create(context.TODO(), desired); err != nil {
			return false, nil, emperror.WrapWith(err, "creating resource failed",
				"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
		}
		log.Info("resource created", "resource", desired.GetObjectKind().GroupVersionKind())
		return true, current, nil
	}
	log.V(1).Info("resource already exists", "resource", desired.GetObjectKind().GroupVersionKind())
	return false, current, nil
}

func (r *GenericResourceReconciler) delete(desired runtime.Object) (bool, error) {
	log := r.Log.WithValues("type", reflect.TypeOf(desired))
	var current = desired.DeepCopyObject()
	switch current.(type) {
	case *v1.ServiceMonitor:
		var crd apiextensions.CustomResourceDefinition
		o := runtimeClient.ObjectKey{
			Name: v1.SchemeGroupVersion.WithResource(v1.ServiceMonitorName).GroupResource().String(),
		}
		err := r.Client.Get(context.TODO(), o, &crd)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return false, emperror.WrapWith(err, "getting crd failed",
					"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
			}
			//log.Info("crd not found", "resource", "servicemonitors.monitoring.coreos.com")
			return false, nil

		}

	}

	key, err := runtimeClient.ObjectKeyFromObject(current)
	if err != nil {
		return false, emperror.With(err)
	}
	err = r.Client.Get(context.TODO(), key, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return false, emperror.WrapWith(err, "getting resource failed",
				"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
		} else {
			log.Info("resource not found skipping delete", "resource", current.GetObjectKind().GroupVersionKind())
			return false, nil

		}

	}
	log.Info("3")
	err = r.Client.Delete(context.TODO(), current)
	if err != nil {
		return false, emperror.With(err)
	}
	log.Info("resource deleted", "resource", current.GetObjectKind().GroupVersionKind())
	return true, nil
}
