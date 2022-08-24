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
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/banzaicloud/logging-operator/pkg/mirror"
)

func NewValidationReconciler(
	ctx context.Context,
	repo client.StatusClient,
	resources LoggingResources,
	secrets SecretLoaderFactory,
) func() (*reconcile.Result, error) {
	return func() (*reconcile.Result, error) {
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

		var errs error
		for _, req := range patchRequests {
			if req.IsEmptyPatch() {
				continue
			}

			if err := repo.Status().Patch(ctx, req.Obj, req.Patch); err != nil {
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
