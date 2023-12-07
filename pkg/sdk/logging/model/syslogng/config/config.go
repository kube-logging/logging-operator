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
	"io"
	"reflect"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/siliconbrain/go-seqs/seqs"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
)

func RenderConfigInto(in Input, out io.Writer) error {
	if in.SecretLoaderFactory == nil {
		return errors.New("no secret loader factory provided")
	}
	rnd, err := configRenderer(in)
	if err != nil {
		return err
	}
	return rnd(render.RenderContext{
		Out:        out,
		IndentWith: "    ",
	})
}

type Input struct {
	Name                string
	Namespace           string
	SyslogNGSpec        *v1beta1.SyslogNGSpec
	ClusterOutputs      []v1beta1.SyslogNGClusterOutput
	Outputs             []v1beta1.SyslogNGOutput
	ClusterFlows        []v1beta1.SyslogNGClusterFlow
	Flows               []v1beta1.SyslogNGFlow
	SecretLoaderFactory SecretLoaderFactory
	SourcePort          int
}

type SecretLoaderFactory interface {
	SecretLoaderForNamespace(namespace string) secret.SecretLoader
}

const configVersion = "current"
const sourceName = "main_input"

func configRenderer(in Input) (render.Renderer, error) {
	if in.SyslogNGSpec == nil {
		return nil, errors.New("missing syslog-ng spec")
	}

	var errs error

	// TODO: this should happen at the spec level, in something like `SyslogNGSpec.FinalGlobalOptions() GlobalOptions`
	if in.SyslogNGSpec.Metrics != nil {
		setDefault(&in.SyslogNGSpec.GlobalOptions, &v1beta1.GlobalOptions{})
		if in.SyslogNGSpec.GlobalOptions.StatsFreq != nil ||
			in.SyslogNGSpec.GlobalOptions.StatsLevel != nil {
			return nil, errors.New("stats_freq and stats_level are not supported anymore, please use stats.level and stats.freq")
		}

		setDefault(&in.SyslogNGSpec.GlobalOptions.Stats, &v1beta1.Stats{})
		setDefault(&in.SyslogNGSpec.GlobalOptions.Stats.Freq, amp(0))
		setDefault(&in.SyslogNGSpec.GlobalOptions.Stats.Level, amp(2))
	}

	globalOptions := renderAny(in.SyslogNGSpec.GlobalOptions, in.SecretLoaderFactory.SecretLoaderForNamespace(in.Namespace))

	destinationDefs := make([]render.Renderer, 0, len(in.ClusterOutputs)+len(in.Outputs))
	clusterOutputRefs := make(map[string]types.NamespacedName, len(in.ClusterOutputs))
	for _, co := range in.ClusterOutputs {
		clusterOutputRefs[co.Name] = types.NamespacedName{
			Namespace: co.Namespace,
			Name:      co.Name,
		}
		destinationDefs = append(destinationDefs, renderClusterOutput(co, in.SecretLoaderFactory))
	}
	for _, o := range in.Outputs {
		destinationDefs = append(destinationDefs, renderOutput(o, in.SecretLoaderFactory))
	}

	logDefs := make([]render.Renderer, 0, len(in.ClusterFlows)+len(in.Flows))
	for _, cf := range in.ClusterFlows {
		if err := validateClusterOutputs(clusterOutputRefs, client.ObjectKeyFromObject(&cf).String(), cf.Spec.GlobalOutputRefs); err != nil {
			errs = errors.Append(errs, err)
		}
		logDefs = append(logDefs, renderClusterFlow(in.Name, clusterOutputRefs, sourceName, cf, in.SecretLoaderFactory))
	}
	for _, f := range in.Flows {
		if err := validateClusterOutputs(clusterOutputRefs, client.ObjectKeyFromObject(&f).String(), f.Spec.GlobalOutputRefs); err != nil {
			errs = errors.Append(errs, err)
		}
		logDefs = append(logDefs, renderFlow(in.Name, clusterOutputRefs, sourceName, keyDelim(in.SyslogNGSpec.JSONKeyDelimiter), f, in.SecretLoaderFactory))
	}

	if in.SyslogNGSpec.JSONKeyPrefix == "" {
		in.SyslogNGSpec.JSONKeyPrefix = "json" + keyDelim(in.SyslogNGSpec.JSONKeyDelimiter)
	}

	sourceParsers := []render.Renderer{
		renderDriver(
			Field{
				Value: reflect.ValueOf(JSONParser{
					Prefix:       in.SyslogNGSpec.JSONKeyPrefix,
					KeyDelimiter: in.SyslogNGSpec.JSONKeyDelimiter,
				}),
			}, nil),
	}

	if in.SyslogNGSpec.SourceDateParser != nil {
		setDefault(&in.SyslogNGSpec.SourceDateParser.Format, amp("%FT%T.%f%z"))
		setDefault(&in.SyslogNGSpec.SourceDateParser.Template,
			amp(fmt.Sprintf("${%stime}", in.SyslogNGSpec.JSONKeyPrefix)))
		sourceParsers = append(sourceParsers, renderDriver(
			Field{
				Value: reflect.ValueOf(DateParser{
					Format:   *in.SyslogNGSpec.SourceDateParser.Format,
					Template: *in.SyslogNGSpec.SourceDateParser.Template,
				}),
			}, nil))
	}

	for _, sm := range in.SyslogNGSpec.SourceMetrics {
		if sm.Labels == nil {
			sm.Labels = make(filter.ArrowMap, 0)
		}
		sm.Labels["logging"] = in.Name
		sourceParsers = append(sourceParsers, renderDriver(Field{
			Value: reflect.ValueOf(sm),
		}, nil))
	}

	return render.AllFrom(seqs.Intersperse(
		seqs.Filter(
			seqs.Concat(
				seqs.FromValues(
					versionStmt(configVersion),
					includeStmt("scl.conf"),
					globalOptionsDefStmt(globalOptions...),
					sourceDefStmt(sourceName,
						channelDefStmt(
							sourceDefStmt("", renderDriver(Field{
								Value: reflect.ValueOf(NetworkSourceDriver{
									Transport:      "tcp",
									Port:           uint16(in.SourcePort),
									MaxConnections: in.SyslogNGSpec.MaxConnections,
									LogIWSize:      logIWSizeCalculator(in),
									Flags:          []string{"no-parse"},
								}),
							}, nil)),
							[]render.Renderer{
								parserDefStmt("", render.AllOf(sourceParsers...)),
							},
						),
					),
				),
				seqs.FromSlice(destinationDefs),
				seqs.FromSlice(logDefs),
			),
			func(rnd render.Renderer) bool { return rnd != nil },
		),
		render.Line(render.Empty),
	)), errs
}

func versionStmt(version string) render.Renderer {
	return render.Line(render.Formatted("@version: %s", version))
}

func includeStmt(file string) render.Renderer {
	return render.Line(render.Formatted("@include %q", file))
}

func globalOptionsDefStmt(options ...render.Renderer) render.Renderer {
	if len(options) == 0 {
		return nil
	}
	return braceDefStmt("options", "", render.AllFrom(seqs.Map(seqs.FromSlice(options),
		func(rnd render.Renderer) render.Renderer {
			return render.Line(render.AllOf(rnd, render.String(";")))
		},
	)))
}

func keyDelim(delim string) string {
	if delim != "" {
		return delim
	}
	return "."
}

func setDefault[T comparable](ptr *T, def T) {
	var zero T
	if *ptr == zero {
		*ptr = def
	}
}

func amp[T any](v T) *T {
	return &v
}

func logIWSizeCalculator(in Input) int {
	if in.SyslogNGSpec.MaxConnections != 0 && in.SyslogNGSpec.LogIWSize == 0 {
		return in.SyslogNGSpec.MaxConnections * 100
	}
	return in.SyslogNGSpec.LogIWSize
}
