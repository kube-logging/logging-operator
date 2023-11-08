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

package controllers_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/pborman/uuid"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var testNamespace = "test-" + uuid.New()[:8]
var controlNamespace = "control"

func TestMain(m *testing.M) {
	err := beforeSuite()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	code := m.Run()
	err = afterSuite()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func beforeSuite() error {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout)))

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		BinaryAssetsDirectory: os.Getenv("ENVTEST_BINARY_ASSETS"),
	}

	var err error

	cfg, err = testEnv.Start()
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("failed to start testenv, config is nil")
	}

	err = v1beta1.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}
	if k8sClient == nil {
		return fmt.Errorf("failed to create k8s config")
	}

	for _, ns := range []string{controlNamespace, testNamespace} {
		err := k8sClient.Create(context.TODO(), &v12.Namespace{
			ObjectMeta: v1.ObjectMeta{
				Name: ns,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func afterSuite() error {
	return testEnv.Stop()
}

// duplicateRequest returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func duplicateRequest(t *testing.T, inner reconcile.Reconciler, stopped *bool, errors chan<- error) reconcile.Reconciler {
	fn := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(ctx, req)
		if err != nil {
			if !*stopped {
				t.Logf("reconcile failure err: %+v req: %+v, result: %+v", err, req, result)
			}
			if errors != nil {
				errors <- err
			}
		}
		return result, err
	})
	return fn
}

// startTestManager adds recFn
func startTestManager(t *testing.T, mgr manager.Manager) (context.CancelFunc, *sync.WaitGroup) {
	stop, cf := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mgr.Start(stop); err != nil {
			t.Logf("%+v", err)
		}
	}()
	return cf, wg
}
