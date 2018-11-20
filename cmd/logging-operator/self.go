package main

import (
	"errors"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrorCodes
var (
	ErrorNoParent    = errors.New("ErrorNoParent")
	ErrorUnknownKind = errors.New("ErrorUnknownKind")
)

// GetSelf returning the running Pod
func GetSelf(name, namespace string) (*corev1.Pod, error) {
	podObject := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := sdk.Get(podObject)
	if err != nil {
		return nil, err
	}
	return podObject, nil
}

// GetDeployment of a running Pod
func GetDeployment(pod *corev1.Pod, namespace string) (metav1.Object, error) {
	rs, err := GetParent(pod, namespace)
	if err != nil {
		return nil, err
	}
	deployment, err := GetParent(rs, namespace)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

// GetParent return an object parent
func GetParent(obj metav1.Object, namespace string) (metav1.Object, error) {
	parent := metav1.GetControllerOf(obj)
	if parent != nil {
		return GetObjectFromOwnerReference(parent, namespace)
	}
	return nil, ErrorNoParent
}

// GetObjectFromOwnerReference get parent from OwnerReference
func GetObjectFromOwnerReference(owner *metav1.OwnerReference, namespace string) (metav1.Object, error) {
	switch owner.Kind {
	case "Pod":
		obj := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       owner.Kind,
				APIVersion: owner.APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      owner.Name,
				Namespace: namespace,
			},
		}
		err := sdk.Get(obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	case "ReplicaSet":
		obj := &v1beta1.ReplicaSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       owner.Kind,
				APIVersion: owner.APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      owner.Name,
				Namespace: namespace,
			},
		}
		err := sdk.Get(obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	case "Deployment":
		obj := &v1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       owner.Kind,
				APIVersion: owner.APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      owner.Name,
				Namespace: namespace,
			},
		}
		err := sdk.Get(obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return nil, ErrorUnknownKind
}
