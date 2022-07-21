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

func sourceDefStmt(def model.SourceDef) Renderer {
	return braceDefStmt("source", def.Name, AllFrom(seqs.Map(seqs.FromSlice(def.Drivers), sourceDriverDefStmt)))
}

func sourceDriverDefStmt(drv model.SourceDriver) Renderer {
	var args []Renderer
	switch drv := drv.(type) {
	case model.SourceDriverAlt[model.NetworkSourceDriver]:
		if transport := drv.Alt.Transport; transport != "" {
			args = append(args, optionExpr("transport", transport))
		}
		if ip := drv.Alt.IP; ip != "" {
			args = append(args, optionExpr("ip", ip))
		}
		if port := drv.Alt.Port; port != 0 {
			args = append(args, optionExpr("port", port))
		}
		if flags := drv.Alt.Flags; len(flags) > 0 {
			args = append(args, flagsOption(flags))
		}
	default:
		return Error(fmt.Errorf("unsupported source driver %q", drv.Name()))
	}
	return parenDefStmt(drv.Name(), args...)
}
