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

package model

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func newRepoWithObjects(t *testing.T, objs ...client.Object) LoggingResourceRepository {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := v1beta1.AddToScheme(scheme); err != nil {
		t.Fatalf("add to scheme: %v", err)
	}
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	return LoggingResourceRepository{Client: cl, Logger: logr.Discard()}
}

func loggingFor() v1beta1.Logging {
	return v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: "test-logging", Namespace: "test"},
		Spec:       v1beta1.LoggingSpec{ControlNamespace: "test"},
	}
}

// Regression: when multiple SyslogNGConfig objects exist, the previously-
// associated one (recorded in Logging.Status.SyslogNGConfigName) must be kept
// as the active Configuration. Returning nil here de-associates the running
// aggregator, removes its finalizer, and breaks log forwarding.
func TestSyslogNGConfigFor_PreservesPreviouslyAssociatedConfigOnExcess(t *testing.T) {
	ns := "test"
	primary := &v1beta1.SyslogNGConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "primary", Namespace: ns},
	}
	excess := &v1beta1.SyslogNGConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "excess", Namespace: ns},
	}
	repo := newRepoWithObjects(t, primary, excess)

	logging := loggingFor()
	logging.Status.SyslogNGConfigName = "primary"

	cfg, excesses, err := repo.SyslogNGConfigFor(context.Background(), logging)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil || cfg.Name != "primary" {
		t.Fatalf("expected primary to be preserved, got: %+v", cfg)
	}
	if len(excesses) != 1 || excesses[0].Name != "excess" {
		t.Fatalf("expected excess to contain only 'excess', got: %+v", excesses)
	}
}

func TestSyslogNGConfigFor_NoPriorAssociationMarksAllExcess(t *testing.T) {
	ns := "test"
	a := &v1beta1.SyslogNGConfig{ObjectMeta: metav1.ObjectMeta{Name: "a", Namespace: ns}}
	b := &v1beta1.SyslogNGConfig{ObjectMeta: metav1.ObjectMeta{Name: "b", Namespace: ns}}
	repo := newRepoWithObjects(t, a, b)

	cfg, excesses, err := repo.SyslogNGConfigFor(context.Background(), loggingFor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Fatalf("expected no Configuration, got: %+v", cfg)
	}
	if len(excesses) != 2 {
		t.Fatalf("expected both as excess, got: %+v", excesses)
	}
}

func TestSyslogNGConfigFor_StaleAssociationFallsBackToAllExcess(t *testing.T) {
	ns := "test"
	a := &v1beta1.SyslogNGConfig{ObjectMeta: metav1.ObjectMeta{Name: "a", Namespace: ns}}
	b := &v1beta1.SyslogNGConfig{ObjectMeta: metav1.ObjectMeta{Name: "b", Namespace: ns}}
	repo := newRepoWithObjects(t, a, b)

	logging := loggingFor()
	logging.Status.SyslogNGConfigName = "deleted-config" // stale reference

	cfg, excesses, err := repo.SyslogNGConfigFor(context.Background(), logging)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Fatalf("expected no Configuration when prior name no longer exists, got: %+v", cfg)
	}
	if len(excesses) != 2 {
		t.Fatalf("expected both as excess, got %d", len(excesses))
	}
}

// Same regression on the Fluentd path.
func TestFluentdConfigFor_PreservesPreviouslyAssociatedConfigOnExcess(t *testing.T) {
	ns := "test"
	primary := &v1beta1.FluentdConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "primary", Namespace: ns},
	}
	excess := &v1beta1.FluentdConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "excess", Namespace: ns},
	}
	repo := newRepoWithObjects(t, primary, excess)

	logging := loggingFor()
	logging.Status.FluentdConfigName = "primary"

	cfg, excesses, err := repo.FluentdConfigFor(context.Background(), logging)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil || cfg.Name != "primary" {
		t.Fatalf("expected primary to be preserved, got: %+v", cfg)
	}
	if len(excesses) != 1 || excesses[0].Name != "excess" {
		t.Fatalf("expected excess to contain only 'excess', got: %+v", excesses)
	}
}
