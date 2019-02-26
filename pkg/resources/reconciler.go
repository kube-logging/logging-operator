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

package resources

import (
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginReconciler struct {
	Client client.Client
	Plugin *loggingv1alpha1.LoggingPlugin
}

type FluentdReconciler struct {
	Client  client.Client
	Fluentd *loggingv1alpha1.Fluentd
}

type FluentbitReconciler struct {
	Client    client.Client
	Fluentbit *loggingv1alpha1.Fluentbit
}

type ComponentReconciler interface {
	Reconcile(log logr.Logger) error
}

type Resource func() runtime.Object

type ResourceVariation func(t string) runtime.Object
