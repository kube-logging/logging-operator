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

	"github.com/siliconbrain/go-seqs/seqs"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config/model"
)

func destinationDefStmt(def model.DestinationDef) Renderer {
	return braceDefStmt("destination", def.Name, AllFrom(seqs.Map(seqs.FromSlice(def.Drivers), destinationDriverDefStmt)))
}

func destinationDriverDefStmt(drv model.DestinationDriver) Renderer {
	var args []Renderer
	switch drv := drv.(type) {
	case model.DestinationDriverAlt[model.SyslogDestinationDriver]:
		args = append(args, Quoted(drv.Alt.Host))
		if drv.Alt.Port != 0 {
			args = append(args, optionExpr("port", drv.Alt.Port))
		}
		if drv.Alt.Transport != "" {
			args = append(args, optionExpr("transport", drv.Alt.Transport))
		}
		if drv.Alt.CADir != "" {
			args = append(args, optionExpr("ca-dir", drv.Alt.CADir))
		}
		if drv.Alt.CAFile != "" {
			args = append(args, optionExpr("ca-file", drv.Alt.CAFile))
		}
		// CloseOnInput   *bool
		if flags := drv.Alt.Flags; len(flags) > 0 {
			args = append(args, flagsOption(flags))
		}
		// FlushLines     int
		// SoKeepalive    *bool
		// Suppress       int
		if drv.Alt.Template != "" {
			args = append(args, optionExpr("template", drv.Alt.Template))
		}
		if drv.Alt.TemplateEscape != nil {
			args = append(args, optionExpr("template-escape", *drv.Alt.TemplateEscape))
		}
		// TLS            *SyslogDestinationDriverTLS
		// TSFormat       string
		if drv.Alt.DiskBuffer != nil {
			args = append(args, diskBufferDef(*drv.Alt.DiskBuffer))
		}
	case model.DestinationDriverAlt[model.FileDestinationDriver]:
		args = append(args, Quoted(drv.Alt.Path))
	case model.DestinationDriverAlt[model.SumologicHTTPDriver]:
		args = append(args, Quoted(drv.Alt.CAFile))
	case model.DestinationDriverAlt[model.SumologicSyslogDriver]:
		args = append(args, Quoted(drv.Alt.CAFile))

	default:
		return Error(fmt.Errorf("unsupported destination driver %q", drv.Name()))
	}
	return parenDefStmt(drv.Name(), args...)
}

func diskBufferDef(def model.DiskBufferDef) Renderer {
	opts := []Renderer{
		optionExpr("disk-buf-size", def.DiskBufSize),
		optionExpr("reliable", def.Reliable),
	}
	if def.Compaction != nil {
		opts = append(opts, optionExpr("compaction", *def.Compaction))
	}
	if def.Dir != "" {
		opts = append(opts, optionExpr("dir", def.Dir))
	}
	if def.MemBufLength != nil {
		opts = append(opts, optionExpr("mem-buf-length", def.MemBufLength))
	}
	if def.MemBufSize != nil {
		opts = append(opts, optionExpr("mem-buf-size", def.MemBufSize))
	}
	if def.QOutSize != nil {
		opts = append(opts, optionExpr("qout-size", def.QOutSize))
	}
	return AllOf(String("disk-buffer("), SpaceSeparated(opts...), String(")"))
}
