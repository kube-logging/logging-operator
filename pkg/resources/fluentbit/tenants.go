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

package fluentbit

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"emperror.dev/errors"
	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/resources/model"
	"github.com/kube-logging/logging-operator/pkg/resources/syslogng"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type Tenant struct {
	Name         string
	AllNamespace bool
	Namespaces   []string
}

func FindTenants(ctx context.Context, target metav1.LabelSelector, reader client.Reader) ([]Tenant, error) {
	var tenants []Tenant

	selector, err := metav1.LabelSelectorAsSelector(&target)
	if err != nil {
		return nil, errors.WrapIf(err, "logrouting targetSelector")
	}
	listOptions := &client.ListOptions{
		LabelSelector: selector,
	}
	loggingList := &v1beta1.LoggingList{}
	if err := reader.List(ctx, loggingList, listOptions); err != nil {
		return nil, errors.WrapIf(err, "listing loggings for targetSelector")
	}
	for _, l := range loggingList.Items {
		l := l
		targetNamespaces, err := model.UniqueWatchNamespaces(ctx, reader, &l)
		if err != nil {
			return nil, err
		}
		if l.WatchAllNamespaces() {
			tenants = append(tenants, Tenant{
				Name:         l.Name,
				AllNamespace: true,
			})
		} else {
			tenants = append(tenants, Tenant{
				Name:       l.Name,
				Namespaces: targetNamespaces,
			})
		}
	}

	sort.Slice(tenants, func(i, j int) bool {
		return tenants[i].Name < tenants[j].Name
	})
	// Make sure our tenant list is stable
	slices.SortStableFunc(tenants, func(a, b Tenant) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name == b.Name {
			return 0
		}
		return 1
	})

	return tenants, nil
}

func (r *Reconciler) configureOutputsForTenants(ctx context.Context, tenants []v1beta1.Tenant, input *fluentBitConfig) error {
	var errs error
	for _, t := range tenants {
		allNamespaces := len(t.Namespaces) == 0
		namespaceRegex := `.`
		if !allNamespaces {
			namespaceRegex = fmt.Sprintf("^[^_]+_(%s)_", strings.Join(t.Namespaces, "|"))
		}
		logging := &v1beta1.Logging{}
		if err := r.resourceReconciler.Client.Get(ctx, types.NamespacedName{Name: t.Name}, logging); err != nil {
			return errors.WrapIf(err, "getting logging resource")
		}
		fluentdSpec := logging.Spec.FluentdSpec
		if detachedFluentd := fluentd.GetFluentd(ctx, r.resourceReconciler.Client, r.logger, logging.Spec.ControlNamespace); detachedFluentd != nil {
			fluentdSpec = &detachedFluentd.Spec
		}
		if fluentdSpec != nil {
			if input.FluentForwardOutput == nil {
				input.FluentForwardOutput = &fluentForwardOutputConfig{}
			}
			input.FluentForwardOutput.Targets = append(input.FluentForwardOutput.Targets, forwardTargetConfig{
				AllNamespaces:  allNamespaces,
				NamespaceRegex: namespaceRegex,
				Host:           aggregatorEndpoint(logging, fluentd.ServiceName),
				Port:           fluentd.ServicePort,
			})
		} else if logging.Spec.SyslogNGSpec != nil {
			if input.SyslogNGOutput == nil {
				input.SyslogNGOutput = newSyslogNGOutputConfig()
			}
			input.SyslogNGOutput.Targets = append(input.SyslogNGOutput.Targets, forwardTargetConfig{
				AllNamespaces:  allNamespaces,
				NamespaceRegex: namespaceRegex,
				Host:           aggregatorEndpoint(logging, syslogng.ServiceName),
				Port:           syslogng.ServicePort,
			})
		} else {
			errs = errors.Append(errs, errors.Errorf("logging %s does not provide any aggregator configured", t.Name))
		}
	}
	return errs
}
