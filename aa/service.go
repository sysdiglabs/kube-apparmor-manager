package aa

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"k8s.io/klog"

	"github.com/sysdiglabs/kube-apparmor-manager/aa/commands"
	"github.com/sysdiglabs/kube-apparmor-manager/client"
	"github.com/sysdiglabs/kube-apparmor-manager/types"
	"github.com/sysdiglabs/kube-apparmor-manager/utils"
)

const (
	envSSHUsername   = "SSH_USERNAME"
	envSSHPERMFile   = "SSH_PERM_FILE"
	envSSHPassPhrase = "SSH_PASSPHRASE"

	Enabled  = "Enabled"
	Disabled = "Disabled"
	Unknown  = "Unknown"

	SSH_PORT = "22"
)

type AppArmor struct {
	k8sClient *client.K8sClient
	sshClient *client.SSHClient
}

// NewAppArmor returns a new AppArmor object
func NewAppArmor() (*AppArmor, error) {
	k8s, err := client.NewK8sClient()

	if err != nil {
		return nil, err
	}

	username := os.Getenv(envSSHUsername)

	if username == "" {
		username = "admin"
	}

	sshPermFile := os.Getenv(envSSHPERMFile)

	if sshPermFile == "" {
		sshPermFile = fmt.Sprintf("%s/.ssh/id_rsa", utils.HomeDir())
	}

	sshPassPhrase := os.Getenv(envSSHPassPhrase)

	ssh, err := client.NewSSHClientConfig(username, sshPermFile, sshPassPhrase)

	if err != nil {
		return nil, err
	}
	return &AppArmor{
		k8sClient: k8s,
		sshClient: ssh,
	}, nil
}

// InstallCRD installs CRD in Kubernetes
func (aa *AppArmor) InstallCRD() error {
	return aa.k8sClient.InstallCRD()
}

// InstallAppArmor installs AppArmor service on worker nodes
func (aa *AppArmor) InstallAppArmor() error {
	nodes, err := aa.k8sClient.GetNodes()

	if err != nil {
		return err
	}

	for _, node := range nodes {
		err = aa.install(node)

		if err != nil {
			return err
		}
	}
	return nil
}

func (aa *AppArmor) install(node *types.Node) error {
	if node.IsMaster() {
		return nil
	}

	err := aa.sshClient.Connect(node.ExternalIP, SSH_PORT)

	if err != nil {
		return err
	}

	defer aa.sshClient.Close()

	if aa.enabledInConnection(node) {
		klog.Infof("AppArmor was enabled on node: %s (external IP: %s)", node.NodeName, node.ExternalIP)
		return nil
	}

	err = aa.sshClient.ExecuteBatch(commands.InstallAppArmor, true)

	if err != nil {
		return err
	}

	return nil
}

// Sync syncs AppArmor profiles from etcd to worker nodes
func (aa *AppArmor) Sync() error {
	nodes, err := aa.k8sClient.GetNodes()

	if err != nil {
		return err
	}

	profiles, err := aa.k8sClient.GetAppArmorProfiles()

	if err != nil {
		return err
	}

	for _, node := range nodes {
		for _, profile := range profiles {
			err := aa.syncProfile(node, profile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (aa *AppArmor) syncProfile(node *types.Node, profile types.AppArmorProfile) error {
	if node.IsMaster() {
		return nil
	}

	err := aa.sshClient.Connect(node.ExternalIP, SSH_PORT)
	if err != nil {
		return err
	}

	defer aa.sshClient.Close()

	if !aa.enabledInConnection(node) {
		klog.Infof("AppArmor was not enabled on node: %s (external IP: %s), no sync happen.", node.NodeName, node.ExternalIP)
		return nil
	}

	err = aa.sshClient.ExecuteBatch(commands.CreateProfileCommands(profile), true)

	if err != nil {
		return err
	}

	if profile.Enforced {
		err = aa.sshClient.ExecuteBatch(commands.EnforceProfileCommands(profile), true)
	} else {
		// turn it into complain mode
		err = aa.sshClient.ExecuteBatch(commands.ComplainProfileCommands(profile), true)
	}

	if err != nil {
		return err
	}

	return nil
}

// AppArmorEnabled get AppArmor enabled status on worker nodes
func (aa *AppArmor) AppArmorEnabled() (types.NodeList, error) {
	nodes, err := aa.k8sClient.GetNodes()

	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		_, err := aa.enabled(node)
		if err != nil {
			return nodes, err
		}
	}

	return nodes, nil
}

func (aa *AppArmor) enabled(node *types.Node) (bool, error) {
	if node.IsMaster() {
		return false, nil
	}

	err := aa.sshClient.Connect(node.ExternalIP, SSH_PORT)
	if err != nil {
		return false, err
	}

	defer aa.sshClient.Close()

	return aa.enabledInConnection(node), nil
}

func (aa *AppArmor) enabledInConnection(node *types.Node) bool {
	stdout, stderr, err := aa.sshClient.ExecuteOne(commands.AAEnable, true)

	if err != nil {
		return false
	}

	if len(stderr) > 0 {
		return false
	}

	if strings.ToLower(stdout) == "yes" {
		node.AppArmorEnabled = true
		return true
	}

	return false
}

// AppArmorStatus gets AppArmor enforced profiles on worker nodes
func (aa *AppArmor) AppArmorStatus() (types.NodeList, error) {
	nodes, err := aa.k8sClient.GetNodes()

	if err != nil {
		return nodes, err
	}

	for _, node := range nodes {
		err := aa.status(node)

		if err != nil {
			return nodes, err
		}
	}

	return nodes, nil
}

func (aa *AppArmor) status(node *types.Node) error {
	if node.IsMaster() {
		return nil
	}

	err := aa.sshClient.Connect(node.ExternalIP, SSH_PORT)
	if err != nil {
		return err
	}

	defer aa.sshClient.Close()

	if !aa.enabledInConnection(node) {
		return nil
	}

	stdout, stderr, err := aa.sshClient.ExecuteOne(commands.AppArmorStatus, true)

	if err != nil {
		return err
	}

	if len(stderr) > 0 {
		return fmt.Errorf(stderr)
	}

	status := types.NewAppArmorStatus()

	err = json.Unmarshal([]byte(stdout), status)

	if err != nil {
		return err
	}

	node.AppArmorStatus = status

	return nil
}
