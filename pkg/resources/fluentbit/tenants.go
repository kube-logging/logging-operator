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
	"strings"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/resources/model"
	"github.com/kube-logging/logging-operator/pkg/resources/syslogng"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type Tenant struct {
	Logging      *v1beta1.Logging
	AllNamespace bool
	Namespaces   []string
}

func FindTenants(ctx context.Context, targets []v1beta1.Target, currentLogging string, reader client.Reader, logger logr.Logger) ([]Tenant, error) {
	var tenantCandidates []Tenant

	for _, t := range targets {
		if t.LoggingName != "" {
			l := &v1beta1.Logging{}
			if err := reader.Get(ctx, client.ObjectKey{Name: t.LoggingName}, l); err != nil {
				return nil, errors.WrapIff(err, "logrouting target %", t.LoggingName)
			}
			tenantCandidates = append(tenantCandidates, Tenant{
				Logging: l,
			})
		}
		if t.LoggingSelector != nil {
			selector, err := metav1.LabelSelectorAsSelector(t.LoggingSelector)
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
				tenantCandidates = append(tenantCandidates, Tenant{
					Logging: &l,
				})
			}
		}
	}

	var validTenants []Tenant
	for _, t := range tenantCandidates {
		targetNamespaces, allNamespaces, err := model.UniqueWatchNamespaces(ctx, reader, t.Logging)
		if err != nil {
			logger.Error(err, "getting target namespaces for logging %s", t.Logging.Name)
			continue
		}
		if len(targetNamespaces) == 0 {
			logger.Info(fmt.Sprintf("WARNING unable to use logging %s as a valid target as no watch namespaces have been found", t.Logging.Name))
			continue
		}
		if allNamespaces && t.Logging.Name != currentLogging {
			logger.Info(fmt.Sprintf("WARNING refusing to send logs from all namespaces to logging %s", t.Logging.Name))
			continue
		}
		validTenants = append(validTenants, Tenant{
			Logging:      t.Logging,
			AllNamespace: allNamespaces,
			Namespaces:   targetNamespaces,
		})
	}
	return validTenants, nil
}

func (r *Reconciler) configureOutputsForTenants(tenants []Tenant, input *fluentBitConfig) error {
	var errs error
	for _, t := range tenants {
		namespaceRegex := fmt.Sprintf("^[^_]+_(%s)_", strings.Join(t.Namespaces, "|"))
		if t.Logging.Spec.FluentdSpec != nil {
			if input.FluentForwardOutput == nil {
				input.FluentForwardOutput = &fluentForwardOutputConfig{}
			}
			input.FluentForwardOutput.Targets = append(input.FluentForwardOutput.Targets, forwardTargetConfig{
				AllNamespaces:  t.AllNamespace,
				NamespaceRegex: namespaceRegex,
				Host:           aggregatorEndpoint(t.Logging, fluentd.ServiceName),
				Port:           fluentd.ServicePort,
			})
		} else if t.Logging.Spec.SyslogNGSpec != nil {
			if input.SyslogNGOutput == nil {
				input.SyslogNGOutput = newSyslogNGOutputConfig()
			}
			input.SyslogNGOutput.Targets = append(input.SyslogNGOutput.Targets, forwardTargetConfig{
				AllNamespaces:  t.AllNamespace,
				NamespaceRegex: namespaceRegex,
				Host:           aggregatorEndpoint(t.Logging, syslogng.ServiceName),
				Port:           syslogng.ServicePort,
			})
		} else {
			errs = errors.Append(errs, errors.Errorf("logging %s does not provide any aggregator configured", t.Logging.Name))
		}
	}
	return errs
}
