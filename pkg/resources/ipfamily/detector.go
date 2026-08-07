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
	"sync"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// dryRunner reports whether the API server would accept a Service pinned to the given family.
type dryRunner func(ctx context.Context, namespace string, family corev1.IPFamily) error

// Detector answers whether the cluster can allocate IPv6 Service addresses. Naming a family the
// cluster has no range for is rejected at create with no recovery path, so the answer has to come
// from the API server rather than from the user's intent.
type Detector struct {
	dryRun dryRunner
	log    logr.Logger

	mu       sync.Mutex
	resolved bool
	ipv6     bool
}

func NewDetector(c client.Client, log logr.Logger) *Detector {
	return &Detector{dryRun: dryRunService(c), log: log}
}

// SupportsIPv6 defaults to false while the answer is unknown: a needless IPv4 primary costs a
// metrics scrape, whereas an IPv6 primary the cluster cannot allocate is unrecoverable.
func (d *Detector) SupportsIPv6(ctx context.Context, namespace string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.resolved {
		return d.ipv6
	}

	err := d.dryRun(ctx, namespace, corev1.IPv6Protocol)
	if err == nil {
		d.resolved, d.ipv6 = true, true
		return true
	}
	if !apierrors.IsInvalid(err) {
		d.log.Error(err, "could not determine cluster IPv6 support, assuming none for now")
		return false
	}

	// An admission webhook rejecting every Service looks identical to a missing IPv6 range.
	if controlErr := d.dryRun(ctx, namespace, corev1.IPv4Protocol); controlErr != nil {
		d.log.Error(controlErr, "cluster rejects the IP family probe outright, assuming no IPv6 for now")
		return false
	}

	d.log.Info("cluster has no IPv6 service range, leaving the IP families to the cluster")
	d.resolved, d.ipv6 = true, false
	return false
}

func dryRunService(c client.Client) dryRunner {
	return func(ctx context.Context, namespace string, family corev1.IPFamily) error {
		probe := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "logging-ipfamily-probe-",
				Namespace:    namespace,
			},
			Spec: corev1.ServiceSpec{
				Ports:      []corev1.ServicePort{{Port: 1}},
				IPFamilies: []corev1.IPFamily{family},
			},
		}
		return c.Create(ctx, probe, client.DryRunAll)
	}
}
