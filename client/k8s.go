package client

import (
	"flag"
	"path/filepath"
	"strings"
	"time"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"

	"github.com/sysdiglabs/kube-apparmor-manager/api/types/v1alpha1"
	aaClientset "github.com/sysdiglabs/kube-apparmor-manager/clientset/v1alpha1"
	"github.com/sysdiglabs/kube-apparmor-manager/types"
	"github.com/sysdiglabs/kube-apparmor-manager/utils"
	corev1 "k8s.io/api/core/v1"
	extClientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	cs        *kubernetes.Clientset
	aaclient  *aaClientset.AppArmorV1Alpha1Client
	extclient *extClientset.Clientset
}

// NewK8sClient return s Kubernetes client that contains the following
// - cs: general k8s client
// - aaclient: specific client to manage AppArmorProfile CRD object
// - extclient: extension client to manage CRD
func NewK8sClient() (*K8sClient, error) {
	var kubeconfig *string
	if home := utils.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	v1alpha1.AddToScheme(scheme.Scheme)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	aaClientset, err := aaClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	extClient, err := extClientset.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return &K8sClient{
		clientset,
		aaClientset,
		extClient,
	}, nil
}

// InstallCRD installs AppArmorProfile CRD
func (c *K8sClient) InstallCRD() error {
	klog.Infof("Creating a CRD: %s\n", v1alpha1.Name)

	crd := &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: v1alpha1.Name,
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: v1alpha1.GroupName,
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: v1alpha1.GroupVersion, Served: true, Storage: true},
			},
			Scope: apiextensions.ClusterScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     v1alpha1.Plural,
				Kind:       v1alpha1.Kind,
				ShortNames: []string{"aap"},
			},
		},
	}

	_, err := c.extclient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}

	klog.Infoln("The CRD created. Need to wait whether it is confirmed.")

	return c.waitForCRD()
}

// RemoveCRD removes AppArmorProfile CRD
func (c *K8sClient) RemoveCRD() error {
	c.extclient.ApiextensionsV1beta1().RESTClient().Delete().Name(v1alpha1.Name)
	return nil
}

func (c *K8sClient) waitForCRD() error {
	klog.Infof("Waiting for a CRD to be created: %s\n", v1alpha1.Name)

	err := wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
		// get CRDs by name
		crd, err := c.extclient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(v1alpha1.Name, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		for _, condition := range crd.Status.Conditions {
			if condition.Type == apiextensions.Established && condition.Status == apiextensions.ConditionTrue {
				// CRD successfully created.
				klog.Infoln("Confirmed that the CRD successfully created.")
				return true, err
			} else if condition.Type == apiextensions.NamesAccepted && condition.Status == apiextensions.ConditionFalse {
				klog.Fatalf("Name conflict while wait for CRD creation: %s, %v\n", condition.Reason, err)
			}
		}

		return false, err
	})
	if err != nil {
		return err
	}

	return nil
}

// GetNodes returns node list
func (c *K8sClient) GetNodes() (types.NodeList, error) {
	nodeList := types.NodeList{}

	list, err := c.cs.CoreV1().Nodes().List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, node := range list.Items {
		nodeReady := false
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				nodeReady = true
				break
			}
		}

		if nodeReady {
			n := types.NewNode()
			role := node.Labels[types.RoleLabel]
			n.Role = role
			n.NodeName = node.Name

			for _, addr := range node.Status.Addresses {
				switch addr.Type {
				case corev1.NodeExternalIP:
					n.ExternalIP = addr.Address
				case corev1.NodeInternalIP:
					n.InternalIP = addr.Address
				default:
				}
			}

			nodeList = append(nodeList, n)
		}
	}

	return nodeList, nil
}

// GetAppArmorProfiles returns apparmor profiles from etcd
func (c *K8sClient) GetAppArmorProfiles() ([]types.AppArmorProfile, error) {
	profileList := []types.AppArmorProfile{}
	list, err := c.aaclient.ApparmorProfiles().List(metav1.ListOptions{})

	if err != nil {
		return profileList, err
	}

	for _, p := range list.Items {
		var profile types.AppArmorProfile
		profile.Name = p.Name
		profile.Rules = p.Spec.Rules
		profile.Enforced = p.Spec.Enforced
		profileList = append(profileList, profile)
	}

	return profileList, nil
}
