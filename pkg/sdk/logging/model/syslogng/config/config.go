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
	"errors"
	"fmt"
	"io"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/output"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/siliconbrain/go-seqs/seqs"
)

func RenderConfigInto(in Input, out io.Writer) error {
	if in.SecretLoaderFactory == nil {
		return errors.New("no secret loader factory provided")
	}
	rnd, err := configRenderer(in)
	if err != nil {
		return err
	}
	return rnd(RenderContext{
		Out:        out,
		IndentWith: "    ",
	})
}

type Input struct {
	Logging             v1beta1.Logging
	ClusterOutputs      []v1beta1.SyslogNGClusterOutput
	Outputs             []v1beta1.SyslogNGOutput
	ClusterFlows        []v1beta1.SyslogNGClusterFlow
	Flows               []v1beta1.SyslogNGFlow
	SecretLoaderFactory SecretLoaderFactory
}

type SecretLoaderFactory interface {
	OutputSecretLoaderForNamespace(namespace string) secret.SecretLoader
}

const configVersion = "3.37"

func configRenderer(in Input) (Renderer, error) {
	if in.Logging.Spec.SyslogNGSpec == nil {
		return nil, errors.New("missing syslog-ng spec")
	}

	srcDef := model.SourceDef{
		Name: "main_input",
		Drivers: []model.SourceDriver{
			model.NewSourceDriver(model.NetworkSourceDriver{
				Transport: "tcp",
				Port:      601,
				Flags:     []string{"no-parse"},
			}),
		},
	}

	destinationDefs := make([]model.DestinationDef, 0, len(in.ClusterOutputs)+len(in.Outputs))
	for _, co := range in.ClusterOutputs {
		def, err := clusterOutputToDestinationDef(in.SecretLoaderFactory, co)
		if err != nil {
			return nil, err
		}
		destinationDefs = append(destinationDefs, def)
	}
	for _, o := range in.Outputs {
		def, err := outputToDestinationDef(in.SecretLoaderFactory, o)
		if err != nil {
			return nil, err
		}
		destinationDefs = append(destinationDefs, def)
	}

	logDefs := make([]model.LogDef, 0, len(in.ClusterFlows)+len(in.Flows))
	for _, cf := range in.ClusterFlows {
		def, err := clusterFlowToLogDef(srcDef.Name, cf)
		if err != nil {
			return nil, err
		}
		logDefs = append(logDefs, def)
	}
	for _, f := range in.Flows {
		def, err := flowToLogDef(in.Logging.Spec.ControlNamespace, srcDef.Name, f)
		if err != nil {
			return nil, err
		}
		logDefs = append(logDefs, def)
	}

	return AllOf(
		versionStmt(configVersion),
		Line(Empty),
		sourceDefStmt(srcDef),
		Line(Empty),
		AllFrom(seqs.Map(seqs.FromSlice(destinationDefs), destinationDefStmt)),
		Line(Empty),
		AllFrom(seqs.Map(seqs.FromSlice(logDefs), logDefStmt)),
	), nil
}

func clusterOutputToDestinationDef(secretLoaderFactory SecretLoaderFactory, o v1beta1.SyslogNGClusterOutput) (def model.DestinationDef, err error) {
	def.Name = fmt.Sprintf("clusteroutput_%s_%s", o.Namespace, o.Name)
	driver, err := outputSpecToDriver(secretLoaderFactory.OutputSecretLoaderForNamespace(o.Namespace), o.Spec.SyslogNGOutputSpec)
	if err != nil {
		return
	}
	def.Drivers = []model.DestinationDriver{driver}
	return
}

func outputToDestinationDef(secretLoaderFactory SecretLoaderFactory, o v1beta1.SyslogNGOutput) (def model.DestinationDef, err error) {
	def.Name = fmt.Sprintf("output_%s_%s", o.Namespace, o.Name)
	driver, err := outputSpecToDriver(secretLoaderFactory.OutputSecretLoaderForNamespace(o.Namespace), o.Spec)
	if err != nil {
		return
	}
	def.Drivers = []model.DestinationDriver{driver}
	return
}

func outputSpecToDriver(secretLoader secret.SecretLoader, s v1beta1.SyslogNGOutputSpec) (model.DestinationDriver, error) {
	switch {
	case s.Syslog != nil:
		caDir, err := loadSecretIfNotNil(secretLoader, s.Syslog.CaDir)
		if err != nil {
			return nil, err
		}
		caFile, err := loadSecretIfNotNil(secretLoader, s.Syslog.CaFile)
		if err != nil {
			return nil, err
		}
		return model.NewDestinationDriver(model.SyslogDestinationDriver{
			Host:           s.Syslog.Host,
			Port:           s.Syslog.Port,
			Transport:      s.Syslog.Transport,
			CADir:          caDir,
			CAFile:         caFile,
			CloseOnInput:   s.Syslog.CloseOnInput,
			Flags:          s.Syslog.Flags,
			FlushLines:     s.Syslog.FlushLines,
			SoKeepalive:    s.Syslog.SoKeepalive,
			Suppress:       s.Syslog.Suppress,
			Template:       s.Syslog.Template,
			TemplateEscape: s.Syslog.TemplateEscape,
			TLS:            (*model.SyslogDestinationDriverTLS)(s.Syslog.TLS), // TODO
			TSFormat:       s.Syslog.TSFormat,
			DiskBuffer:     outputDiskBufferToModelDiskBuffer(s.Syslog.DiskBuffer),
		}), nil
	default:
		return nil, errors.New("unsupported output type")
	}
}

func clusterFlowToLogDef(sourceName string, f v1beta1.SyslogNGClusterFlow) (def model.LogDef, err error) {
	//def.Name = fmt.Sprintf("clusterflow_%s_%s", f.Namespace, f.Name)
	def.SourceNames = []string{sourceName}
	def.OptionalElements = append(def.OptionalElements, model.NewLogElement(model.ParserDef{
		Parsers: []model.Parser{model.NewParser(model.JSONParser{})},
	}))
	if match := f.Spec.Match; match != nil {
		def.OptionalElements = append(def.OptionalElements, model.NewLogElement(model.FilterDef{
			Expr: filterExprFromMatchExpr(filter.MatchExpr(*match)),
		}))
	}
	def.OptionalElements = append(def.OptionalElements, seqs.ToSlice(seqs.Map(seqs.FromSlice(f.Spec.Filters), loggingFilterToLogElement))...)
	for _, o := range f.Spec.GlobalOutputRefs {
		def.DestinationNames = append(def.DestinationNames, fmt.Sprintf("clusteroutput_%s_%s", f.Namespace, o))
	}
	return
}

func flowToLogDef(controlNS string, sourceName string, f v1beta1.SyslogNGFlow) (def model.LogDef, err error) {
	//def.Name = fmt.Sprintf("flow_%s_%s", f.Namespace, f.Name)
	def.SourceNames = []string{sourceName}
	def.OptionalElements = append(def.OptionalElements, model.NewLogElement(model.ParserDef{
		Parsers: []model.Parser{model.NewParser(model.JSONParser{})},
	}))
	def.OptionalElements = append(def.OptionalElements, model.NewLogElement(model.FilterDef{
		Expr: model.NewFilterExpr(model.FilterExprMatch{
			Pattern: f.Namespace,
			Scope:   model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue(".kubernetes.namespace_name")),
			Type:    "string",
		}),
	}))
	if match := f.Spec.Match; match != nil {
		def.OptionalElements = append(def.OptionalElements, model.NewLogElement(model.FilterDef{
			Expr: filterExprFromMatchExpr(filter.MatchExpr(*match)),
		}))
	}
	def.OptionalElements = append(def.OptionalElements, seqs.ToSlice(seqs.Map(seqs.FromSlice(f.Spec.Filters), loggingFilterToLogElement))...)
	for _, o := range f.Spec.GlobalOutputRefs {
		def.DestinationNames = append(def.DestinationNames, fmt.Sprintf("clusteroutput_%s_%s", controlNS, o))
	}
	for _, o := range f.Spec.LocalOutputRefs {
		def.DestinationNames = append(def.DestinationNames, fmt.Sprintf("output_%s_%s", f.Namespace, o))
	}
	return
}

func filterExprFromMatchExpr(expr filter.MatchExpr) model.FilterExpr {
	switch {
	case len(expr.And) > 0:
		return model.NewFilterExpr(model.FilterExprAnd(seqs.ToSlice(seqs.Map(seqs.FromSlice(expr.And), filterExprFromMatchExpr))))
	case expr.Not != nil:
		return model.NewFilterExpr(model.FilterExprNot{Expr: filterExprFromMatchExpr(filter.MatchExpr(*expr.Not))})
	case len(expr.Or) > 0:
		return model.NewFilterExpr(model.FilterExprOr(seqs.ToSlice(seqs.Map(seqs.FromSlice(expr.Or), filterExprFromMatchExpr))))
	case expr.Regexp != nil:
		m := model.FilterExprMatch{
			Pattern: expr.Regexp.Pattern,
			Type:    expr.Regexp.Type,
			Flags:   expr.Regexp.Flags,
		}
		switch {
		case expr.Regexp.Template != "":
			m.Scope = model.NewFilterExprMatchScope(model.FilterExprMatchScopeTemplate(expr.Regexp.Template))
		case expr.Regexp.Value != "":
			m.Scope = model.NewFilterExprMatchScope(model.FilterExprMatchScopeValue(expr.Regexp.Value))
		}
		return model.NewFilterExpr(m)
	default:
		return nil
	}
}

func loggingFilterToLogElement(f v1beta1.SyslogNGFilter) model.LogElement {
	switch {
	case f.Match != nil:
		return model.NewLogElement(model.FilterDef{
			Expr: filterExprFromMatchExpr(filter.MatchExpr(*f.Match)),
		})
	case f.Rewrite != nil:
		return model.NewLogElement(model.RewriteDef{
			Rules: []model.RewriteRule{rewriteRuleFromRewriteConfig(*f.Rewrite)},
		})
	default:
		return nil
	}
}

func rewriteRuleFromRewriteConfig(cfg filter.RewriteConfig) model.RewriteRule {
	switch {
	case cfg.Rename != nil:
		return model.NewRewriteRule(model.RenameRule{
			OldFieldName: cfg.Rename.OldFieldName,
			NewFieldName: cfg.Rename.NewFieldName,
			Condition:    rewriteConditionFromMatchExpr(cfg.Rename.Condition),
		})
	case cfg.Set != nil:
		return model.NewRewriteRule(model.SetRule{
			FieldName: cfg.Set.FieldName,
			Value:     cfg.Set.Value,
			Condition: rewriteConditionFromMatchExpr(cfg.Set.Condition),
		})
	case cfg.Substitute != nil:
		return model.NewRewriteRule(model.SubstituteRule{
			FieldName:   cfg.Substitute.FieldName,
			Pattern:     cfg.Substitute.Pattern,
			Replacement: cfg.Substitute.Replacement,
			Type:        cfg.Substitute.Type,
			Flags:       cfg.Substitute.Flags,
			Condition:   rewriteConditionFromMatchExpr(cfg.Substitute.Condition),
		})
	case cfg.Unset != nil:
		return model.NewRewriteRule(model.UnsetRule{
			FieldName: cfg.Unset.FieldName,
			Condition: rewriteConditionFromMatchExpr(cfg.Unset.Condition),
		})
	default:
		return nil
	}
}

func outputDiskBufferToModelDiskBuffer(b *output.SyslogNGDiskBuffer) *model.DiskBufferDef {
	if b == nil {
		return nil
	}
	return &model.DiskBufferDef{
		Reliable:     b.Reliable,
		Compaction:   b.Compaction,
		Dir:          b.Dir,
		DiskBufSize:  b.DiskBufSize,
		MemBufLength: b.MemBufLength,
		MemBufSize:   b.MemBufSize,
		QOutSize:     b.QOutSize,
	}
}

func rewriteConditionFromMatchExpr(c *filter.MatchExpr) *model.RewriteCondition {
	if c == nil {
		return nil
	}
	return &model.RewriteCondition{
		Expr: filterExprFromMatchExpr(*c),
	}
}

func loadSecretIfNotNil(loader secret.SecretLoader, secret *secret.Secret) (string, error) {
	if secret == nil {
		return "", nil
	}
	return loader.Load(secret)
}
