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
	"fmt"
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	"golang.org/x/exp/slices"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources/configcheck"
	loggingv1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

	"github.com/kube-logging/logging-operator/pkg/mirror"
)

const LoggingRefConflict = "Other logging resources exist with the same loggingRef"

func NewValidationReconciler(
	repo client.StatusClient,
	resources LoggingResources,
	secrets SecretLoaderFactory,
	logger logr.Logger,
) func(ctx context.Context) (*reconcile.Result, error) {
	return func(ctx context.Context) (*reconcile.Result, error) {
		// Make sure that you call registerForPatching() before modifying the object
		var patchRequests []patchRequest
		registerForPatching := func(obj client.Object) {
			patchRequests = append(patchRequests, patchRequest{
				Obj:   obj,
				Patch: client.MergeFrom(obj.DeepCopyObject().(client.Object)),
			})
		}

		for i := range resources.Fluentd.ClusterOutputs {
			output := &resources.Fluentd.ClusterOutputs[i]
			registerForPatching(output)

			output.Status.Active = utils.BoolPointer(false)
			output.Status.Problems = nil

			if output.Name == resources.Logging.Spec.ErrorOutputRef {
				output.Status.Active = utils.BoolPointer(true)
			}

			output.Status.Problems = append(output.Status.Problems,
				validateOutputSpec(output.Spec.OutputSpec, secrets.OutputSecretLoaderForNamespace(output.Namespace))...)
			output.Status.ProblemsCount = len(output.Status.Problems)
		}

		for i := range resources.Fluentd.Outputs {
			output := &resources.Fluentd.Outputs[i]
			registerForPatching(output)

			output.Status.Active = utils.BoolPointer(false)
			output.Status.Problems = nil

			output.Status.Problems = append(output.Status.Problems,
				validateOutputSpec(output.Spec, secrets.OutputSecretLoaderForNamespace(output.Namespace))...)
			output.Status.ProblemsCount = len(output.Status.Problems)
		}

		for i := range resources.SyslogNG.ClusterOutputs {
			output := &resources.SyslogNG.ClusterOutputs[i]
			registerForPatching(output)

			output.Status.Active = utils.BoolPointer(false)
			output.Status.Problems = nil

			if output.Name == resources.Logging.Spec.ErrorOutputRef {
				output.Status.Active = utils.BoolPointer(true)
			}

			output.Status.Problems = append(output.Status.Problems,
				validateOutputSpec(output.Spec.SyslogNGOutputSpec, secrets.OutputSecretLoaderForNamespace(output.Namespace))...)
			output.Status.ProblemsCount = len(output.Status.Problems)
		}

		for i := range resources.SyslogNG.Outputs {
			output := &resources.SyslogNG.Outputs[i]
			registerForPatching(output)

			output.Status.Active = utils.BoolPointer(false)
			output.Status.Problems = nil

			output.Status.Problems = append(output.Status.Problems,
				validateOutputSpec(output.Spec, secrets.OutputSecretLoaderForNamespace(output.Namespace))...)
			output.Status.ProblemsCount = len(output.Status.Problems)
		}

		for i := range resources.Fluentd.ClusterFlows {
			flow := &resources.Fluentd.ClusterFlows[i]
			registerForPatching(flow)

			flow.Status.Active = utils.BoolPointer(false)
			flow.Status.Problems = nil

			if len(flow.Spec.GlobalOutputRefs) == 0 && len(flow.Spec.OutputRefs) > 0 {
				flow.Status.Problems = append(flow.Status.Problems, "\"outputRefs\" field is deprecated, use \"globalOutputRefs\" instead")
			}

			for _, ref := range flow.Spec.GlobalOutputRefs {
				if output := resources.Fluentd.ClusterOutputs.FindByName(ref); output != nil {
					flow.Status.Active = utils.BoolPointer(true)
					output.Status.Active = utils.BoolPointer(true)
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling global output reference: %s", ref))
				}
			}
			flow.Status.ProblemsCount = len(flow.Status.Problems)
		}

		for i := range resources.Fluentd.Flows {
			flow := &resources.Fluentd.Flows[i]
			registerForPatching(flow)

			flow.Status.Active = utils.BoolPointer(false)
			flow.Status.Problems = nil

			if len(flow.Spec.LocalOutputRefs)+len(flow.Spec.GlobalOutputRefs) == 0 && len(flow.Spec.OutputRefs) > 0 {
				flow.Status.Problems = append(flow.Status.Problems, "\"outputRefs\" field is deprecated, use \"globalOutputRefs\" and \"localOutputRefs\" instead")
			}

			for _, ref := range flow.Spec.GlobalOutputRefs {
				if output := resources.Fluentd.ClusterOutputs.FindByName(ref); output != nil {
					flow.Status.Active = utils.BoolPointer(true)
					output.Status.Active = utils.BoolPointer(true)
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling global output reference: %s", ref))
				}
			}

			for _, ref := range flow.Spec.LocalOutputRefs {
				if output := resources.Fluentd.Outputs.FindByNamespacedName(flow.Namespace, ref); output != nil {
					flow.Status.Active = utils.BoolPointer(true)
					output.Status.Active = utils.BoolPointer(true)
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling local output reference: %s", ref))
				}
			}
			flow.Status.ProblemsCount = len(flow.Status.Problems)
		}

		for i := range resources.SyslogNG.ClusterFlows {
			flow := &resources.SyslogNG.ClusterFlows[i]
			registerForPatching(flow)

			flow.Status.Active = utils.BoolPointer(false)
			flow.Status.Problems = nil

			for _, ref := range flow.Spec.GlobalOutputRefs {
				if output := resources.SyslogNG.ClusterOutputs.FindByName(ref); output != nil {
					flow.Status.Active = utils.BoolPointer(true)
					output.Status.Active = utils.BoolPointer(true)
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling global output reference: %s", ref))
				}
			}
			flow.Status.ProblemsCount = len(flow.Status.Problems)
		}

		for i := range resources.SyslogNG.Flows {
			flow := &resources.SyslogNG.Flows[i]
			registerForPatching(flow)

			flow.Status.Active = utils.BoolPointer(false)
			flow.Status.Problems = nil

			for _, ref := range flow.Spec.GlobalOutputRefs {
				if output := resources.SyslogNG.ClusterOutputs.FindByName(ref); output != nil {
					flow.Status.Active = utils.BoolPointer(true)
					output.Status.Active = utils.BoolPointer(true)
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling global output reference: %s", ref))
				}
			}

			for _, ref := range flow.Spec.LocalOutputRefs {
				if output := resources.SyslogNG.Outputs.FindByNamespacedName(flow.Namespace, ref); output != nil {
					flow.Status.Active = utils.BoolPointer(true)
					output.Status.Active = utils.BoolPointer(true)
				} else {
					flow.Status.Problems = append(flow.Status.Problems, fmt.Sprintf("dangling local output reference: %s", ref))
				}
			}
			flow.Status.ProblemsCount = len(flow.Status.Problems)
		}

		registerForPatching(&resources.Logging)
		resources.Logging.Status.Problems = nil
		resources.Logging.Status.WatchNamespaces = nil

		if len(resources.Fluentd.ExcessFluentds) != 0 {
			logger.Info("Excess Fluentd CRDs found")
			resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, "multiple fluentd configurations found, couldn't associate it with logging")
			for i := range resources.Fluentd.ExcessFluentds {
				excessFluentd := &resources.Fluentd.ExcessFluentds[i]
				registerForPatching(excessFluentd)
				excessFluentd.Status.Problems = nil
				excessFluentd.Status.Active = utils.BoolPointer(false)
				excessFluentd.Status.Logging = ""

				if len(resources.Logging.Status.FluentdConfigName) == 0 {
					excessFluentd.Status.Problems = append(excessFluentd.Status.Problems, "multiple fluentd configurations found, couldn't associate it with logging")
				} else if resources.Logging.Status.FluentdConfigName != excessFluentd.Name {
					excessFluentd.Status.Problems = append(excessFluentd.Status.Problems, "logging already has a detached fluentd configuration, remove excess configuration objects")
				}
				excessFluentd.Status.ProblemsCount = len(excessFluentd.Status.Problems)
			}
		}
		if resources.Fluentd.Configuration != nil {
			registerForPatching(resources.Fluentd.Configuration)

			if resources.Logging.Spec.FluentdSpec != nil {
				resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, fmt.Sprintf("Fluentd configuration reference set (name=%s), but inline fluentd configuration is set as well, clearing inline", resources.Fluentd.Configuration.Name))
				resources.Logging.Spec.FluentdSpec = nil
			}
			logger.Info("found detached fluentd aggregator, making association", "name", resources.Fluentd.Configuration.Name)
			resources.Logging.Status.FluentdConfigName = resources.Fluentd.Configuration.Name

			resources.Fluentd.Configuration.Status.Active = utils.BoolPointer(true)
			resources.Fluentd.Configuration.Status.Logging = resources.Logging.Name
		} else {
			resources.Logging.Status.FluentdConfigName = ""
		}

		if len(resources.SyslogNG.ExcessSyslogNGs) != 0 {
			logger.Info("Excess SyslogNG CRDs found")
			resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, "multiple syslog-ng configurations found, couldn't associate it with logging")
			for i := range resources.SyslogNG.ExcessSyslogNGs {
				excessSyslogNG := &resources.SyslogNG.ExcessSyslogNGs[i]
				registerForPatching(excessSyslogNG)
				excessSyslogNG.Status.Problems = nil
				excessSyslogNG.Status.Active = utils.BoolPointer(false)
				excessSyslogNG.Status.Logging = ""

				if len(resources.Logging.Status.SyslogNGConfigName) == 0 {
					excessSyslogNG.Status.Problems = append(excessSyslogNG.Status.Problems, "multiple fluentd configurations found, couldn't associate it with logging")
				} else if resources.Logging.Status.FluentdConfigName != excessSyslogNG.Name {
					excessSyslogNG.Status.Problems = append(excessSyslogNG.Status.Problems, "logging already has a detached syslog-ng configuration, remove excess configuration objects")
				}
				excessSyslogNG.Status.ProblemsCount = len(excessSyslogNG.Status.Problems)
			}
		}

		if resources.SyslogNG.Configuration != nil {
			registerForPatching(resources.SyslogNG.Configuration)

			if resources.Logging.Spec.SyslogNGSpec != nil {
				resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, fmt.Sprintf("syslog-ng configuration reference set (name=%s), but inline syslog-ng configuration is set as well, clearing inline", resources.SyslogNG.Configuration.Name))
				resources.Logging.Spec.SyslogNGSpec = nil
			}
			logger.Info("found detached syslog-ng aggregator, making association", "name=", resources.SyslogNG.Configuration.Name)
			resources.Logging.Status.SyslogNGConfigName = resources.SyslogNG.Configuration.Name
			logger.Info("found detached syslog-ng aggregator, making association, done: ", "name=", resources.Logging.Status.SyslogNGConfigName)
			resources.SyslogNG.Configuration.Status.Active = utils.BoolPointer(true)
			resources.SyslogNG.Configuration.Status.Logging = resources.Logging.Name
		} else {
			resources.Logging.Status.SyslogNGConfigName = ""
		}

		if !resources.Logging.WatchAllNamespaces() {
			resources.Logging.Status.WatchNamespaces = resources.WatchNamespaces
		}

		if resources.Logging.Spec.WatchNamespaceSelector != nil &&
			len(resources.Logging.Status.WatchNamespaces) == 0 {
			resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, "Defined watchNamespaceSelector did not match any namespaces")
		}

		loggingsForTheSameRef := make([]string, 0)
		for _, l := range resources.AllLoggings {
			if l.Name == resources.Logging.Name {
				continue
			}
			if l.Spec.LoggingRef == resources.Logging.Spec.LoggingRef && hasIntersection(l.Status.WatchNamespaces, resources.Logging.Status.WatchNamespaces) {
				loggingsForTheSameRef = append(loggingsForTheSameRef, l.Name)
			}
		}

		if len(loggingsForTheSameRef) > 0 {
			problem := fmt.Sprintf("%s (%s) and their watchNamespaces conflict", LoggingRefConflict,
				strings.Join(loggingsForTheSameRef, ","))
			logger.Info(fmt.Sprintf("WARNING %s", problem))
			resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, problem)
		}

		for hash, r := range resources.Logging.Status.ConfigCheckResults {
			if !r {
				problem := fmt.Sprintf("Configuration with checksum %s has failed. "+
					"Config secrets: `kubectl get secret -n %s -l %s=%s`. "+
					"Configcheck pod log: `kubectl logs -n %s -l %s=%s --tail -1`",
					hash,
					resources.Logging.Spec.ControlNamespace, configcheck.HashLabel, hash,
					resources.Logging.Spec.ControlNamespace, configcheck.HashLabel, hash)
				resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, problem)
			}
		}

		if len(resources.Logging.Spec.NodeAgents) > 0 || len(resources.NodeAgents) > 0 {
			// load agents from standalone NodeAgent resources and additionally with inline nodeAgents from the logging resource
			// for compatibility reasons
			agents := make(map[string]loggingv1beta1.NodeAgentConfig)
			for _, a := range resources.NodeAgents {
				agents[a.Name] = a.Spec.NodeAgentConfig
			}
			for _, a := range resources.Logging.Spec.NodeAgents {
				if _, exists := agents[a.Name]; !exists {
					agents[a.Name] = a.NodeAgentConfig
					problem := fmt.Sprintf("inline nodeAgent definition (%s) in Logging resource is deprecated, use standalone NodeAgent CRD instead!", a.Name)
					resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, problem)
				} else {
					problem := fmt.Sprintf("NodeAgent resource overrides inline nodeAgent definition (%s) in Logging resource", a.Name)
					resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, problem)
				}
			}
		}

		if resources.Logging.Spec.FluentbitSpec != nil && len(resources.LoggingRoutes) > 0 {
			resources.Logging.Status.Problems = append(resources.Logging.Status.Problems, "Logging routes are not supported for embedded fluentbit configs, please use a separate FluentbitAgent resource!")
		}

		slices.Sort(resources.Logging.Status.Problems)
		resources.Logging.Status.ProblemsCount = len(resources.Logging.Status.Problems)

		var errs error
		for _, req := range patchRequests {
			if req.IsEmptyPatch() {
				continue
			}

			obj := req.Obj.DeepCopyObject().(client.Object) // copy object so that the original is not changed by the call to Patch
			if err := repo.Status().Patch(ctx, obj, req.Patch); err != nil {
				errs = errors.Append(errs, err)
			}
		}

		return nil, errs
	}
}

func validateOutputSpec(spec interface{}, secrets secret.SecretLoader) (problems []string) {
	var configuredFields []string
	it := mirror.StructRange(spec)
	for it.Next() {
		if it.Field().Type.Kind() == reflect.Ptr && !it.Value().IsNil() {
			configuredFields = append(configuredFields, jsonFieldName(it.Field()))
			problems = append(problems, checkSecrets(it.Value().Elem(), secrets)...)
		}
	}

	switch len(configuredFields) {
	case 0:
		problems = append(problems, "no output target configured")
	case 1:
		// OK
	default:
		problems = append(problems, fmt.Sprintf("multiple output targets configured: %s", configuredFields))
	}
	return
}

func checkSecrets(v reflect.Value, secrets secret.SecretLoader) (problems []string) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			problems = append(problems, checkSecrets(v.Index(i), secrets)...)
		}
	case reflect.Pointer:
		problems = checkSecrets(v.Elem(), secrets)
	case reflect.Struct:
		it := mirror.NewStructIter(v)
		for it.Next() {
			if s, _ := it.Value().Interface().(*secret.Secret); s != nil {
				if _, err := secrets.Load(s); err != nil {
					problems = append(problems, err.Error())
				}
			}
		}
	}
	return
}

type patchRequest struct {
	Obj   client.Object
	Patch client.Patch
}

func (r patchRequest) IsEmptyPatch() bool {
	data, err := r.Patch.Data(r.Obj)
	return err == nil && string(data) == "{}"
}

func jsonFieldName(f reflect.StructField) string {
	t := f.Tag.Get("json")
	n := strings.Split(t, ",")[0]
	if n != "" {
		return n
	}
	return f.Name
}

func hasIntersection(a, b []string) bool {
	for _, i := range a {
		for _, j := range b {
			if i == j {
				return true
			}
		}
	}
	return false
}
