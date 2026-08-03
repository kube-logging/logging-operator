// Copyright © 2026 Kube logging authors
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

package fixture

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

// recorder keeps whatever is handed to Create. The embedded nil Client is
// never reached: the helpers under test only create.
type recorder struct {
	client.Client
	created []client.Object
}

func (r *recorder) Create(_ context.Context, obj client.Object, _ ...client.CreateOption) error {
	r.created = append(r.created, obj)
	return nil
}

// The builders have to keep producing what common creates today, object for
// object and in the same order, or a suite that moves onto them stops testing
// what it used to. Nothing else pins that: the shape tests above assert the
// fields this package chose, not that the choice matches the live path.
func TestBuildersMatchCommon(t *testing.T) {
	const (
		nsInfra  = "infra"
		nsTenant = "tenant"
		release  = "fluentbit-multitenant"
	)
	labels := map[string]string{"my-unique-label": "log-producer"}
	buffer := func() *output.Buffer {
		tags := "time"
		return &output.Buffer{Tags: &tags, Timekey: "1s", TimekeyWait: "0s"}
	}

	t.Run("infra", func(t *testing.T) {
		rec := &recorder{}
		common.LoggingInfra(context.Background(), t, rec, nsInfra, release, "tag_infra", buffer(), labels)

		require.Equal(t, rec.created, LoggingInfra(nsInfra, release, "tag_infra", buffer(), labels))
	})

	t.Run("tenant", func(t *testing.T) {
		rec := &recorder{}
		common.LoggingTenant(context.Background(), t, rec, nsTenant, nsInfra, release, "tag_tenant", buffer(), labels)

		require.Equal(t, rec.created, LoggingTenant(nsTenant, nsInfra, release, "tag_tenant", buffer(), labels))
	})

	t.Run("route", func(t *testing.T) {
		rec := &recorder{}
		common.LoggingRoute(context.Background(), t, rec)

		require.Equal(t, rec.created, []client.Object{LoggingRoute()})
	})
}
