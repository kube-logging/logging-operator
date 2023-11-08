// Copyright © 2020 Banzai Cloud
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

package config

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func CheckConfigForOutput(t *testing.T, output v1beta1.SyslogNGOutput, expected string, opts ...OutputConfigCheckOption) {
	options := OutputConfigCheckOptions{
		IndentWith: "    ",
	}
	for _, opt := range opts {
		opt(&options)
	}
	if options.SecretLoaderFactory == nil {
		options.SecretLoaderFactory = &TestSecretLoaderFactory{}
	}
	renderer := renderOutput(output, options.SecretLoaderFactory)
	result := &strings.Builder{}
	err := renderer(render.RenderContext{
		Out:        result,
		IndentWith: options.IndentWith,
	})
	CheckError(t, options.ExpectedError, err)
	assert.Equal(t, strings.TrimSpace(Untab(expected))+"\n", result.String())
}

func CheckConfigForClusterOutput(t *testing.T, output v1beta1.SyslogNGClusterOutput, expected string, opts ...OutputConfigCheckOption) {
	options := OutputConfigCheckOptions{
		IndentWith: "    ",
	}
	for _, opt := range opts {
		opt(&options)
	}
	if options.SecretLoaderFactory == nil {
		options.SecretLoaderFactory = &TestSecretLoaderFactory{}
	}
	renderer := renderClusterOutput(output, options.SecretLoaderFactory)
	result := &strings.Builder{}
	err := renderer(render.RenderContext{
		Out:        result,
		IndentWith: options.IndentWith,
	})
	CheckError(t, options.ExpectedError, err)
	assert.Equal(t, strings.TrimSpace(Untab(expected))+"\n", result.String())
}

type OutputConfigCheckOptions struct {
	ExpectedError       interface{}
	IndentWith          string
	SecretLoaderFactory SecretLoaderFactory
}

type OutputConfigCheckOption func(options *OutputConfigCheckOptions)

func CheckError(t *testing.T, expected interface{}, actual error, msgAndArgs ...interface{}) {
	t.Helper()
	switch expected := expected.(type) {
	case nil:
		require.NoError(t, actual, msgAndArgs...)
	case bool:
		if expected {
			require.Error(t, actual, msgAndArgs...)
		} else {
			require.NoError(t, actual, msgAndArgs...)
		}
	case func(error) bool:
		require.True(t, expected(actual), msgAndArgs...)
	default:
		require.Equal(t, expected, actual, msgAndArgs...)
	}
}

var leadingTabs = regexp.MustCompile("(?m:^\t+)")

func Untab(s string) string {
	return leadingTabs.ReplaceAllStringFunc(s, func(match string) string {
		return strings.Repeat("    ", len(match))
	})
}

type TestSecretLoaderFactory struct {
	Reader    client.Reader
	MountPath string
	Secrets   secret.MountSecrets
}

func (f *TestSecretLoaderFactory) SecretLoaderForNamespace(ns string) secret.SecretLoader {
	return secret.NewSecretLoader(f.Reader, ns, f.MountPath, &f.Secrets)
}

type SecretReader struct {
	Secrets []corev1.Secret
}

func (r SecretReader) Get(_ context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if secret, ok := obj.(*corev1.Secret); ok {
		if secret == nil {
			return nil
		}
		for _, s := range r.Secrets {
			if s.Namespace == key.Namespace && s.Name == key.Name {
				*secret = s
				return nil
			}
		}
		return apierrors.NewNotFound(corev1.Resource("secret"), key.String())
	}
	return apierrors.NewNotFound(schema.GroupResource{
		Group:    obj.GetObjectKind().GroupVersionKind().Group,
		Resource: strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind),
	}, key.String())
}

func (r SecretReader) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	panic("not implemented")
}

var _ client.Reader = (*SecretReader)(nil)

func NewTrue() *bool {
	b := true
	return &b
}

func NewFalse() *bool {
	b := false
	return &b
}
