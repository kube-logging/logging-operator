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
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// dryRunner reports whether the API server would accept a Service pinned to the given family.
type dryRunner func(ctx context.Context, namespace string, family corev1.IPFamily) error

// Detector reports which Service IP families the cluster can allocate. Naming a family the cluster
// has no range for is rejected at create with no recovery path, so the answer has to come from the
// API server rather than from the user's intent.
type Detector struct {
	dryRun dryRunner
	log    logr.Logger

	mu       sync.Mutex
	resolved bool
	families []corev1.IPFamily
}

func NewDetector(c client.Client, log logr.Logger) *Detector {
	return &Detector{dryRun: dryRunService(c), log: log}
}

// Families lists the allocatable families, IPv6 first, and nil for a plain IPv4 cluster so the
// caller leaves the choice to the API server. An error means the cluster could not be asked: the
// caller must not carry on, because nil would then be indistinguishable from a real IPv4-only
// answer and the Service and the listener would disagree about which family to use.
func (d *Detector) Families(ctx context.Context, namespace string) ([]corev1.IPFamily, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.resolved {
		return d.families, nil
	}

	hasIPv6, err := d.allocatable(ctx, namespace, corev1.IPv6Protocol)
	if err != nil {
		return nil, err
	}
	hasIPv4, err := d.allocatable(ctx, namespace, corev1.IPv4Protocol)
	if err != nil {
		return nil, err
	}

	// No cluster serves neither family, so something rejects every Service and the answer is not
	// ours to cache. An admission webhook doing that would otherwise pin nil until a restart.
	if !hasIPv6 && !hasIPv4 {
		return nil, errors.New("every service IP family probe was rejected, cannot determine the cluster families")
	}

	if !hasIPv6 {
		// Re-probed every time, so a range added to a running cluster is picked up without a restart.
		d.log.V(1).Info("cluster has no IPv6 service range, leaving the IP families to the cluster")
		return nil, nil
	}

	d.families = []corev1.IPFamily{corev1.IPv6Protocol}
	if hasIPv4 {
		d.families = append(d.families, corev1.IPv4Protocol)
	}
	d.resolved = true
	return d.families, nil
}

func (d *Detector) allocatable(ctx context.Context, namespace string, family corev1.IPFamily) (bool, error) {
	err := d.dryRun(ctx, namespace, family)
	if err == nil {
		return true, nil
	}
	if apierrors.IsInvalid(err) {
		return false, nil
	}
	return false, fmt.Errorf("probing the %s service range: %w", family, err)
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
