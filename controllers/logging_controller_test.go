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
	"strings"
	"sync"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/controllers"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/output"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	err              error
	mgr              ctrl.Manager
	requests         chan reconcile.Request
	stopMgr          chan struct{}
	mgrStopped       *sync.WaitGroup
	reconcilerErrors chan error
	g                gomega.GomegaWithT
)

func TestFluentdResourcesCreatedAndRemoved(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	defer ensureCreated(t, logging)()

	cm := &corev1.Secret{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.SecretConfigName), cm)()

	g.Expect(cm.Data["fluent.conf"]).Should(gomega.And(
		gomega.ContainSubstring("@include /fluentd/etc/input.conf"),
		gomega.ContainSubstring("@include /fluentd/app-config/*"),
		gomega.ContainSubstring("@include /fluentd/etc/devnull.conf"),
	))

	deployment := &appsv1.StatefulSet{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.StatefulSetName), deployment)()
}

func TestSingleFlowWithoutOutputRefs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestSingleFlowWithoutExistingLoggingRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			LoggingRef: "nonexistent",
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).ShouldNot(gomega.ContainSubstring("namespace " + testNamespace))
}

func TestSingleFlowWithOutputRefDefaultLoggingRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{"test-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestSingleFlowWithClusterOutput(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-cluster-output",
			Namespace: controlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: output.NewNullOutputConfig(),
			},
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{"test-cluster-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestClusterFlowWithNamespacedOutput(t *testing.T) {
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: controlNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{"test-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	err := wait.Poll(time.Second, time.Second*3, func() (bool, error) {
		select {
		case err := <-reconcilerErrors:
			expected := "referenced output not found: test-output"
			if !strings.Contains(err.Error(), expected) {
				return false, errors.Errorf("expected `%s` but received `%s`", expected, err.Error())
			} else {
				return true, nil
			}
		case <-time.After(100 * time.Millisecond):
			return false, nil
		}
	})

	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestSingleFlowWithOutputRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:              "someloggingref",
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			LoggingRef:       "someloggingref",
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			LoggingRef: "someloggingref",
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{"test-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestSingleFlowDefaultLoggingRefInvalidOutputRef(t *testing.T) {
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{"test-output-nonexistent"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, flow)()

	err := wait.Poll(time.Second, time.Second*3, func() (bool, error) {
		select {
		case err := <-reconcilerErrors:
			expected := "referenced output not found: test-output-nonexistent"
			if !strings.Contains(err.Error(), expected) {
				return false, errors.Errorf("expected `%s` but received `%s`", expected, err.Error())
			} else {
				return true, nil
			}
		case <-time.After(100 * time.Millisecond):
			return false, nil
		}
	})

	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestSingleFlowWithSecretInOutput(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			S3OutputConfig: &output.S3OutputConfig{
				AwsAccessKey: &secret.Secret{
					ValueFrom: &secret.ValueFrom{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "topsecret",
							},
							Key: "key",
						},
					},
				},
				AwsSecretKey: &secret.Secret{
					MountFrom: &secret.ValueFrom{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "topsecret",
							},
							Key: "key",
						},
					},
				},
				SharedCredentials: &output.S3SharedCredentials{},
			},
		},
	}
	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			OutputRefs: []string{
				"test-output",
			},
		},
	}
	topsecret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "topsecret",
			Namespace: testNamespace,
		},
		StringData: map[string]string{
			"key": "topsecretdata",
		},
	}
	defer ensureCreated(t, logging)()
	defer ensureCreated(t, topsecret)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	secretKey := fmt.Sprintf("%s-topsecret-key", testNamespace)

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()
	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("aws_key_id topsecretdata"))
	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring(
		fmt.Sprintf("aws_sec_key /fluentd/secret/%s", secretKey)))

	outputSecret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.OutputSecretName), outputSecret)()

	g.Expect(outputSecret.Data).Should(gomega.HaveKeyWithValue(secretKey, []byte("topsecretdata")))
}

// TODO add following tests:
// - resources from non watched namespaces are not incorporated
// - namespaced flow cannot use an output not enabled for the given namespace

func beforeEach(t *testing.T) func() {
	mgr, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	flowReconciler := &controllers.LoggingReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Flow"),
	}

	var wrappedReconciler reconcile.Reconciler
	wrappedReconciler, requests, _, reconcilerErrors = duplicateRequest(t, flowReconciler)

	err := controllers.SetupLoggingWithManager(mgr, ctrl.Log.WithName("manager").WithName("Setup")).Complete(wrappedReconciler)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped = startTestManager(t, mgr)

	return func() {
		close(stopMgr)
		mgrStopped.Wait()
	}
}

func ensureCreated(t *testing.T, object runtime.Object) func() {
	err := mgr.GetClient().Create(context.TODO(), object)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return func() {
		mgr.GetClient().Delete(context.TODO(), object)
	}
}

func ensureCreatedEventually(t *testing.T, ns, name string, object runtime.Object) func() {
	err := wait.Poll(time.Second, time.Second*3, func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), types.NamespacedName{
			Name: name, Namespace: ns,
		}, object)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	})
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	return func() {
		mgr.GetClient().Delete(context.TODO(), object)
	}
}
