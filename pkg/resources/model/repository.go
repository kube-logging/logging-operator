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
	"sort"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewLoggingResourceRepository(client client.Reader) *LoggingResourceRepository {
	return &LoggingResourceRepository{
		Client: client,
	}
}

type LoggingResourceRepository struct {
	Client client.Reader
}

func (r LoggingResourceRepository) LoggingResourcesFor(ctx context.Context, logging v1beta1.Logging) (res LoggingResources, errs error) {
	res.Logging = logging

	var err error

	res.ClusterFlows, err = r.ClusterFlowsFor(ctx, logging)
	errs = errors.Append(errs, err)

	res.ClusterOutputs, err = r.ClusterOutputsFor(ctx, logging)
	errs = errors.Append(errs, err)

	watchNamespaces := logging.Spec.WatchNamespaces
	if len(watchNamespaces) == 0 {
		var nsList corev1.NamespaceList
		if err := r.Client.List(ctx, &nsList); err != nil {
			errs = errors.Append(errs, errors.WrapIf(err, "listing namespaces"))
			return
		}

		for _, i := range nsList.Items {
			watchNamespaces = append(watchNamespaces, i.Name)
		}
	}
	sort.Strings(watchNamespaces)

	for _, ns := range watchNamespaces {
		flows, err := r.FlowsInNamespaceFor(ctx, ns, logging)
		res.Flows = append(res.Flows, flows...)
		errs = errors.Append(errs, err)

		outputs, err := r.OutputsInNamespaceFor(ctx, ns, logging)
		res.Outputs = append(res.Outputs, outputs...)
		errs = errors.Append(errs, err)
	}

	return
}

func (r LoggingResourceRepository) ClusterFlowsFor(ctx context.Context, logging v1beta1.Logging) ([]v1beta1.ClusterFlow, error) {
	var list v1beta1.ClusterFlowList
	if err := r.Client.List(ctx, &list, clusterResourceListOpts(logging)...); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.ClusterFlow
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) ClusterOutputsFor(ctx context.Context, logging v1beta1.Logging) ([]v1beta1.ClusterOutput, error) {
	var list v1beta1.ClusterOutputList
	if err := r.Client.List(ctx, &list, clusterResourceListOpts(logging)...); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.ClusterOutput
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) FlowsInNamespaceFor(ctx context.Context, namespace string, logging v1beta1.Logging) ([]v1beta1.Flow, error) {
	var list v1beta1.FlowList
	if err := r.Client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.Flow
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) OutputsInNamespaceFor(ctx context.Context, namespace string, logging v1beta1.Logging) ([]v1beta1.Output, error) {
	var list v1beta1.OutputList
	if err := r.Client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.Output
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) SyslogNGLoggingResourcesFor(ctx context.Context, logging v1beta1.Logging) (res SyslogNGLoggingResources, errs error) {
	res.Logging = logging

	var err error

	res.ClusterFlows, err = r.SyslogNGClusterFlowsFor(ctx, logging)
	errs = errors.Append(errs, err)

	res.ClusterOutputs, err = r.SyslogNGClusterOutputsFor(ctx, logging)
	errs = errors.Append(errs, err)

	watchNamespaces := logging.Spec.WatchNamespaces
	if len(watchNamespaces) == 0 {
		var nsList corev1.NamespaceList
		if err := r.Client.List(ctx, &nsList); err != nil {
			errs = errors.Append(errs, errors.WrapIf(err, "listing namespaces"))
			return
		}

		for _, i := range nsList.Items {
			watchNamespaces = append(watchNamespaces, i.Name)
		}
	}
	sort.Strings(watchNamespaces)

	for _, ns := range watchNamespaces {
		flows, err := r.SyslogNGFlowsInNamespaceFor(ctx, ns, logging)
		res.Flows = append(res.Flows, flows...)
		errs = errors.Append(errs, err)

		outputs, err := r.SyslogNGOutputsInNamespaceFor(ctx, ns, logging)
		res.Outputs = append(res.Outputs, outputs...)
		errs = errors.Append(errs, err)
	}

	return
}

func (r LoggingResourceRepository) SyslogNGClusterFlowsFor(ctx context.Context, logging v1beta1.Logging) ([]v1beta1.SyslogNGClusterFlow, error) {
	var list v1beta1.SyslogNGClusterFlowList
	if err := r.Client.List(ctx, &list, clusterResourceListOpts(logging)...); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.SyslogNGClusterFlow
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) SyslogNGClusterOutputsFor(ctx context.Context, logging v1beta1.Logging) ([]v1beta1.SyslogNGClusterOutput, error) {
	var list v1beta1.SyslogNGClusterOutputList
	if err := r.Client.List(ctx, &list, clusterResourceListOpts(logging)...); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.SyslogNGClusterOutput
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) SyslogNGFlowsInNamespaceFor(ctx context.Context, namespace string, logging v1beta1.Logging) ([]v1beta1.SyslogNGFlow, error) {
	var list v1beta1.SyslogNGFlowList
	if err := r.Client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.SyslogNGFlow
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func (r LoggingResourceRepository) SyslogNGOutputsInNamespaceFor(ctx context.Context, namespace string, logging v1beta1.Logging) ([]v1beta1.SyslogNGOutput, error) {
	var list v1beta1.SyslogNGOutputList
	if err := r.Client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return lessByNamespacedName(&list.Items[i], &list.Items[j])
	})

	var res []v1beta1.SyslogNGOutput
	for _, i := range list.Items {
		if i.Spec.LoggingRef == logging.Spec.LoggingRef {
			res = append(res, i)
		}
	}
	return res, nil
}

func clusterResourceListOpts(logging v1beta1.Logging) []client.ListOption {
	var opts []client.ListOption
	if !logging.Spec.AllowClusterResourcesFromAllNamespaces {
		opts = append(opts, client.InNamespace(logging.Spec.ControlNamespace))
	}
	return opts
}

func lessByNamespacedName(a, b interface {
	GetNamespace() string
	GetName() string
}) bool {
	if a.GetNamespace() != b.GetNamespace() {
		return a.GetNamespace() < b.GetNamespace()
	}
	return a.GetName() < b.GetName()
}
