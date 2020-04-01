package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ApparmorProfileSpec struct {
	Rules    string `json:"rules"`
	Enforced bool   `json:"enforced"`
}

type ApparmorProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ApparmorProfileSpec `json:"spec"`
}

type ApparmorProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ApparmorProfile `json:"items"`
}
