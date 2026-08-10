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

package fluentd

import (
	"context"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// newCheckPodReconciler builds a Reconciler that is just complete enough for
// newCheckPod, which only needs the logging resource and the fluentd spec.
// Pass a client only for the paths that talk to the API server, e.g. configCheck.
func newCheckPodReconciler(t *testing.T, fluentdSpec *v1beta1.FluentdSpec, c client.Client) *Reconciler {
	t.Helper()

	logging := &v1beta1.Logging{}
	logging.Name = "test"
	logging.Spec.ControlNamespace = "logging"
	logging.Spec.FluentdSpec = fluentdSpec
	require.NoError(t, logging.SetDefaults())

	config := "<system>\n</system>\n"

	return New(c, log.Log, logging, logging.Spec.FluentdSpec, nil, &config, nil, reconciler.ReconcilerOpts{}, nil)
}

// newFakeClient returns a client backed by an empty in-memory store for the
// configcheck paths that create and delete pods.
func newFakeClient(t *testing.T) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))

	return fake.NewClientBuilder().WithScheme(scheme).Build()
}

func TestNewCheckPodDNSSettings(t *testing.T) {
	dnsConfig := &corev1.PodDNSConfig{
		Nameservers: []string{"10.0.0.53"},
		Searches:    []string{"logging.svc.cluster.local"},
	}

	tests := []struct {
		name              string
		dnsPolicy         corev1.DNSPolicy
		dnsConfig         *corev1.PodDNSConfig
		expectedDNSPolicy corev1.DNSPolicy
	}{
		{
			// Left empty the API server applies ClusterFirst, which is what the
			// configcheck pod did unconditionally before dnsPolicy was propagated.
			name:              "Unset",
			expectedDNSPolicy: "",
		},
		{
			// The interesting case: with None the pod resolves *only* through
			// dnsConfig, so dropping the policy while keeping dnsConfig gives the
			// configcheck pod a different resolver than the aggregator it checks for.
			name:              "NoneWithDNSConfig",
			dnsPolicy:         corev1.DNSNone,
			dnsConfig:         dnsConfig,
			expectedDNSPolicy: corev1.DNSNone,
		},
		{
			name:              "ClusterFirstWithHostNet",
			dnsPolicy:         corev1.DNSClusterFirstWithHostNet,
			expectedDNSPolicy: corev1.DNSClusterFirstWithHostNet,
		},
		{
			name:              "Default",
			dnsPolicy:         corev1.DNSDefault,
			expectedDNSPolicy: corev1.DNSDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &v1beta1.FluentdSpec{
				DNSPolicy: tt.dnsPolicy,
				DNSConfig: tt.dnsConfig,
			}
			r := newCheckPodReconciler(t, spec, nil)

			pod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedDNSPolicy, pod.Spec.DNSPolicy)
			assert.Equal(t, tt.dnsConfig, pod.Spec.DNSConfig)
		})
	}
}

// TestNewCheckPodDNSSettingsMatchStatefulSet pins the configcheck pod and the
// aggregator to the same DNS settings, since a configcheck that resolves names
// differently from the aggregator can pass or fail for the wrong reason.
func TestNewCheckPodDNSSettingsMatchStatefulSet(t *testing.T) {
	ndots := "2"
	spec := &v1beta1.FluentdSpec{
		DNSPolicy: corev1.DNSNone,
		DNSConfig: &corev1.PodDNSConfig{
			Nameservers: []string{"10.0.0.53"},
			Options:     []corev1.PodDNSConfigOption{{Name: "ndots", Value: &ndots}},
		},
	}
	r := newCheckPodReconciler(t, spec, nil)

	checkPod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
	require.NoError(t, err)

	obj, _, err := r.statefulset()
	require.NoError(t, err)
	sts, ok := obj.(*appsv1.StatefulSet)
	require.True(t, ok)

	assert.Equal(t, sts.Spec.Template.Spec.DNSPolicy, checkPod.Spec.DNSPolicy)
	assert.Equal(t, sts.Spec.Template.Spec.DNSConfig, checkPod.Spec.DNSConfig)
}

// TestNewCheckPodDoesNotMutateSharedAffinity pins that clearing PodAntiAffinity
// on the check pod does not affect the StatefulSet - pod.Spec.Affinity starts
// out as the same pointer as fluentdSpec.Affinity, shared with statefulset(),
// so mutating it in place would silently strip anti-affinity from the
// aggregator too and let replicas co-schedule.
func TestNewCheckPodDoesNotMutateSharedAffinity(t *testing.T) {
	spec := &v1beta1.FluentdSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{TopologyKey: "kubernetes.io/hostname"},
				},
			},
		},
	}
	r := newCheckPodReconciler(t, spec, nil)

	checkPod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
	require.NoError(t, err)
	assert.Nil(t, checkPod.Spec.Affinity.PodAntiAffinity, "the check pod itself must not carry pod anti-affinity")

	obj, _, err := r.statefulset()
	require.NoError(t, err)
	sts, ok := obj.(*appsv1.StatefulSet)
	require.True(t, ok)
	assert.NotNil(t, sts.Spec.Template.Spec.Affinity.PodAntiAffinity, "clearing the check pod's anti-affinity must not strip it from the StatefulSet")
}

// TestNewCheckPodConfigCheckPodOverrides pins configCheckPod.{initContainers,
// volumes,activeDeadlineSeconds} to land on the generated configcheck pod, and
// the pod's own RestartPolicy to stay Never - ConfigCheckPodOverrides has no
// restartPolicy field, so it can't be changed by the merge, regardless of the
// (container-level, unrelated) restartPolicy on a merged init container.
func TestNewCheckPodConfigCheckPodOverrides(t *testing.T) {
	restartAlways := corev1.ContainerRestartPolicyAlways
	sidecar := corev1.Container{
		Name:          "geoip-refresh",
		Image:         "busybox:1.37",
		RestartPolicy: &restartAlways,
	}
	extraVolume := corev1.Volume{
		Name: "geoip-db",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	deadline := int64(120)
	spec := &v1beta1.FluentdSpec{
		ConfigCheckPod: &v1beta1.ConfigCheckPodOverrides{
			InitContainers:        []corev1.Container{sidecar},
			Volumes:               []corev1.Volume{extraVolume},
			ActiveDeadlineSeconds: &deadline,
		},
	}
	r := newCheckPodReconciler(t, spec, nil)

	checkPod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
	require.NoError(t, err)

	assert.Contains(t, checkPod.Spec.InitContainers, sidecar)
	assert.Contains(t, checkPod.Spec.Volumes, extraVolume)
	assert.Equal(t, &deadline, checkPod.Spec.ActiveDeadlineSeconds)
	assert.Equal(t, corev1.RestartPolicyNever, checkPod.Spec.RestartPolicy)
}

// TestNewCheckPodConfigCheckPodOverridesAreAdditive pins configCheckPod to be
// merged on top of the pod the operator would otherwise generate, not replace
// it - e.g. it would catch a regression to
// pod.Spec.InitContainers = overrides.InitContainers, which would silently
// drop tmp-dir-hack, config-reloader and the operator's own volumes.
func TestNewCheckPodConfigCheckPodOverridesAreAdditive(t *testing.T) {
	overrideInitContainer := corev1.Container{Name: "geoip-refresh", Image: "busybox:1.37"}
	overrideVolume := corev1.Volume{
		Name:         "geoip-db",
		VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
	}
	overrides := &v1beta1.ConfigCheckPodOverrides{
		InitContainers: []corev1.Container{overrideInitContainer},
		Volumes:        []corev1.Volume{overrideVolume},
	}

	tests := []struct {
		name               string
		compressConfigFile bool
		tlsEnabled         bool
		configCheckPod     *v1beta1.ConfigCheckPodOverrides
		expectNoOp         bool
	}{
		{name: "NilOverride", expectNoOp: true},
		{name: "EmptyOverride", configCheckPod: &v1beta1.ConfigCheckPodOverrides{}, expectNoOp: true},
		{name: "WithOverride", configCheckPod: overrides},
		{name: "WithOverrideCompressAndTLS", compressConfigFile: true, tlsEnabled: true, configCheckPod: overrides},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseSpec := &v1beta1.FluentdSpec{CompressConfigFile: tt.compressConfigFile}
			if tt.tlsEnabled {
				baseSpec.TLS = v1beta1.FluentdTLS{Enabled: true, SecretName: "fluentd-tls-secret"}
			}
			baseReconciler := newCheckPodReconciler(t, baseSpec, nil)
			basePod, err := baseReconciler.newCheckPod("deadbeef", *baseReconciler.fluentdSpec)
			require.NoError(t, err)

			spec := baseSpec.DeepCopy()
			spec.ConfigCheckPod = tt.configCheckPod
			r := newCheckPodReconciler(t, spec, nil)

			pod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
			require.NoError(t, err)

			if tt.expectNoOp {
				// Nothing to merge must mean nothing changed: this is the whole
				// backwards-compatibility claim of the configCheckPod field.
				assert.Equal(t, basePod.Spec, pod.Spec, "an empty configCheckPod must leave the pod spec untouched")
				return
			}

			for _, c := range basePod.Spec.InitContainers {
				assert.Contains(t, pod.Spec.InitContainers, c, "configCheckPod must not drop the operator's own init containers")
			}
			for _, v := range basePod.Spec.Volumes {
				assert.Contains(t, pod.Spec.Volumes, v, "configCheckPod must not drop the operator's own volumes")
			}
			assert.Contains(t, pod.Spec.InitContainers, overrideInitContainer)
			assert.Contains(t, pod.Spec.Volumes, overrideVolume)
		})
	}
}

// TestConfigCheckDeadlineExceededIsRetried pins that a configcheck pod failed
// by activeDeadlineSeconds (Status.Reason == "DeadlineExceeded") is deleted
// and reported not-ready for retry, rather than latched as an invalid config.
func TestConfigCheckFailedPodVerdict(t *testing.T) {
	tests := []struct {
		name        string
		reason      string
		expectRetry bool
	}{
		{name: "InvalidConfig", reason: ""},
		{name: "DeadlineExceeded", reason: "DeadlineExceeded", expectRetry: true},
		{name: "Evicted", reason: "Evicted", expectRetry: true},
		{name: "NodeAffinity", reason: "NodeAffinity", expectRetry: true},
		{name: "Shutdown", reason: "Shutdown", expectRetry: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := newFakeClient(t)
			r := newCheckPodReconciler(t, &v1beta1.FluentdSpec{}, fakeClient)
			ctx := context.Background()

			_, err := r.configCheck(ctx)
			require.NoError(t, err)

			pods := &corev1.PodList{}
			require.NoError(t, fakeClient.List(ctx, pods))
			require.Len(t, pods.Items, 1)

			pod := &pods.Items[0]
			pod.Status.Phase = corev1.PodFailed
			pod.Status.Reason = tt.reason
			require.NoError(t, fakeClient.Status().Update(ctx, pod))

			result, err := r.configCheck(ctx)
			require.NoError(t, err)

			remaining := &corev1.PodList{}
			require.NoError(t, fakeClient.List(ctx, remaining))

			if !tt.expectRetry {
				assert.True(t, result.Ready, "an invalid config is a final verdict, not something to retry")
				assert.False(t, result.Valid)
				assert.Len(t, remaining.Items, 1, "the pod that proved the config invalid must be kept for diagnosis")
				return
			}

			assert.False(t, result.Ready, "a pod killed around the check says nothing about the config")
			assert.False(t, result.Valid)
			assert.Contains(t, result.Message, tt.reason)
			assert.Empty(t, remaining.Items, "the pod must be deleted so a fresh check can be created")
		})
	}
}

// TestNewCheckPodExtraVolumeClaims pins which extraVolumes the check pod can mount as-is. A spec
// is only ever realized as a volumeClaimTemplate on the StatefulSet, so on a plain Pod it names a
// claim that never exists and must become an emptyDir; a bare claimName refers to a claim the user
// brought themselves, which does exist and must be left alone.
func TestNewCheckPodExtraVolumeClaims(t *testing.T) {
	pvcSpec := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Resources: corev1.VolumeResourceRequirements{
			Requests: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("1Gi")},
		},
	}

	tests := []struct {
		name            string
		claim           *volume.PersistentVolumeClaim
		expectEmptyDir  bool
		expectClaimName string
	}{
		{
			name: "TemplatedClaimBecomesEmptyDir",
			claim: &volume.PersistentVolumeClaim{
				PersistentVolumeClaimSpec: pvcSpec,
				PersistentVolumeSource:    corev1.PersistentVolumeClaimVolumeSource{ClaimName: "geoip-data"},
			},
			expectEmptyDir: true,
		},
		{
			name: "ExistingClaimIsKept",
			claim: &volume.PersistentVolumeClaim{
				PersistentVolumeSource: corev1.PersistentVolumeClaimVolumeSource{ClaimName: "geoip-data"},
			},
			expectClaimName: "geoip-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &v1beta1.FluentdSpec{
				ExtraVolumes: []v1beta1.ExtraVolume{
					{
						VolumeName:    "geoip-data",
						ContainerName: "fluentd",
						Path:          "/geoip",
						Volume:        &volume.KubernetesVolume{PersistentVolumeClaim: tt.claim},
					},
				},
			}
			r := newCheckPodReconciler(t, spec, nil)

			pod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
			require.NoError(t, err)

			var got *corev1.Volume
			for i := range pod.Spec.Volumes {
				if pod.Spec.Volumes[i].Name == "geoip-data" {
					got = &pod.Spec.Volumes[i]
				}
			}
			require.NotNil(t, got, "extraVolume must still be attached to the check pod")

			if tt.expectEmptyDir {
				assert.NotNil(t, got.EmptyDir)
				assert.Nil(t, got.PersistentVolumeClaim, "must not reference a claim that will never exist for the check pod")
			} else {
				require.NotNil(t, got.PersistentVolumeClaim, "a claim the user brought themselves must be mounted as-is")
				assert.Equal(t, tt.expectClaimName, got.PersistentVolumeClaim.ClaimName)
			}

			// The substitution must not lose the mount it exists for.
			assert.Contains(t, pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      "geoip-data",
				MountPath: "/geoip",
			})
		})
	}
}
