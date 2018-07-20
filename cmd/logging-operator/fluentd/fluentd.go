package fluentd

import (
	corev1 "k8s.io/api/core/v1"
	extensionv1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
    "github.com/operator-framework/operator-sdk/pkg/sdk"
    "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var labels = map[string]string{
	"app": "fluentd",
}

func InitFluentd() {
	if viper.GetBool("logging-operator.rbac") {
	}
    sdk.Create(newFluentdConfigmap())
    sdk.Create(newFluentdPVC())
    sdk.Create(newFluentdDeployment())
    sdk.Create(newFluentdService())
    // Create fluentd services
    // Possible options
    //  replica: x
    //  tag_rewrite config: ? it should be possible to give labels
    //  input port
    //  TLS?
    //  monitoring
    //    enabled:
    //    port:
    //    path:
}

func CheckIfDeploymentSetExist() bool {
	fluentdDeployment := &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Namespace: "default",
		},
	}
	if err := sdk.Get(fluentdDeployment); err != nil {
		logrus.Info("Fluentd Deployment does not exists!")
		return false
	}
	logrus.Info("Fluentd Deployment already exists!")
	return true
}

func newFluentdRole() {

}
func newFluentdService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd",
			Namespace: "default",
			Labels: map[string]string{
				"app": "fluentd",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       24240,
					TargetPort: intstr.IntOrString{IntVal: 24240},
				},
			},
			Selector: map[string]string{
				"app": "fluentd",
			},
			Type: "ClusterIP",
		},
	}

}

// Generate configmap into folder

// TODO This has to be a Golang template with proper values gathered
func newFluentdConfigmap() *corev1.ConfigMap {
	config := viper.GetString("fluentd.config")
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-config",
			Namespace: "default",
			Labels:    labels,
		},

		Data: map[string]string{
			"fluentd.conf": config,
		},
	}
	return configMap
}

type fluentdDeploymentConfig struct {
	Namespace string
	Replicas  int32
}

func newFluentdPVC() *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-buffer",
			Namespace: "default",
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse("10G"),
				},
			},
		},
	}
}

// /fluentd/etc/fluent.conf  is the default config
// TODO the options should come from the operator configuration
func newFluentdDeployment() *extensionv1.Deployment {
	labels := map[string]string{
		"app": "fluentd",
	}
	var replicas int32 = 1
	return &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd",
			Namespace: "default",
			Labels:    labels,
		},
		Spec: extensionv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "fluent-bit",
					Labels: labels,
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path":   "/metrics",
						"prometheus.io/port":   "25000",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "fluentd-config",
									},
								},
							},
						},
						{
							Name: "buffer",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "fluentd-buffer",
									ReadOnly:  false,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "fluentd",
							Image: "banzaicloud/fluentd:v1.1.2",
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 25000,
									Protocol:      "TCP",
								},
								{
									Name:          "fluent-input",
									ContainerPort: 24224,
									Protocol:      "TCP",
								},
							},

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/fluentd/etc/conf.d",
								},
								{
									Name:      "buffer",
									MountPath: "/buffers",
								},
							},
						},
					},
					//ServiceAccountName: "",
				},
			},
		},
	}
}
