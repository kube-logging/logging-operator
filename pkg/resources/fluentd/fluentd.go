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

package fluentd

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/util"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	SecretConfigName      = "fluentd"
	AppSecretConfigName   = "fluentd-app"
	AppConfigKey          = "fluentd.conf"
	StatefulSetName       = "fluentd"
	PodSecurityPolicyName = "fluentd"
	ServiceName           = "fluentd"
	OutputSecretName      = "fluentd-output"
	OutputSecretPath      = "/fluentd/secret"

	bufferVolumeName          = "fluentd-buffer"
	defaultServiceAccountName = "fluentd"
	roleBindingName           = "fluentd"
	roleName                  = "fluentd"
)

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging *v1beta1.Logging
	*k8sutil.GenericResourceReconciler
	config  *string
	secrets *secret.MountSecrets
}

func (r *Reconciler) getFluentdLabels() map[string]string {
	return util.MergeLabels(r.Logging.Labels, map[string]string{
		"app.kubernetes.io/name": "fluentd"}, generataLoggingRefLabels(r.Logging.ObjectMeta.GetName()))
}

func (r *Reconciler) getServiceAccount() string {
	if r.Logging.Spec.FluentdSpec.Security.ServiceAccount != "" {
		return r.Logging.Spec.FluentdSpec.Security.ServiceAccount
	}
	return r.Logging.QualifiedName(defaultServiceAccountName)
}

func New(client client.Client, log logr.Logger, logging *v1beta1.Logging, config *string, secrets *secret.MountSecrets) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		GenericResourceReconciler: k8sutil.NewReconciler(client, log),
		config:                    config,
		secrets:                   secrets,
	}
}

// Reconcile reconciles the fluentd resource
func (r *Reconciler) Reconcile() (*reconcile.Result, error) {
	// Config check and cleanup if enabled
	if !r.Logging.Spec.FlowConfigCheckDisabled {
		hash, err := r.configHash()
		if err != nil {
			return nil, err
		}
		if result, ok := r.Logging.Status.ConfigCheckResults[hash]; ok {
			// We already have an existing configcheck result:
			// - bail out if it was unsuccessful
			// - cleanup previous results if it's successful
			if !result {
				return nil, errors.Errorf("current config is invalid")
			}
			var removedHashes []string
			if removedHashes, err = r.configCheckCleanup(hash); err != nil {
				r.Log.Error(err, "failed to cleanup resources")
			}
			if len(removedHashes) > 0 {
				for _, removedHash := range removedHashes {
					delete(r.Logging.Status.ConfigCheckResults, removedHash)
				}
				if err := r.Client.Status().Update(context.TODO(), r.Logging); err != nil {
					return nil, errors.WrapWithDetails(err, "failed to update status", "logging", r.Logging)
				} else {
					// explicitly ask for a requeue to short circuit the controller loop after the status update
					return &reconcile.Result{Requeue: true}, nil
				}
			}
		} else {
			// We don't have an existing result
			// - let's create what's necessary to have one
			// - if the result is ready write it into the status
			result, err := r.configCheck()
			if err != nil {
				return nil, errors.WrapIf(err, "failed to validate config")
			}
			if result.Ready {
				r.Logging.Status.ConfigCheckResults[hash] = result.Valid
				if err := r.Client.Status().Update(context.TODO(), r.Logging); err != nil {
					return nil, errors.WrapWithDetails(err, "failed to update status", "logging", r.Logging)
				} else {
					// explicitly ask for a requeue to short circuit the controller loop after the status update
					return &reconcile.Result{Requeue: true}, nil
				}
			} else {
				r.Log.Info("still waiting for the configcheck result...")
				return &reconcile.Result{RequeueAfter: time.Second}, nil
			}
		}
	}
	// Prepare output secret
	err := r.ReconcileResource(r.outputSecret(r.secrets, OutputSecretPath))
	if err != nil {
		return nil, errors.WrapIf(err, "failed to reconcile resource")
	}
	// Mark watched secrets
	secretList, state := r.markSecrets(r.secrets)
	for _, obj := range secretList {
		err := r.ReconcileResource(obj, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
	}
	for _, res := range []resources.Resource{
		r.serviceAccount,
		r.role,
		r.roleBinding,
		r.clusterPodSecurityPolicy,
		r.pspRole,
		r.pspRoleBinding,
		r.secretConfig,
		r.appconfigMap,
		r.statefulset,
		r.service,
		r.serviceMetrics,
		r.monitorServiceMetrics,
	} {
		o, state := res()
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %s returns with nil object", res)
		}
		err := r.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
	}

	return nil, nil
}
