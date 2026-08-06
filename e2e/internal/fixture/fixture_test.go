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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestImageNamesMatchCommon(t *testing.T) {
	require.Equal(t, common.FluentdImageRepo, FluentdImageRepo)
	require.Equal(t, common.ConfigReloaderRepo, ConfigReloaderRepo)
	require.Equal(t, common.SyslogNGReloaderRepo, SyslogNGReloaderRepo)
	require.Equal(t, common.FluentdDrainWatchRepo, FluentdDrainWatchRepo)
	require.Equal(t, common.NodeExporterRepo, NodeExporterRepo)

	for _, tag := range []string{
		common.FluentdImageTag, common.ConfigReloaderTag, common.SyslogNGReloaderTag,
		common.FluentdDrainWatchTag, common.NodeExporterTag,
	} {
		require.Equal(t, localTag, tag)
	}
}

// Each helper has to equal the literal it replaces at the call sites, or the
// sweep changes which image a suite runs.
func TestImageHelpersMatchTheLiterals(t *testing.T) {
	require.Equal(t, v1beta1.ImageSpec{Repository: common.FluentdImageRepo, Tag: common.FluentdImageTag}, FluentdImage())
	require.Equal(t, v1beta1.ImageSpec{Repository: common.ConfigReloaderRepo, Tag: common.ConfigReloaderTag}, ConfigReloaderImage())
	require.Equal(t, v1beta1.ImageSpec{Repository: common.SyslogNGReloaderRepo, Tag: common.SyslogNGReloaderTag}, SyslogNGReloaderImage())
	require.Equal(t, v1beta1.ImageSpec{Repository: common.FluentdDrainWatchRepo, Tag: common.FluentdDrainWatchTag}, DrainWatchImage())
	require.Equal(t, v1beta1.ImageSpec{Repository: common.NodeExporterRepo, Tag: common.NodeExporterTag}, NodeExporterImage())

	require.Equal(t,
		&v1beta1.BasicImageSpec{Repository: common.NodeExporterRepo, Tag: common.NodeExporterTag},
		Basic(NodeExporterImage()))
}

func TestReceiverURL(t *testing.T) {
	require.Equal(t, "http://e2e-test-receiver:8080/test.tag", ReceiverURL("e2e", "test.tag"))
}
