// Copyright © 2025 Kube logging authors
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

package main

import (
	"log"

	"github.com/banzaicloud/config-reloader/pkg/configreloader"
	"github.com/banzaicloud/config-reloader/pkg/metrics"
)

func main() {
	cfg, err := configreloader.New()
	if err != nil {
		log.Fatalln(err)
	}

	err = cfg.Run()
	if err != nil {
		log.Fatalln(err)
	}

	if !*cfg.InitMode {
		err = metrics.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
