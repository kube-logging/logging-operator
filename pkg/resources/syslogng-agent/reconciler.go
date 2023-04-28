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

package syslogng_agent

import "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

func Reconcile(agent v1beta1.SyslogNGAgent) {
    // 1. load defaults and merge actual values on top
    // 2. generate resources
    //   2.1 generate most of static resources
    //   2.2 generate config
}
