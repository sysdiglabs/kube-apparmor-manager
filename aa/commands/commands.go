package commands

import (
	"fmt"

	"github.com/sysdiglabs/kube-apparmor-manager/types"
)

var (
	AAEnable        = "aa-enabled"
	InstallAppArmor = []string{
		`apt install -y apparmor-profiles apparmor-utils`,
		`sed -i -e '/^GRUB_CMDLINE_LINUX_DEFAULT/s/"$/ apparmor=1 security=apparmor"/' /etc/default/grub`,
		`update-grub`,
		`reboot`,
	}

	CreateAppArmorProfileTemplate = []string{
		`echo '%s' > /tmp/%s`,
		`mv /tmp/%s /etc/apparmor.d/%s`,
	}

	EnforceAppArmorProfileTemplate = []string{
		`rm /etc/apparmor.d/disable/%s 2>/dev/null`,
		`aa-enforce /etc/apparmor.d/%s`,
	}

	AppArmorStatus = "apparmor_status --json"

	DisableAppArmorProfileTempalte = []string{
		`ln -s /etc/apparmor.d/%s /etc/apparmor.d/disable/`,
		`apparmor_parser -R /etc/apparmor.d/%s`,
	}
)

// CreateProfileCommands returns a list of commands to create AppArmor profiles on worker nodes
func CreateProfileCommands(profile types.AppArmorProfile) []string {
	commands := make([]string, 2)

	commands[0] = fmt.Sprintf(CreateAppArmorProfileTemplate[0], profile, profile.Name)

	commands[1] = fmt.Sprintf(CreateAppArmorProfileTemplate[1], profile.Name, profile.Name)

	//commands[4] = CreateAppArmorProfileTemplate[4]

	return commands
}

// EnforceProfileCommands returns a list of commands to enforce profile on worker nodes
func EnforceProfileCommands(profile types.AppArmorProfile) []string {
	commands := make([]string, 2)

	commands[0] = fmt.Sprintf(EnforceAppArmorProfileTemplate[0], profile.Name)

	commands[1] = fmt.Sprintf(EnforceAppArmorProfileTemplate[1], profile.Name)

	return commands
}

// DisableProfileCommands returns a list of commands to disable profile on worker nodes
func DisableProfileCommands(profile types.AppArmorProfile) []string {
	commands := make([]string, 2)

	commands[0] = fmt.Sprintf(DisableAppArmorProfileTempalte[0], profile.Name)

	commands[1] = fmt.Sprintf(DisableAppArmorProfileTempalte[1], profile.Name)

	return commands
}
