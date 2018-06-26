package main
import (
corev1 "k8s.io/api/core/v1"
extensionv1 "k8s.io/api/extensions/v1beta1"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/apimachinery/pkg/runtime/schema"
"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
)


func newFluentdRole() {

}

func newFluentdService() {

}

func newFluentdConfigmap() {

}

type fluentdDeploymentConfig  struct {
	Namespace string
	Replicas int32
}

// TODO the options should come from the operator configuration
func newFluentdDeployment(config *fluentdDeploymentConfig, cr *v1alpha1.LoggingOperator) *extensionv1.Deployment{
	labels := map[string]string {
		"app": "fluentd",
	}
	return &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "fluentd",
			Namespace: config.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "LoggingOperator",
				}),
			},
			Labels: labels,
		},
		Spec: extensionv1.DeploymentSpec{
			Replicas: &config.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fluentd",
					Namespace: config.Namespace,
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(cr, schema.GroupVersionKind{
							Group:   v1alpha1.SchemeGroupVersion.Group,
							Version: v1alpha1.SchemeGroupVersion.Version,
							Kind:    "LoggingOperator",
						}),
					},
					Labels: labels,
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path": "/metrics",
						"prometheus.io/port": "24231",
					},
				},
				Spec: corev1.PodSpec{
					Volumes:                       []corev1.Volume{
						{
							Name: "",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "",
									Type: nil,
								},
							},
						},
					},
					InitContainers:                nil,
					Containers:                    []corev1.Container{
						{
							Name:                     "",
							Image:                    "",
							Command:                  nil,
							Args:                     nil,
							WorkingDir:               "",
							Ports:                    nil,
							EnvFrom:                  nil,
							Env:                      nil,
							Resources:                corev1.ResourceRequirements{},
							VolumeMounts:             nil,
							VolumeDevices:            nil,
							LivenessProbe:            nil,
							ReadinessProbe:           nil,
							Lifecycle:                nil,
							TerminationMessagePath:   "",
							TerminationMessagePolicy: "",
							ImagePullPolicy:          "",
							SecurityContext:          nil,
							Stdin:                    false,
							StdinOnce:                false,
							TTY:                      false,
						}
					},
					RestartPolicy:                 "",
					TerminationGracePeriodSeconds: nil,
					ActiveDeadlineSeconds:         nil,
					DNSPolicy:                     "",
					NodeSelector:                  nil,
					ServiceAccountName:            "",
					DeprecatedServiceAccount:      "",
					AutomountServiceAccountToken:  nil,
					NodeName:                      "",
					HostNetwork:                   false,
					HostPID:                       false,
					HostIPC:                       false,
					SecurityContext: &corev1.PodSecurityContext{
						SELinuxOptions: &corev1.SELinuxOptions{
							User:  "",
							Role:  "",
							Type:  "",
							Level: "",
						},
						RunAsUser:          nil,
						RunAsNonRoot:       nil,
						SupplementalGroups: nil,
						FSGroup:            nil,
					},
					Hostname:         "",
					Subdomain:        "",
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: nil,
							},
							PreferredDuringSchedulingIgnoredDuringExecution: nil,
						},
						PodAffinity: &corev1.PodAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution:  nil,
							PreferredDuringSchedulingIgnoredDuringExecution: nil,
						},
						PodAntiAffinity: &corev1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution:  nil,
							PreferredDuringSchedulingIgnoredDuringExecution: nil,
						},
					},
					SchedulerName:     "",
					Tolerations:       nil,
					HostAliases:       nil,
					PriorityClassName: "",
					Priority:          nil,
					DNSConfig: &corev1.PodDNSConfig{
						Nameservers: nil,
						Searches:    nil,
						Options:     nil,
					},
				}
			},
		},
	}
}