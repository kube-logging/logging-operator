// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package eventtailer

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ServiceAccount resource for reconciler
func (e *EventTailer) ServiceAccount() (runtime.Object, reconciler.DesiredState, error) {
	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: e.objectMeta(),
	}
	return &serviceAccount, reconciler.StatePresent, nil
}
