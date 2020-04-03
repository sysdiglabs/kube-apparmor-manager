package types

import (
	"fmt"
	"sort"
	"strings"
)

const (
	enforced = "enforce"
)

type AppArmorProfileStatus struct {
	Profiles map[string]string `json:"profiles"`
}

// NewAppArmorStatus return apparmor profile status from workder nodes
func NewAppArmorStatus() *AppArmorProfileStatus {
	return &AppArmorProfileStatus{
		Profiles: map[string]string{},
	}
}

// GetEnforcedProfiles get enforced profile names
func (s *AppArmorProfileStatus) GetEnforcedProfiles() []string {
	profiles := []string{}

	for k, v := range s.Profiles {
		if v == enforced {
			profiles = append(profiles, k)
		}
	}

	sort.Strings(profiles)

	return profiles
}

type AppArmorProfile struct {
	Name  string
	Rules string
	// Profile contains lines of the profile as following

	/*
			capability net_raw,
			capability setuid,
			capability setgid,
			capability dac_override,
			network raw,
			network packet,

			# for -D
			capability sys_module,
			@{PROC}/bus/usb/ r,
			@{PROC}/bus/usb/** r,

		  	audit deny @{HOME}/bin/ rw,
		  	audit deny @{HOME}/bin/** mrwkl,
		 	 @{HOME}/ r,
		  	@{HOME}/** rw,

		  	/usr/sbin/tcpdump r,
	*/
	Enforced bool
}

func (p AppArmorProfile) String() string {
	ret := ""

	ret += fmt.Sprintf("profile %s flags=(attach_disconnected,mediate_deleted) {\n", p.Name)

	lines := strings.Split(p.Rules, "\n")

	for _, line := range lines {
		// skip empty line and commented out rule
		if line == "" || line[0] == '#' {
			continue
		}

		last := line[len(line)-1:]

		// append ',' at the end if it doesn't exists
		if last != "," {
			line += ","
		}
		ret += fmt.Sprintf("\t%s\n", line)
	}

	ret += "}"
	return ret
}
