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

package ipfamily

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestSupportsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		results  map[corev1.IPFamily]error
		expected bool
		resolved bool
	}{
		{
			name:     "dual-stack cluster accepts the IPv6 probe",
			results:  map[corev1.IPFamily]error{corev1.IPv6Protocol: nil},
			expected: true,
			resolved: true,
		},
		{
			name: "single-stack cluster rejects IPv6 but accepts the IPv4 control",
			results: map[corev1.IPFamily]error{
				corev1.IPv6Protocol: invalidFamilyErr(),
				corev1.IPv4Protocol: nil,
			},
			expected: false,
			resolved: true,
		},
		{
			name: "a webhook rejecting every service leaves the answer unknown",
			results: map[corev1.IPFamily]error{
				corev1.IPv6Protocol: invalidFamilyErr(),
				corev1.IPv4Protocol: invalidFamilyErr(),
			},
			expected: false,
			resolved: false,
		},
		{
			name:     "a denied probe leaves the answer unknown",
			results:  map[corev1.IPFamily]error{corev1.IPv6Protocol: apierrors.NewForbidden(schema.GroupResource{Resource: "services"}, "probe", errors.New("nope"))},
			expected: false,
			resolved: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			d := &Detector{
				log: logr.Discard(),
				dryRun: func(_ context.Context, _ string, family corev1.IPFamily) error {
					calls++
					return test.results[family]
				},
			}

			require.Equal(t, test.expected, d.SupportsIPv6(context.Background(), "default"))

			before := calls
			require.Equal(t, test.expected, d.SupportsIPv6(context.Background(), "default"))
			if test.resolved {
				require.Equal(t, before, calls, "a resolved answer must not re-probe")
			} else {
				require.Greater(t, calls, before, "an unresolved answer must be retried")
			}
		})
	}
}

// A cluster that gains an IPv6 range must be picked up without restarting the operator.
func TestSupportsIPv6RecoversAfterATransientFailure(t *testing.T) {
	failing := true
	d := &Detector{
		log: logr.Discard(),
		dryRun: func(_ context.Context, _ string, _ corev1.IPFamily) error {
			if failing {
				return apierrors.NewServiceUnavailable("apiserver is having a moment")
			}
			return nil
		},
	}

	require.False(t, d.SupportsIPv6(context.Background(), "default"))

	failing = false
	require.True(t, d.SupportsIPv6(context.Background(), "default"))
}

func invalidFamilyErr() error {
	return apierrors.NewInvalid(
		schema.GroupKind{Kind: "Service"},
		"probe",
		field.ErrorList{field.Invalid(
			field.NewPath("spec", "ipFamilies").Index(0),
			"IPv6",
			"not configured on this cluster",
		)},
	)
}
