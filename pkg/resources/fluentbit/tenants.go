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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"emperror.dev/errors"
	"golang.org/x/exp/maps"
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

	sort.SliceStable(tenants, func(i, j int) bool {
		return tenants[i].Name < tenants[j].Name
	})

	return tenants, nil
}

func (r *Reconciler) configureOutputsForTenants(ctx context.Context, tenants []v1beta1.Tenant, input *fluentBitConfig) error {
	var errs error
	for _, t := range tenants {
		match := fmt.Sprintf("kubernetes.%s.*", hashFromTenantName(t.Name))
		logging := &v1beta1.Logging{}
		if err := r.resourceReconciler.Client.Get(ctx, types.NamespacedName{Name: t.Name}, logging); err != nil {
			return errors.WrapIf(err, "getting logging resource")
		}

		loggingResources, err := r.loggingResourcesRepo.LoggingResourcesFor(ctx, *logging)
		if err != nil {
			errs = errors.Append(errs, errors.WrapIff(err, "querying related resources for logging %s", logging.Name))
			continue
		}

		_, fluentdSpec := loggingResources.GetFluentd()
		if fluentdSpec != nil {
			if input.FluentForwardOutput == nil {
				input.FluentForwardOutput = &fluentForwardOutputConfig{}
			}
			input.FluentForwardOutput.Targets = append(input.FluentForwardOutput.Targets, forwardTargetConfig{
				Match: match,
				Host:  aggregatorEndpoint(logging, fluentd.ServiceName),
				Port:  fluentd.ServicePort,
			})
		} else if _, syslogNGSPec := loggingResources.GetSyslogNGSpec(); syslogNGSPec != nil {
			if input.SyslogNGOutput == nil {
				input.SyslogNGOutput = newSyslogNGOutputConfig()
			}
			input.SyslogNGOutput.Targets = append(input.SyslogNGOutput.Targets, forwardTargetConfig{
				Match: match,
				Host:  aggregatorEndpoint(logging, syslogng.ServiceName),
				Port:  syslogng.ServicePort,
			})
		} else {
			errs = errors.Append(errs, errors.Errorf("logging %s does not provide any aggregator configured", t.Name))
		}
	}
	return errs
}

func (r *Reconciler) configureInputsForTenants(tenants []v1beta1.Tenant, input *fluentBitConfig) error {
	var errs error
	for _, t := range tenants {
		allNamespaces := len(t.Namespaces) == 0
		tenantValues := maps.Clone(input.Input.Values)
		if !allNamespaces {
			var paths []string
			for _, n := range t.Namespaces {
				paths = append(paths, fmt.Sprintf("/var/log/containers/*_%s_*.log", n))
			}
			tenantValues["Path"] = strings.Join(paths, ",")
		} else {
			tenantValues["Path"] = "/var/log/containers/*.log"
		}

		tenantValues["DB"] = fmt.Sprintf("/tail-db/tail-containers-state-%s.db", t.Name)
		tenantValues["Tag"] = fmt.Sprintf("kubernetes.%s.*", hashFromTenantName(t.Name))
		// This helps to make sure we apply backpressure on the input, see https://docs.fluentbit.io/manual/administration/backpressure
		tenantValues["storage.pause_on_chunks_overlimit"] = "on"
		input.Inputs = append(input.Inputs, fluentbitInputConfigWithTenant{
			Tenant:          t.Name,
			Values:          tenantValues,
			ParserN:         input.Input.ParserN,
			MultilineParser: input.Input.MultilineParser,
		})
	}
	// the regex will work only if we cut the prefix off. fluentbit doesn't care about the content, just the length
	input.KubernetesFilter["Kube_Tag_Prefix"] = `kubernetes.0000000000.var.log.containers.`
	return errs
}

func hashFromTenantName(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString[0:10]
}
