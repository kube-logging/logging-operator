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

type FluentdDeployment struct {
	Namespace string
	Labels    map[string]string
}

func InitFluentd(fluentd *FluentdDeployment) {
	if !fluentd.checkIfDeploymentExist() {
		if viper.GetBool("logging-operator.rbac") {
		}
		sdk.Create(fluentd.newFluentdConfigmap())
		sdk.Create(fluentd.newFluentdPVC())
		sdk.Create(fluentd.newFluentdDeployment())
		sdk.Create(fluentd.newFluentdService())
	}
    // Create fluentd services
    // Possible options
    //  replica: x
    //  tag_rewrite config: ? it should be possible to give Labels
    //  input port
    //  TLS?
    //  monitoring
    //    enabled:
    //    port:
    //    path:
}

func (d *FluentdDeployment)checkIfDeploymentExist() bool {
	fluentdDeployment := &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    d.Labels,
			Namespace: d.Namespace,
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
func (d *FluentdDeployment)newFluentdService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd",
			Namespace: d.Namespace,
			Labels: d.Labels,
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

// TODO This has to be a Golang template with proper values gathered
func (d *FluentdDeployment)newFluentdConfigmap() *corev1.ConfigMap {
	config := viper.GetString("fluentd.config")
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-config",
			Namespace: d.Namespace,
			Labels:    d.Labels,
		},

		Data: map[string]string{
			"fluentd.conf": config,
		},
	}
	return configMap
}

func (d *FluentdDeployment)newFluentdPVC() *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-buffer",
			Namespace: d.Namespace,
			Labels:    d.Labels,
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

// TODO in case of rbac add created serviceAccount name
func (d *FluentdDeployment)newFluentdDeployment() *extensionv1.Deployment {
	var replicas int32 = 1
	return &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd",
			Namespace: d.Namespace,
			Labels:    d.Labels,
		},
		Spec: extensionv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "fluentd",
					Labels: d.Labels,
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
