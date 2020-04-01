package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Name = "apparmorprofiles.crd.security.sysdig.com"

	GroupName = "crd.security.sysdig.com"

	Kind = "AppArmorProfile"

	ListKind = "AppArmorProfileList"

	Plural = "apparmorprofiles"

	Singular = "apparmorprofile"

	GroupVersion = "v1alpha1"
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AppArmorProfile{},
		&AppArmorProfileList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
