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
	"fmt"
	"reflect"
	"time"

	"emperror.dev/errors"
	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	StateAbsent  StaticDesiredState = "Absent"
	StatePresent StaticDesiredState = "Present"
)

type DesiredState interface {
	BeforeUpdate(object runtime.Object) error
}

type StaticDesiredState string

func (s StaticDesiredState) BeforeUpdate(object runtime.Object) error {
	return nil
}

type DesiredStateHook func(object runtime.Object) error

func (d DesiredStateHook) BeforeUpdate(object runtime.Object) error {
	return d(object)
}

// GenericResourceReconciler generic resource reconciler
type GenericResourceReconciler struct {
	Log     logr.Logger
	Client  runtimeClient.Client
	Options ReconcilerOpts
}

type ReconcilerOpts struct {
	EnableRecreateWorkloadOnImmutableFieldChange     bool
	EnableRecreateWorkloadOnImmutableFieldChangeHelp string
}

// NewReconciler returns GenericResourceReconciler
func NewReconciler(client runtimeClient.Client, log logr.Logger, opts ReconcilerOpts) *GenericResourceReconciler {
	return &GenericResourceReconciler{
		Log:     log,
		Client:  client,
		Options: opts,
	}
}

// CreateResource creates a resource if it doesn't exist
func (r *GenericResourceReconciler) CreateResource(desired runtime.Object) error {
	_, _, err := r.createIfNotExists(desired)
	return err
}

// ReconcileResource reconciles various kubernetes types
func (r *GenericResourceReconciler) ReconcileResource(desired runtime.Object, desiredState DesiredState) (*reconcile.Result, error) {
	log := r.Log.WithValues("type", reflect.TypeOf(desired))

	switch desiredState {
	default:
		created, current, err := r.createIfNotExists(desired)
		if err == nil && created {
			return nil, nil
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create resource %+v", desired)
		}

		// last chance to hook into the desired state armed with the knowledge of the current state
		err = desiredState.BeforeUpdate(current)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get desired state dynamically")
		}
		key, err := runtimeClient.ObjectKeyFromObject(current)
		if err != nil {
			return nil, errors.Wrapf(err, "meta accessor failed %+v", current)
		}
		if err == nil {
			if metaObject, ok := current.(metav1.Object); ok {
				if metaObject.GetDeletionTimestamp() != nil {
					r.Log.Info(fmt.Sprintf("object %s is being deleted, backing off", metaObject.GetSelfLink()))
					return &reconcile.Result{RequeueAfter: time.Second * 2}, nil
				}
			}
			patchResult, err := patch.DefaultPatchMaker.Calculate(current, desired)
			if err != nil {
				log.Error(err, "could not match objects",
					"kind", desired.GetObjectKind().GroupVersionKind(), "name", key.Name)
			} else if patchResult.IsEmpty() {
				log.V(1).Info("resource is in sync",
					"kind", desired.GetObjectKind().GroupVersionKind(), "name", key.Name)
				return nil, nil
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
				return nil, errors.Wrap(err, "failed to access resourceVersion from metadata")
			}
			metaAccessor.SetResourceVersion(desired, currentResourceVersion)

			var name string
			if name, err = metaAccessor.Name(current); err != nil {
				return nil, errors.Wrap(err, "failed to access Name from metadata")
			}

			log.V(1).Info("Updating resource",
				"gvk", desired.GetObjectKind().GroupVersionKind(), "name", name)
			if err := r.Client.Update(context.TODO(), desired); err != nil {
				sErr, ok := err.(*apierrors.StatusError)
				if ok && sErr.ErrStatus.Code == 422 && sErr.ErrStatus.Reason == metav1.StatusReasonInvalid {
					if r.Options.EnableRecreateWorkloadOnImmutableFieldChange {
						r.Log.Error(err, "failed to update resource, trying to recreate")
						err := r.Client.Delete(context.TODO(), current,
							// wait until all dependent resources gets cleared up
							runtimeClient.PropagationPolicy(metav1.DeletePropagationForeground),
						)
						if err != nil {
							return nil, errors.Wrapf(err, "failed to delete resource %+v", current)
						}
						return &reconcile.Result{
							Requeue:      true,
							RequeueAfter: time.Second * 10,
						}, nil
					} else {
						return nil, errors.New(r.Options.EnableRecreateWorkloadOnImmutableFieldChangeHelp)
					}
				}
				return nil, emperror.WrapWith(err, "updating resource failed",
					"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
			}
			log.Info("resource updated", "resource", desired.GetObjectKind().GroupVersionKind())
		}
	case StateAbsent:
		_, err := r.delete(desired)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to delete resource %+v", desired)
		}
	}
	return nil, nil
}

func (r *GenericResourceReconciler) createIfNotExists(desired runtime.Object) (bool, runtime.Object, error) {
	log := r.Log.WithValues("type", reflect.TypeOf(desired))
	var current = desired.DeepCopyObject()
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
	key, err := runtimeClient.ObjectKeyFromObject(current)
	if err != nil {
		return false, emperror.With(err)
	}
	err = r.Client.Get(context.TODO(), key, current)
	if err != nil {
		// If the resource type does not exist we should be ok to move on
		if meta.IsNoMatchError(err) {
			return false, nil
		}
		if !apierrors.IsNotFound(err) {
			return false, emperror.WrapWith(err, "getting resource failed",
				"resource", desired.GetObjectKind().GroupVersionKind(), "type", reflect.TypeOf(desired))
		} else {
			log.V(1).Info("resource not found skipping delete", "resource", current.GetObjectKind().GroupVersionKind())
			return false, nil
		}
	}
	err = r.Client.Delete(context.TODO(), current)
	if err != nil {
		return false, emperror.With(err)
	}
	log.Info("resource deleted", "resource", current.GetObjectKind().GroupVersionKind())
	return true, nil
}
