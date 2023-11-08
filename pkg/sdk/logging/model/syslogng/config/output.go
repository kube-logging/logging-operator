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

package config

import (
	"fmt"
	"reflect"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/siliconbrain/go-seqs/seqs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func renderClusterOutput(o v1beta1.SyslogNGClusterOutput, secretLoaderFactory SecretLoaderFactory) render.Renderer {
	name := clusterOutputDestName(o.Namespace, o.Name)
	return destinationDefStmt(
		name,
		renderOutputSpec(o.Spec.SyslogNGOutputSpec, name, &o, secretLoaderFactory.SecretLoaderForNamespace(o.Namespace)),
	)
}

func clusterOutputDestName(ns string, name string) string {
	return fmt.Sprintf("clusteroutput_%s_%s", ns, name)
}

func renderOutput(o v1beta1.SyslogNGOutput, secretLoaderFactory SecretLoaderFactory) render.Renderer {
	name := outputDestName(o.Namespace, o.Name)
	return destinationDefStmt(
		name,
		renderOutputSpec(o.Spec, name, &o, secretLoaderFactory.SecretLoaderForNamespace(o.Namespace)),
	)
}

func outputDestName(ns string, name string) string {
	return fmt.Sprintf("output_%s_%s", ns, name)
}

func renderOutputSpec(spec v1beta1.SyslogNGOutputSpec, destName string, output metav1.Object, secretLoader secret.SecretLoader) render.Renderer {
	specValue := reflect.ValueOf(spec)
	driverFields := seqs.ToSlice(seqs.Filter(seqs.FromSlice(fieldsOf(specValue)), isActiveDestinationDriver))
	switch len(driverFields) {
	case 0:
		return render.Error(fmt.Errorf("no destination driver specified on output %s/%s", output.GetNamespace(), output.GetName()))
	case 1:
		driverField := driverFields[0]
		defaultPersistName(driverField.Value, destName) // HACK: defaulting should be done properly
		if br, _ := driverField.Value.Interface().(interface{ BeforeRender() }); br != nil {
			br.BeforeRender()
		}
		return renderDriver(driverField, secretLoader)
	default:
		return render.Error(fmt.Errorf(
			"multiple drivers (%v) specified on output %s/%s",
			seqs.ToSlice(seqs.Map(seqs.FromSlice(driverFields), func(f Field) string { return f.Meta.Name })),
			output.GetNamespace(), output.GetName(),
		))
	}
}

func defaultPersistName(value reflect.Value, name string) {
	switch value.Kind() {
	case reflect.Pointer:
		defaultPersistName(derefAll(value), name)
	case reflect.Struct:
		if persistName := value.FieldByName("PersistName"); persistName.IsValid() && persistName.Kind() == reflect.String {
			if persistName.String() == "" {
				persistName.SetString(name)
			}
		}
	}
}
