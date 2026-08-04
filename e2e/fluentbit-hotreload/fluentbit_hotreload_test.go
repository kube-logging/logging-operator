// Copyright © 2024 Kube logging authors
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

package fluentbit_hotreload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kube-logging/logging-operator/e2e/internal/fixture"
	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	release   = "fluentbit-hotreload"
	nsInfra   = "infra"
	nsTenant  = "tenant"
	tagInfra  = "tag_infra"
	tagTenant = "tag_tenant"
)

var producerLabels = map[string]string{"my-unique-label": "log-producer"}

// Type is left unset here as it is in fluentbit-multitenant; fixture.Buffer's
// "file" default would switch the outputs from memory buffering.
func realTimeBuffer() *output.Buffer {
	tags := "time"
	return &output.Buffer{Tags: &tags, Timekey: "1s", TimekeyWait: "0s"}
}

func TestFluentbitHotReload(t *testing.T) {
	env := harness.New(t).
		WithCluster(release).
		WithRelease(release).
		WithControlNamespace(nsInfra).
		WithNamespaces(nsTenant).
		Start()

	buffer := realTimeBuffer()
	env.Create(fixture.LoggingInfra(nsInfra, release, tagInfra, buffer, producerLabels)...)
	env.Create(fixture.LoggingTenant(nsTenant, nsInfra, release, tagTenant, buffer, producerLabels)...)

	env.StartLogProducer(nsTenant, producerLabels)

	env.WaitForRunning(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(nsInfra),
		wait.FluentdAggregator(nsTenant),
	)

	// No route yet, so the tenant's logs have nowhere to go.
	env.Receiver.MustReceive(tagInfra)
	env.Receiver.MustNotReceive(tagTenant)

	env.Create(fixture.LoggingRoute())
	env.Receiver.MustReceive(tagTenant)

	ds := &appsv1.DaemonSet{}
	assert.NoError(t, env.Client.Get(env.Ctx, types.NamespacedName{
		Namespace: nsInfra,
		Name:      "infra-fluentbit",
	}, ds))
	assert.Equal(t, int64(1), ds.Generation, "generation should not be incremented for a reloadable agent")
}
