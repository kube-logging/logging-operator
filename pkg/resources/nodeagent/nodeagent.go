// Copyright © 2021 Banzai Cloud
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

package nodeagent

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/resources"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentddataprovider"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/typeoverride"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	defaultServiceAccountName      = "fluentbit"
	clusterRoleBindingName         = "fluentbit"
	clusterRoleName                = "fluentbit"
	fluentBitSecretConfigName      = "fluentbit"
	fluentbitDaemonSetName         = "fluentbit"
	fluentbitPodSecurityPolicyName = "fluentbit"
	fluentbitServiceName           = "fluentbit"
	containerName                  = "fluent-bit"
)

func NodeAgentFluentbitDefaults(userDefined **v1beta1.NodeAgent) (*v1beta1.NodeAgent, error) {
	programDefault := &v1beta1.NodeAgent{
		FluentbitSpec: &v1beta1.NodeAgentFluentbit{
			DaemonSetOverrides: &typeoverride.DaemonSet{
				Spec: typeoverride.DaemonSetSpec{
					Template: typeoverride.PodTemplateSpec{
						ObjectMeta: typeoverride.ObjectMeta{
							Annotations: map[string]string{},
						},
						Spec: typeoverride.PodSpec{
							Containers: []v1.Container{
								{
									Name:            containerName,
									Image:           "fluent/fluent-bit:1.9.3",
									Command:         []string{"/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/conf_operator/fluent-bit.conf"},
									ImagePullPolicy: v1.PullIfNotPresent,
									Resources: v1.ResourceRequirements{
										Limits: v1.ResourceList{
											v1.ResourceMemory: resource.MustParse("100M"),
											v1.ResourceCPU:    resource.MustParse("200m"),
										},
										Requests: v1.ResourceList{
											v1.ResourceMemory: resource.MustParse("50M"),
											v1.ResourceCPU:    resource.MustParse("100m"),
										},
									},
								},
							},
						},
					},
				},
			},
			Flush:         1,
			Grace:         5,
			LogLevel:      "info",
			CoroStackSize: 24576,
			InputTail: v1beta1.InputTail{
				Path:            "/var/log/containers/*.log",
				RefreshInterval: "5",
				SkipLongLines:   "On",
				DB:              util.StringPointer("/tail-db/tail-containers-state.db"),
				MemBufLimit:     "5MB",
				Tag:             "kubernetes.*",
			},
			Security: &v1beta1.Security{
				RoleBasedAccessControlCreate: util.BoolPointer(true),
				SecurityContext:              &v1.SecurityContext{},
				PodSecurityContext:           &v1.PodSecurityContext{},
			},
			ContainersPath: "/var/lib/docker/containers",
			VarLogsPath:    "/var/log",
			BufferStorage: v1beta1.BufferStorage{
				StoragePath: "/buffers",
			},

			ForwardOptions: &v1beta1.ForwardOptions{
				RetryLimit: "False",
			},
		},
	}
	if (*userDefined).FluentbitSpec == nil {
		(*userDefined).FluentbitSpec = &v1beta1.NodeAgentFluentbit{}
	}

	if (*userDefined).FluentbitSpec.FilterAws != nil {

		programDefault.FluentbitSpec.FilterAws = &v1beta1.FilterAws{
			ImdsVersion:     "v2",
			AZ:              util.BoolPointer(true),
			Ec2InstanceID:   util.BoolPointer(true),
			Ec2InstanceType: util.BoolPointer(false),
			PrivateIP:       util.BoolPointer(false),
			AmiID:           util.BoolPointer(false),
			AccountID:       util.BoolPointer(false),
			Hostname:        util.BoolPointer(false),
			VpcID:           util.BoolPointer(false),
			Match:           "*",
		}

		err := merge.Merge(programDefault.FluentbitSpec.FilterAws, (*userDefined).FluentbitSpec.FilterAws)
		if err != nil {
			return nil, err
		}

	}
	if (*userDefined).FluentbitSpec.LivenessDefaultCheck == nil || *(*userDefined).FluentbitSpec.LivenessDefaultCheck {
		if (*userDefined).Profile != "windows" {
			programDefault.FluentbitSpec.Metrics = &v1beta1.Metrics{
				Port: 2020,
				Path: "/",
			}
		}
	}

	if (*userDefined).FluentbitSpec.Metrics != nil {

		programDefault.FluentbitSpec.Metrics = &v1beta1.Metrics{
			Interval: "15s",
			Timeout:  "5s",
			Port:     2020,
			Path:     "/api/v1/metrics/prometheus",
		}
		err := merge.Merge(programDefault.FluentbitSpec.Metrics, (*userDefined).FluentbitSpec.Metrics)
		if err != nil {
			return nil, err
		}
	}
	if programDefault.FluentbitSpec.Metrics != nil && (*userDefined).FluentbitSpec.Metrics != nil && (*userDefined).FluentbitSpec.Metrics.PrometheusAnnotations {
		defaultPrometheusAnnotations := &typeoverride.ObjectMeta{
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/path":   programDefault.FluentbitSpec.Metrics.Path,
				"prometheus.io/port":   fmt.Sprintf("%d", programDefault.FluentbitSpec.Metrics.Port),
			},
		}
		err := merge.Merge(&(programDefault.FluentbitSpec.DaemonSetOverrides.Spec.Template.ObjectMeta), defaultPrometheusAnnotations)
		if err != nil {
			return nil, err
		}
	}
	if programDefault.FluentbitSpec.Metrics != nil {
		defaultLivenessProbe := &v1.Probe{
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Path: programDefault.FluentbitSpec.Metrics.Path,
					Port: intstr.IntOrString{
						IntVal: programDefault.FluentbitSpec.Metrics.Port,
					},
				}},
			InitialDelaySeconds: 10,
			TimeoutSeconds:      0,
			PeriodSeconds:       10,
			SuccessThreshold:    0,
			FailureThreshold:    3,
		}
		if programDefault.FluentbitSpec.DaemonSetOverrides.Spec.Template.Spec.Containers[0].LivenessProbe == nil {
			programDefault.FluentbitSpec.DaemonSetOverrides.Spec.Template.Spec.Containers[0].LivenessProbe = &v1.Probe{}
		}

		err := merge.Merge(programDefault.FluentbitSpec.DaemonSetOverrides.Spec.Template.Spec.Containers[0].LivenessProbe, defaultLivenessProbe)
		if err != nil {
			return nil, err
		}
	}

	return programDefault, nil
}

var NodeAgentFluentbitWindowsDefaults = &v1beta1.NodeAgent{
	FluentbitSpec: &v1beta1.NodeAgentFluentbit{
		FilterKubernetes: v1beta1.FilterKubernetes{
			KubeURL:       "https://kubernetes.default.svc.cluster.local:443",
			KubeCAFile:    "c:\\var\\run\\secrets\\kubernetes.io\\serviceaccount\\ca.crt",
			KubeTokenFile: "c:\\var\\run\\secrets\\kubernetes.io\\serviceaccount\\token",
			KubeTagPrefix: "kubernetes.C.var.log.containers.",
		},
		InputTail: v1beta1.InputTail{
			Path: "C:\\var\\log\\containers\\*.log",
		},
		ContainersPath: "C:\\ProgramData\\docker",
		VarLogsPath:    "C:\\var\\log",
		DaemonSetOverrides: &typeoverride.DaemonSet{
			Spec: typeoverride.DaemonSetSpec{
				Template: typeoverride.PodTemplateSpec{
					Spec: typeoverride.PodSpec{
						Containers: []v1.Container{
							{
								Name:    containerName,
								Image:   "rancher/fluent-bit:1.6.10-rc7",
								Command: []string{"fluent-bit", "-c", "fluent-bit\\conf_operator\\fluent-bit.conf"},
								Resources: v1.ResourceRequirements{
									Limits: v1.ResourceList{
										v1.ResourceMemory: resource.MustParse("200M"),
										v1.ResourceCPU:    resource.MustParse("200m"),
									},
									Requests: v1.ResourceList{
										v1.ResourceMemory: resource.MustParse("100M"),
										v1.ResourceCPU:    resource.MustParse("100m"),
									},
								},
							}},
						NodeSelector: map[string]string{
							"kubernetes.io/os": "windows",
						},
					}},
			}},
	},
}
var NodeAgentFluentbitLinuxDefaults = &v1beta1.NodeAgent{
	FluentbitSpec: &v1beta1.NodeAgentFluentbit{},
}

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func (n *nodeAgentInstance) getFluentBitLabels() map[string]string {
	return util.MergeLabels(n.nodeAgent.Metadata.Labels, map[string]string{
		"app.kubernetes.io/name":     "fluentbit",
		"app.kubernetes.io/instance": n.nodeAgent.Name,
	}, generateLoggingRefLabels(n.logging.ObjectMeta.GetName()))
}

func (n *nodeAgentInstance) getServiceAccount() string {
	if n.nodeAgent.FluentbitSpec.Security.ServiceAccount != "" {
		return n.nodeAgent.FluentbitSpec.Security.ServiceAccount
	}
	return n.QualifiedName(defaultServiceAccountName)
}

//
//type DesiredObject struct {
//	Object runtime.Object
//	State  reconciler.DesiredState
//}
//
// Reconciler holds info what resource to reconcile
type Reconciler struct {
	Logging *v1beta1.Logging
	*reconciler.GenericResourceReconciler
	configs             map[string][]byte
	fluentdDataProvider fluentddataprovider.FluentdDataProvider
}

// NewReconciler creates a new NodeAgent reconciler
func New(client client.Client, logger logr.Logger, logging *v1beta1.Logging, opts reconciler.ReconcilerOpts, fluentdDataProvider fluentddataprovider.FluentdDataProvider) *Reconciler {
	return &Reconciler{
		Logging:                   logging,
		GenericResourceReconciler: reconciler.NewGenericReconciler(client, logger, opts),
		fluentdDataProvider:       fluentdDataProvider,
	}
}

type nodeAgentInstance struct {
	nodeAgent           *v1beta1.NodeAgent
	reconciler          *reconciler.GenericResourceReconciler
	logging             *v1beta1.Logging
	configs             map[string][]byte
	fluentdDataProvider fluentddataprovider.FluentdDataProvider
}

// Reconcile reconciles the NodeAgent resource
func (r *Reconciler) Reconcile() (*reconcile.Result, error) {
	for _, userDefinedAgent := range r.Logging.Spec.NodeAgents {
		var instance nodeAgentInstance
		NodeAgentFluentbitDefaults, err := NodeAgentFluentbitDefaults(&userDefinedAgent)
		if err != nil {
			return nil, err
		}

		switch userDefinedAgent.Profile {
		case "windows":
			err := merge.Merge(NodeAgentFluentbitDefaults, NodeAgentFluentbitWindowsDefaults)
			if err != nil {
				return nil, err
			}
		default:
			err := merge.Merge(NodeAgentFluentbitDefaults, NodeAgentFluentbitLinuxDefaults)
			if err != nil {
				return nil, err
			}

		}
		err = merge.Merge(NodeAgentFluentbitDefaults, &userDefinedAgent)
		if err != nil {
			return nil, err
		}

		instance = nodeAgentInstance{
			nodeAgent:           NodeAgentFluentbitDefaults,
			reconciler:          r.GenericResourceReconciler,
			logging:             r.Logging,
			fluentdDataProvider: r.fluentdDataProvider,
		}

		result, err := instance.Reconcile()
		if err != nil {
			return nil, errors.WrapWithDetails(err,
				"failed to reconcile instances", "NodeName", userDefinedAgent.Name)
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
}

// Reconcile reconciles the nodeAgentInstance resource
func (n *nodeAgentInstance) Reconcile() (*reconcile.Result, error) {
	for _, factory := range []resources.Resource{
		n.serviceAccount,
		n.clusterRole,
		n.clusterRoleBinding,
		n.clusterPodSecurityPolicy,
		n.pspClusterRole,
		n.pspClusterRoleBinding,
		n.configSecret,
		n.daemonSet,
		n.serviceMetrics,
		n.monitorServiceMetrics,
	} {
		o, state, err := factory()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", factory)
		}
		result, err := n.reconciler.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapWithDetails(err,
				"failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func RegisterWatches(builder *builder.Builder) *builder.Builder {
	return builder.
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{})
}

// nodeAgent QualifiedName
func (n *nodeAgentInstance) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s-%s", n.logging.Name, n.nodeAgent.Name, name)
}

// nodeAgent FluentdQualifiedName
func (n *nodeAgentInstance) FluentdQualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", n.logging.Name, name)
}
