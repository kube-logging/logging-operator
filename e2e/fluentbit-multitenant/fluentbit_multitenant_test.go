// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package fluentbit_multitenant

import (
	"testing"

	"github.com/kube-logging/logging-operator/e2e/internal/fixture"
	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	release   = "fluentbit-multitenant"
	nsInfra   = "infra"
	nsTenant  = "tenant"
	tagInfra  = "tag_infra"
	tagTenant = "tag_tenant"
)

var producerLabels = map[string]string{"my-unique-label": "log-producer"}

// This is the one suite whose buffer leaves Type unset, so it is written out
// rather than taken from fixture.Buffer, whose "file" default would switch the
// outputs from memory buffering.
func realTimeBuffer() *output.Buffer {
	tags := "time"
	return &output.Buffer{Tags: &tags, Timekey: "1s", TimekeyWait: "0s"}
}

func TestFluentbitSingleTenantPlusInfra(t *testing.T) {
	env := harness.New(t).
		WithCluster(release).
		WithRelease(release).
		WithControlNamespace(nsInfra).
		WithNamespaces(nsTenant).
		Start()

	buffer := realTimeBuffer()
	env.Create(fixture.LoggingInfra(nsInfra, release, tagInfra, buffer, producerLabels)...)
	env.Create(fixture.LoggingTenant(nsTenant, nsInfra, release, tagTenant, buffer, producerLabels)...)
	env.Create(fixture.LoggingRoute())

	env.StartLogProducer(nsTenant, producerLabels)

	env.WaitForRunning(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.Aggregator(nsInfra),
		wait.Aggregator(nsTenant),
	)
	env.WaitForReceiverLogs(tagInfra, tagTenant)
}
