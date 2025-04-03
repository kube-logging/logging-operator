// Copyright © 2024 Kube logging authors
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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	telemetry_controller "github.com/kube-logging/logging-operator/pkg/resources/telemetry-controller"
	loggingv1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	TelemetryControllerFinalizer = "telemetrycontroller.logging.banzaicloud.io/finalizer"
)

// +kubebuilder:rbac:groups=telemetry.kube-logging.dev,resources=collectors;tenants;subscriptions;outputs;bridges;,verbs=get;list;watch;create;update;patch;delete

func NewTelemetryControllerReconciler(client client.Client, log logr.Logger) *TelemetryControllerReconciler {
	return &TelemetryControllerReconciler{
		Client: client,
		Log:    log,
	}
}

// TelemetryControllerReconciler reconciles Logging resources for the Telemetry controller
type TelemetryControllerReconciler struct {
	client.Client
	Log logr.Logger
}

func (r *TelemetryControllerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("telemetry-controller", req.Name)

	var logging loggingv1beta1.Logging
	if err := r.Get(ctx, req.NamespacedName, &logging); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if logging.Spec.RouteConfig.EnableTelemetryControllerRoute {
		log.Info("Reconciling Logging resource for Telemetry controller", "name", logging.Name)

		objectsToCreate := r.createTelemetryControllerResources(log, &logging)

		if err := r.finalizeLoggingForTelemetryController(ctx, log, &logging, &objectsToCreate); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.isAggregatorReady(ctx, log, logging); err != nil {
			r.Log.Info(fmt.Sprintf("Aggregator pod is not ready yet: %s", err))
			return ctrl.Result{RequeueAfter: 5}, nil
		}

		if err := r.deployTelemetryControllerResources(ctx, log, &objectsToCreate); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func SetupTelemetryControllerWithManager(mgr ctrl.Manager, logger logr.Logger) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Logging{}).
		Named("telemetrycontroller").
		Complete(NewTelemetryControllerReconciler(mgr.GetClient(), logger))
}

func (r *TelemetryControllerReconciler) createTelemetryControllerResources(logger logr.Logger, logging *loggingv1beta1.Logging) []client.Object {
	logger.Info("Creating Telemetry controller resources")

	objectsToCreate := []client.Object{}
	objectsToCreate = append(objectsToCreate, telemetry_controller.CreateTenant(logging))
	objectsToCreate = append(objectsToCreate, telemetry_controller.CreateSubscription(logging))
	objectsToCreate = append(objectsToCreate, telemetry_controller.CreateOutput(logging))

	return objectsToCreate
}

func (r *TelemetryControllerReconciler) finalizeLoggingForTelemetryController(ctx context.Context, logger logr.Logger, logging *loggingv1beta1.Logging, objectsToCreate *[]client.Object) error {
	logger.Info("Finalizing Telemetry controller resources")

	if logging.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(logging, TelemetryControllerFinalizer) {
			r.Log.Info("adding telemetrycontroller finalizer")
			controllerutil.AddFinalizer(logging, TelemetryControllerFinalizer)
			if err := r.Update(ctx, logging); err != nil {
				return err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(logging, TelemetryControllerFinalizer) {
			if err := r.deleteTelemetryControllerResources(ctx, logger, objectsToCreate); err != nil {
				return err
			}

			r.Log.Info("removing telemetrycontroller finalizer")
			controllerutil.RemoveFinalizer(logging, TelemetryControllerFinalizer)
			if err := r.Update(ctx, logging); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *TelemetryControllerReconciler) deployTelemetryControllerResources(ctx context.Context, logger logr.Logger, objectsToCreate *[]client.Object) error {
	logger.Info("Deploying Telemetry controller resources")

	for _, objectToCreate := range *objectsToCreate {
		if err := r.Get(ctx, client.ObjectKeyFromObject(objectToCreate), objectToCreate); err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			if err := r.Create(ctx, objectToCreate); err != nil {
				return err
			}
			logger.Info("Created object", "object", objectToCreate.GetName())
		} else {
			logger.Info("Object already exists", "object", objectToCreate.GetName())
		}
	}

	return nil
}

func (r *TelemetryControllerReconciler) deleteTelemetryControllerResources(ctx context.Context, logger logr.Logger, objectsToCreate *[]client.Object) error {
	logger.Info("Logging resource is being deleted, deleting Telemetry controller resources")

	for _, obj := range *objectsToCreate {
		if err := r.Delete(ctx, obj); err != nil {
			return client.IgnoreNotFound(err)
		}
		logger.Info("Deleted object", "object", obj.GetName())
	}

	return nil
}

func (r *TelemetryControllerReconciler) isAggregatorReady(ctx context.Context, logger logr.Logger, logging loggingv1beta1.Logging) error {
	logger.Info("Waiting for aggregator pod to be ready")

	podName := fmt.Sprintf("%s-fluentd-0", logging.Name)
	pod := &corev1.Pod{}
	err := r.Get(ctx, client.ObjectKey{Name: podName, Namespace: logging.Spec.ControlNamespace}, pod)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("aggregator pod: %s not found", podName)
		}
		return err
	}

	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			r.Log.Info("Aggregator pod is ready", "pod", pod.Name)
			return nil
		}
	}

	return fmt.Errorf("aggregator pod: %s is not ready", podName)
}
