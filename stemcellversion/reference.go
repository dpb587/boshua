package stemcellversion

import "fmt"

type Reference struct {
	IaaS       string
	Hypervisor string
	OS         string
	Version    string
	Light      bool
	// DiskFormat string
}

func (r Reference) Name() string {
	var prefix string

	if r.Light {
		prefix = "light-"
	}

	return fmt.Sprintf("%sbosh-%s-%s-%s-go_agent", prefix, r.IaaS, r.Hypervisor, r.OS)
}
