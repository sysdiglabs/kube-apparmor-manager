package v1alpha1

import (
	"github.com/sysdiglabs/kube-apparmor-manager/api/types/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	apparmorProfiles = "apparmorprofiles"
)

type AppArmorProfileInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.AppArmorProfileList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.AppArmorProfile, error)
	Create(*v1alpha1.AppArmorProfile) (*v1alpha1.AppArmorProfile, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type appArmorProfileClient struct {
	restClient rest.Interface
}

func (c *appArmorProfileClient) List(opts metav1.ListOptions) (*v1alpha1.AppArmorProfileList, error) {
	result := v1alpha1.AppArmorProfileList{}
	err := c.restClient.
		Get().
		Resource(apparmorProfiles).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *appArmorProfileClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.AppArmorProfile, error) {
	result := v1alpha1.AppArmorProfile{}
	err := c.restClient.
		Get().
		Resource(apparmorProfiles).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *appArmorProfileClient) Create(profile *v1alpha1.AppArmorProfile) (*v1alpha1.AppArmorProfile, error) {
	result := v1alpha1.AppArmorProfile{}
	err := c.restClient.
		Post().
		Resource(apparmorProfiles).
		Body(profile).
		Do().
		Into(&result)

	return &result, err
}

func (c *appArmorProfileClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource(apparmorProfiles).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
