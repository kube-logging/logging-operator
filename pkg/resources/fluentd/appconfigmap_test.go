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
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// newCheckPodReconciler builds a Reconciler that is just complete enough for
// newCheckPod, which only needs the logging resource and the fluentd spec.
func newCheckPodReconciler(t *testing.T, fluentdSpec *v1beta1.FluentdSpec) *Reconciler {
	t.Helper()

	logging := &v1beta1.Logging{}
	logging.Name = "test"
	logging.Spec.ControlNamespace = "logging"
	logging.Spec.FluentdSpec = fluentdSpec
	require.NoError(t, logging.SetDefaults())

	config := "<system>\n</system>\n"

	return New(nil, log.Log, logging, logging.Spec.FluentdSpec, nil, &config, nil, reconciler.ReconcilerOpts{})
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
			r := newCheckPodReconciler(t, spec)

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
	r := newCheckPodReconciler(t, spec)

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
	r := newCheckPodReconciler(t, spec)

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
	r := newCheckPodReconciler(t, spec)

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
	}{
		{name: "NilOverride"},
		{name: "WithOverride", configCheckPod: overrides},
		{name: "WithOverrideAndCompress", compressConfigFile: true, configCheckPod: overrides},
		{name: "WithOverrideAndTLS", tlsEnabled: true, configCheckPod: overrides},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseSpec := &v1beta1.FluentdSpec{CompressConfigFile: tt.compressConfigFile}
			if tt.tlsEnabled {
				baseSpec.TLS = v1beta1.FluentdTLS{Enabled: true, SecretName: "fluentd-tls-secret"}
			}
			baseReconciler := newCheckPodReconciler(t, baseSpec)
			basePod, err := baseReconciler.newCheckPod("deadbeef", *baseReconciler.fluentdSpec)
			require.NoError(t, err)

			spec := baseSpec.DeepCopy()
			spec.ConfigCheckPod = tt.configCheckPod
			r := newCheckPodReconciler(t, spec)

			pod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
			require.NoError(t, err)

			for _, c := range basePod.Spec.InitContainers {
				assert.Contains(t, pod.Spec.InitContainers, c, "configCheckPod must not drop the operator's own init containers")
			}
			for _, v := range basePod.Spec.Volumes {
				assert.Contains(t, pod.Spec.Volumes, v, "configCheckPod must not drop the operator's own volumes")
			}

			if tt.configCheckPod != nil {
				assert.Contains(t, pod.Spec.InitContainers, overrideInitContainer)
				assert.Contains(t, pod.Spec.Volumes, overrideVolume)
			}
		})
	}
}

// TestConfigCheckDeadlineExceededIsRetried pins that a configcheck pod failed
// by activeDeadlineSeconds (Status.Reason == "DeadlineExceeded") is deleted
// and reported not-ready for retry, rather than latched as an invalid config.
func TestConfigCheckDeadlineExceededIsRetried(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))

	logging := &v1beta1.Logging{}
	logging.Name = "test"
	logging.Spec.ControlNamespace = "logging"
	logging.Spec.FluentdSpec = &v1beta1.FluentdSpec{}
	require.NoError(t, logging.SetDefaults())

	config := "<system>\n</system>\n"
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	r := New(fakeClient, log.Log, logging, logging.Spec.FluentdSpec, nil, &config, nil, reconciler.ReconcilerOpts{})
	ctx := context.Background()

	// first pass creates the configcheck pod
	_, err := r.configCheck(ctx)
	require.NoError(t, err)

	pods := &corev1.PodList{}
	require.NoError(t, fakeClient.List(ctx, pods))
	require.Len(t, pods.Items, 1)
	pod := &pods.Items[0]

	pod.Status.Phase = corev1.PodFailed
	pod.Status.Reason = "DeadlineExceeded"
	require.NoError(t, fakeClient.Status().Update(ctx, pod))

	result, err := r.configCheck(ctx)
	require.NoError(t, err)
	assert.False(t, result.Ready)
	assert.False(t, result.Valid)
	assert.Contains(t, result.Message, "DeadlineExceeded")

	remaining := &corev1.PodList{}
	require.NoError(t, fakeClient.List(ctx, remaining))
	assert.Empty(t, remaining.Items, "the pod that exceeded its deadline should have been deleted so a fresh one can be created")
}

// TestNewCheckPodPVCBackedExtraVolumeBecomesEmptyDir pins that a PVC-backed
// extraVolume (one whose PersistentVolumeClaimSpec is set, meaning it is only
// ever realized via volumeClaimTemplates on the StatefulSet) is downgraded to
// an emptyDir on the check pod, since the check pod is a plain Pod with no
// matching PVC - referencing the claim name outright would make the pod
// permanently uncreatable.
func TestNewCheckPodPVCBackedExtraVolumeBecomesEmptyDir(t *testing.T) {
	spec := &v1beta1.FluentdSpec{
		ExtraVolumes: []v1beta1.ExtraVolume{
			{
				VolumeName:    "geoip-data",
				ContainerName: "fluentd",
				Path:          "/geoip",
				Volume: &volume.KubernetesVolume{
					PersistentVolumeClaim: &volume.PersistentVolumeClaim{
						PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
							Resources: corev1.VolumeResourceRequirements{
								Requests: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("1Gi")},
							},
						},
						PersistentVolumeSource: corev1.PersistentVolumeClaimVolumeSource{ClaimName: "geoip-data"},
					},
				},
			},
		},
	}
	r := newCheckPodReconciler(t, spec)

	pod, err := r.newCheckPod("deadbeef", *r.fluentdSpec)
	require.NoError(t, err)

	var got *corev1.Volume
	for i := range pod.Spec.Volumes {
		if pod.Spec.Volumes[i].Name == "geoip-data" {
			got = &pod.Spec.Volumes[i]
		}
	}
	require.NotNil(t, got, "extraVolume must still be attached to the check pod")
	assert.NotNil(t, got.EmptyDir, "a PVC meant for volumeClaimTemplates has no matching PVC on a plain Pod, so the check pod must fall back to an emptyDir")
	assert.Nil(t, got.PersistentVolumeClaim, "must not reference a claim name that will never exist for the check pod")
}
