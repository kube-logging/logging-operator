/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fluentbit

import (
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	serviceAccountName      = "logging"
	clusterRoleBindingName  = "logging"
	clusterRoleName         = "LoggingRole"
	fluentbitConfigMapName  = "fluent-bit-config"
	fluentbitDeaemonSetName = "fluent-bit-daemon"
)

var labelSelector = map[string]string{
	"app": "fluent-bit",
}

// Reconciler holds info what resource to reconcile
type Reconciler struct {
	resources.FluentbitReconciler
}

// New creates a new Fluentbit reconciler
func New(client client.Client, fluentbit *loggingv1alpha1.Fluentbit) *Reconciler {
	return &Reconciler{
		FluentbitReconciler: resources.FluentbitReconciler{
			Client:    client,
			Fluentbit: fluentbit,
		},
	}
}

// Reconcile reconciles the fluentBit resource
func (r *Reconciler) Reconcile(log logr.Logger) error {
	for _, res := range []resources.Resource{
		r.serviceAccount,
		r.clusterRole,
		r.clusterRoleBinding,
		r.configMap, r.daemonSet,
	} {
		o := res()
		err := k8sutil.Reconcile(log, r.Client, o)
		if err != nil {
			return emperror.WrapWith(err, "failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
	}
	return nil
}
