package fluentbit

import (
	"bytes"
	"github.com/banzaicloud/logging-operator/cmd/logging-operator/sdkdecorator"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	extensionv1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"text/template"
)

var config *fluentBitDeploymentConfig

func initConfig() *fluentBitDeploymentConfig {
	if config == nil {
		config = &fluentBitDeploymentConfig{
			Name:      "fluent-bit",
			Namespace: viper.GetString("fluent-bit.namespace"),
			Labels:    map[string]string{"app": "fluent-bit"},
		}
	}
	return config
}

// InitFluentBit initialize fluent-bit
func InitFluentBit() {
	cfg := initConfig()
	if !checkIfDeamonSetExist(cfg) {
		logrus.Info("Deploying fluent-bit")
		if viper.GetBool("logging-operator.rbac") {
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newServiceAccount(cfg))
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newClusterRole(cfg))
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newClusterRoleBinding(cfg))
		}
		cfgMap, err := newFluentBitConfig(cfg)
		if err != nil {
			logrus.Error(err)
		}
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(cfgMap)
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(newFluentBitDaemonSet(cfg))
		logrus.Info("Fluent-bit deployed successfully")
	}
}

// DeleteFluentBit deletes fluent-bit if it exists
func DeleteFluentBit() {
	cfg := initConfig()
	if checkIfDeamonSetExist(cfg) {
		logrus.Info("Deleting fluent-bit")
		if viper.GetBool("logging-operator.rbac") {
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newServiceAccount(cfg))
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newClusterRole(cfg))
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newClusterRoleBinding(cfg))
		}
		cfgMap, err := newFluentBitConfig(cfg)
		if err != nil {
			logrus.Error(err)
		}
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(cfgMap)
		foregroundDeletion := metav1.DeletePropagationForeground
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentBitDaemonSet(cfg),
			sdk.WithDeleteOptions(&metav1.DeleteOptions{
				PropagationPolicy: &foregroundDeletion,
			}))
		logrus.Info("Fluent-bit deleted successfully")
	}
}

type fluentBitDeploymentConfig struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type fluentBitConfig struct {
	TLS struct {
		Enabled   bool
		SharedKey string
	}
	Monitor map[string]string
	Output  map[string]string
}

func newServiceAccount(cr *fluentBitDeploymentConfig) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "logging",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
	}
}

func newClusterRole(cr *fluentBitDeploymentConfig) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "LoggingRole",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs: []string{
					"get",
				},
				APIGroups: []string{""},
				Resources: []string{
					"pods",
				},
			},
		},
	}

}

func newClusterRoleBinding(cr *fluentBitDeploymentConfig) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "logging",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "logging",
				Namespace: cr.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "LoggingRole",
		},
	}
}

func generateConfig(input fluentBitConfig) (*string, error) {
	output := new(bytes.Buffer)
	text := viper.GetString("fluent-bit.config")

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

func newFluentBitConfig(cr *fluentBitDeploymentConfig) (*corev1.ConfigMap, error) {
	input := fluentBitConfig{
		TLS: struct {
			Enabled   bool
			SharedKey string
		}{
			Enabled:   viper.GetBool("fluent-bit.enabled"),
			SharedKey: "foobar",
		},
		Monitor: map[string]string{
			"Port": "2020",
		},
	}
	config, err := generateConfig(input)
	if err != nil {
		return nil, err
	}
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit-config",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},

		Data: map[string]string{
			"fluent-bit.conf": *config,
		},
	}
	return configMap, nil
}

func checkIfDeamonSetExist(cr *fluentBitDeploymentConfig) bool {
	fluentbitDaemonSet := &extensionv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Labels:    cr.Labels,
			Namespace: cr.Namespace,
		},
	}
	if err := sdk.Get(fluentbitDaemonSet); err != nil {
		logrus.Info("FluentBit DaemonSet does not exists!")
		logrus.Error(err)
		return false
	}
	logrus.Info("FluentBit DaemonSet already exists!")
	return true
}

// TODO in case of rbac add created serviceAccount name
func newFluentBitDaemonSet(cr *fluentBitDeploymentConfig) *extensionv1.DaemonSet {
	return &extensionv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Spec: extensionv1.DaemonSetSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.Labels,
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path":   "/metrics",
						"prometheus.io/port":   "2020",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "varlibcontainers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/docker/containers",
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "fluent-bit-config",
									},
								},
							},
						},
						{
							Name: "varlogs",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log",
								},
							},
						},
						{
							Name: "positions",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "fluent-tls",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "tls-for-logging-operator",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							// TODO move to configuration
							Name:  "fluent-bit",
							Image: "fluent/fluent-bit:latest",
							// TODO get from config translate to const
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 2020,
									Protocol:      "TCP",
								},
							},
							// TODO Get this from config
							Resources: corev1.ResourceRequirements{
								Limits:   nil,
								Requests: nil,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "varlibcontainers",
									ReadOnly:  true,
									MountPath: "/var/lib/docker/containers",
								},
								{
									Name:      "config",
									MountPath: "/fluent-bit/etc/fluent-bit.conf",
									SubPath:   "fluent-bit.conf",
								},
								{
									Name:      "positions",
									MountPath: "/tail-db",
								},
								{
									Name:      "varlogs",
									ReadOnly:  true,
									MountPath: "/var/log/",
								},
								{
									Name:      "fluent-tls",
									MountPath: "/fluent-bit/tls/caCert",
									SubPath:   "caCert",
								},
								{
									Name:      "fluent-tls",
									MountPath: "/fluent-bit/tls/clientCert",
									SubPath:   "clientCert",
								},
								{
									Name:      "fluent-tls",
									MountPath: "/fluent-bit/tls/clientKey",
									SubPath:   "clientKey",
								},
							},
						},
					},
				},
			},
		},
	}
}
