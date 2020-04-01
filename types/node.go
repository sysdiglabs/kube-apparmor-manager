package types

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	RoleLabel = "kubernetes.io/role"
	Worker    = "node"
	Master    = "master"
)

type NodeList []*Node

type Node struct {
	NodeName        string
	ExternalIP      string
	InternalIP      string
	Role            string
	AppArmorEnabled bool
	AppArmorStatus  *AppArmorProfileStatus
}

// NewNode returns a new node object
func NewNode() *Node {
	return &Node{
		AppArmorStatus: NewAppArmorStatus(),
	}
}

func (nl NodeList) String() string {
	ret := ""

	ret += strings.Join([]string{"Node Name", "Internal IP", "External IP", "Role", "AppArmor Enabled"}, "\t")
	ret += "\n"

	for _, n := range nl {
		ret += strings.Join([]string{n.NodeName, n.InternalIP, n.ExternalIP, n.Role, fmt.Sprintf("%t", n.AppArmorEnabled)}, "\t")
		ret += "\n"
	}

	return ret
}

func (nl NodeList) GetEnforcedProfiles() string {
	ret := ""

	ret += strings.Join([]string{"Node Name", "Role", "Enforced Profiles"}, "\t")
	ret += "\n"

	for _, n := range nl {
		ret += strings.Join([]string{n.NodeName, n.Role, strings.Join(n.AppArmorStatus.GetEnforcedProfiles(), ",")}, "\t")
		ret += "\n"
	}

	return ret
}

// IsMaster checks whether a node is master node
func (n *Node) IsMaster() bool {
	return n.Role == Master
}

// PrintEnforcementStatus prints enforced AppArmor profile on worker nodes
func (nl NodeList) PrintEnforcementStatus() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node Name", "Role", "Enforced Profiles"})

	data := [][]string{}

	for _, n := range nl {
		data = append(data, []string{n.NodeName, n.Role, strings.Join(n.AppArmorStatus.GetEnforcedProfiles(), ",")})
	}

	table.AppendBulk(data)
	table.Render()
}

// PrintEnabledStatus prints AppArmor enabled status on worker nodes
func (nl NodeList) PrintEnabledStatus() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node Name", "Internal IP", "External IP", "Role", "AppArmor Enabled"})

	data := [][]string{}

	for _, n := range nl {
		data = append(data, []string{n.NodeName, n.InternalIP, n.ExternalIP, n.Role, fmt.Sprintf("%t", n.AppArmorEnabled)})
	}

	table.AppendBulk(data)
	table.Render()
}
