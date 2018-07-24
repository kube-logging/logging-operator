package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoggingOperatorList auto generated by the sdk
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LoggingOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []LoggingOperator `json:"items"`
}

// LoggingOperator auto generated by the sdk
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LoggingOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              LoggingOperatorSpec   `json:"spec"`
	Status            LoggingOperatorStatus `json:"status,omitempty"`
}
// LoggingOperatorSpec holds the spec for the operator
type LoggingOperatorSpec struct {
	// Fill me
}
// LoggingOperatorStatus holds the status info for the operator
type LoggingOperatorStatus struct {
	// Fill me
}
