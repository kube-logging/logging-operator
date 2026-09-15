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

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

// TestHealthProbesRespondOK exercises the same manager wiring used in main():
// a HealthProbeBindAddress plus AddHealthzCheck/AddReadyzCheck registrations,
// and verifies the resulting /healthz and /readyz endpoints respond 200 OK.
func TestHealthProbesRespondOK(t *testing.T) {
	testEnv := &envtest.Environment{
		BinaryAssetsDirectory: os.Getenv("ENVTEST_BINARY_ASSETS"),
	}
	cfg, err := testEnv.Start()
	if err != nil {
		t.Fatalf("unable to start test environment: %v", err)
	}
	t.Cleanup(func() {
		if err := testEnv.Stop(); err != nil {
			t.Errorf("unable to stop test environment: %v", err)
		}
	})

	const healthAddr = "127.0.0.1:18734"

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: "0"},
		HealthProbeBindAddress: healthAddr,
	})
	if err != nil {
		t.Fatalf("unable to create manager: %v", err)
	}

	if err := setupHealthChecks(mgr); err != nil {
		t.Fatalf("unable to set up health checks: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	errCh := make(chan error, 1)
	go func() { errCh <- mgr.Start(ctx) }()

	for _, path := range []string{"healthz", "readyz"} {
		var resp *http.Response
		var reqErr error
		for range 50 {
			resp, reqErr = http.Get(fmt.Sprintf("http://%s/%s", healthAddr, path))
			if reqErr == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if reqErr != nil {
			t.Fatalf("/%s endpoint never became available: %v", path, reqErr)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("/%s returned status %d, want %d", path, resp.StatusCode, http.StatusOK)
		}
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Errorf("manager exited with error: %v", err)
	}
}
