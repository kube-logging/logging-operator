package fluentd

import (
	"bytes"
	"github.com/banzaicloud/logging-operator/cmd/logging-operator/sdkdecorator"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	extensionv1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"text/template"
)

type fluentdDeploymentConfig struct {
	Name      string
	Namespace string
	Replicas  int32
	Labels    map[string]string
}

type fluentdConfig struct {
	TLS struct {
		Enabled   bool
		SharedKey string
	}
}

var config *fluentdDeploymentConfig

func initConfig() *fluentdDeploymentConfig {
	if config == nil {
		config = &fluentdDeploymentConfig{
			Name:      "fluentd",
			Namespace: viper.GetString("fluentd.namespace"),
			Replicas:  1,
			Labels:    map[string]string{"app": "fluentd"},
		}
	}
	return config
}

// InitFluentd initialize fluentd
func InitFluentd() {
	fdc := initConfig()
	if !checkIfDeploymentExist(fdc) {
		if viper.GetBool("logging-operator.rbac") {
		}
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newFluentdConfigmap(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newFluentdAppConfigMap(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newFluentdPVC(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newFluentdDeployment(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newFluentdService(fdc))
		logrus.Info("Fluentd Deployment initialized!")
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

// DeleteFluentd deletes fluentd if exists
func DeleteFluentd() {
	fdc := initConfig()
	if checkIfDeploymentExist(fdc) {
		logrus.Info("Deleting fluentd")
		if viper.GetBool("logging-operator.rbac") {
		}
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentdConfigmap(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentdAppConfigMap(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentdPVC(fdc))
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentdService(fdc))
		foregroundDeletion := metav1.DeletePropagationForeground
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentdDeployment(fdc),
			sdk.WithDeleteOptions(&metav1.DeleteOptions{
				PropagationPolicy: &foregroundDeletion,
			}))
		logrus.Info("Fluentd Deployment deleted!")
	}
}

func checkIfDeploymentExist(fdc *fluentdDeploymentConfig) bool {
	fluentdDeployment := &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fdc.Name,
			Labels:    fdc.Labels,
			Namespace: fdc.Namespace,
		},
	}
	if err := sdk.Get(fluentdDeployment); err != nil {
		logrus.Info("Fluentd Deployment does not exists!")
		logrus.Error(err)
		return false
	}
	logrus.Info("Fluentd Deployment already exists!")
	return true
}

func newFluentdRole() {

}
func newFluentdService(fdc *fluentdDeploymentConfig) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fdc.Name,
			Namespace: fdc.Namespace,
			Labels:    fdc.Labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       24240,
					TargetPort: intstr.IntOrString{IntVal: 24240},
				},
			},
			Selector: fdc.Labels,
			Type:     "ClusterIP",
		},
	}

}

func generateConfig(input fluentdConfig) (*string, error) {
	output := new(bytes.Buffer)
	text := viper.GetString("fluentd.input")

	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return nil, err
	}
	outputString := output.String()
	return &outputString, nil
}

func newFluentdAppConfigMap(fdc *fluentdDeploymentConfig) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-app-config",
			Namespace: fdc.Namespace,
			Labels:    fdc.Labels,
		},
	}
}

func newFluentdConfigmap(fdc *fluentdDeploymentConfig) *corev1.ConfigMap {
	input := fluentdConfig{
		TLS: struct {
			Enabled   bool
			SharedKey string
		}{
			Enabled:   viper.GetBool("fluentd.tls_enabled"),
			SharedKey: "foobar",
		},
	}
	inputConfig, err := generateConfig(input)
	if err != nil {
		logrus.Errorf("Error during generating config %v", err)
		return nil
	}
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-config",
			Namespace: fdc.Namespace,
			Labels:    fdc.Labels,
		},

		Data: map[string]string{
			"fluentd.conf": viper.GetString("fluentd.default"),
			"input.conf":   *inputConfig,
			"devnull.conf": viper.GetString("fluentd.devnull"),
		},
	}
	return configMap
}

func newFluentdPVC(fdc *fluentdDeploymentConfig) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd-buffer",
			Namespace: fdc.Namespace,
			Labels:    fdc.Labels,
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

func newConfigMapReloader() *corev1.Container {
	return &corev1.Container{
		Name:  "config-reloader",
		Image: "jimmidyson/configmap-reload:v0.2.2",
		Args: []string{
			"-volume-dir=/fluentd/etc/conf.d",
			"-webhook-url=http://127.0.0.1:24444/api/config.reload",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "config",
				MountPath: "/fluentd/etc/conf.d",
			},
		},
	}
}

// TODO in case of rbac add created serviceAccount name
func newFluentdDeployment(fdc *fluentdDeploymentConfig) *extensionv1.Deployment {
	return &extensionv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fdc.Name,
			Namespace: fdc.Namespace,
			Labels:    fdc.Labels,
		},
		Spec: extensionv1.DeploymentSpec{
			Replicas: &fdc.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: fdc.Labels,
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
							Name: "app-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "fluentd-app-config",
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
						{
							Name: "fluentd-tls",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "tls-for-logging-operator",
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
									ContainerPort: 24240,
									Protocol:      "TCP",
								},
							},

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/fluentd/etc/fluent.conf",
									SubPath:   "fluentd.conf",
								},
								{
									Name:      "config",
									MountPath: "/fluentd/etc/input.conf",
									SubPath:   "input.conf",
								},
								{
									Name:      "config",
									MountPath: "/fluentd/etc/devnull.conf",
									SubPath:   "devnull.conf",
								},
								{
									Name:      "app-config",
									MountPath: "/fluentd/etc/app_config/",
								},
								{
									Name:      "buffer",
									MountPath: "/buffers",
								},
								{
									Name:      "fluentd-tls",
									MountPath: "/fluentd/etc/tls/caCert",
									SubPath:   "caCert",
								},
								{
									Name:      "fluentd-tls",
									MountPath: "/fluentd/etc/tls/serverCert",
									SubPath:   "serverCert",
								},
								{
									Name:      "fluentd-tls",
									MountPath: "/fluentd/etc/tls/serverKey",
									SubPath:   "serverKey",
								},
							},
						},
						*newConfigMapReloader(),
					},
					//ServiceAccountName: "",
				},
			},
		},
	}
}
