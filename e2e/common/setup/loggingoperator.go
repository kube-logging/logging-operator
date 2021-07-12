// Copyright © 2021 Banzai Cloud
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

package setup

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/banzaicloud/logging-operator/pkg/sdk/resourcebuilder"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	logrtesting "github.com/go-logr/logr/testing"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/banzaicloud/logging-operator/e2e/common"
)

func LoggingOperator(t *testing.T, c common.Cluster, opts ...LoggingOperatorOption) {
	options := LoggingOperatorOptions{
		Config: resourcebuilder.ComponentConfig{
			Enabled:   utils.BoolPointer(true),
			Namespace: "default",
		},
		Logger: logrtesting.TestLogger{
			T: t,
		},
		Parent: &parentObject{
			Name: "test",
		},
		PollInterval: 5 * time.Second,
		Timeout:      2 * time.Minute,
	}

	if img := os.Getenv("LOGGING_OPERATOR_IMAGE"); img != "" {
		require.NoError(t, c.LoadImages(img))

		if options.Config.ContainerOverrides == nil {
			options.Config.ContainerOverrides = new(types.ContainerBase)
		}
		options.Config.ContainerOverrides.Image = img
		options.Config.ContainerOverrides.PullPolicy = corev1.PullNever
	}

	for _, opt := range opts {
		opt.ApplyToLoggingOperatorOptions(&options)
	}

	resourceBuilders := resourcebuilder.ResourceBuilders(options.Parent, &options.Config)
	reconciler := reconciler.NewGenericReconciler(c.GetClient(), options.Logger, reconciler.ReconcilerOpts{
		Scheme: c.GetScheme(),
	})
	for _, rb := range resourceBuilders {
		obj, ds, err := rb()
		require.NoError(t, err)
		res, err := reconciler.ReconcileResource(obj, ds)
		require.NoError(t, err)
		require.Nil(t, res)
	}

	require.Eventually(t, func() bool {
		var ps corev1.PodList
		if err := c.GetClient().List(context.Background(), &ps, client.MatchingLabels{
			"banzaicloud.io/operator": options.Parent.GetName() + "-logging-operator",
		}); err != nil {
			t.Logf("failed to list logging operator pods: %v", err)
			return false
		}
		for _, p := range ps.Items {
			if p.Status.Phase == corev1.PodRunning {
				return true
			}
		}
		if len(ps.Items) > 0 {
			t.Log("logging operator is not running")
		}
		return false
	}, options.Timeout, options.PollInterval)
}

type LoggingOperatorOption interface {
	ApplyToLoggingOperatorOptions(options *LoggingOperatorOptions)
}

type LoggingOperatorOptionFunc func(*LoggingOperatorOptions)

func (fn LoggingOperatorOptionFunc) ApplyToLoggingOperatorOptions(options *LoggingOperatorOptions) {
	fn(options)
}

type LoggingOperatorOptions struct {
	Config       resourcebuilder.ComponentConfig
	Logger       logr.Logger
	Parent       reconciler.ResourceOwner
	PollInterval time.Duration
	Timeout      time.Duration
}

type parentObject struct {
	common.PanicObject
	Name string
}

func (o *parentObject) GetName() string {
	return o.Name
}
