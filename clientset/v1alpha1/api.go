package v1alpha1

import (
	"github.com/sysdiglabs/kube-apparmor-manager/api/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type AppArmorV1Alpha1Interface interface {
	AppArmorProfiles(namespace string) AppArmorProfileInterface
}

type AppArmorV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*AppArmorV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	//config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &AppArmorV1Alpha1Client{restClient: client}, nil
}

func (c *AppArmorV1Alpha1Client) ApparmorProfiles() AppArmorProfileInterface {
	return &appArmorProfileClient{
		restClient: c.restClient,
	}
}
