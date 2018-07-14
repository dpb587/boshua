package stemcellversion

import (
	"crypto/sha1"
	"fmt"
)

type Reference struct {
	Name       string `json:"name"`
	IaaS       string `json:"iaas"`
	Hypervisor string `json:"hypervisor"`
	OS         string `json:"os"`
	Version    string `json:"version"`
	Light      bool   `json:"light"`
	// DiskFormat string `json:"disk_format"`
}

func (r Reference) FullName() string {
	var prefix string

	if r.Light {
		prefix = "light-"
	}

	return fmt.Sprintf("%sbosh-%s-%s-%s-go_agent", prefix, r.IaaS, r.Hypervisor, r.OS)
}

func (r Reference) UniqueID() string {
	id := sha1.New()
	id.Write([]byte(r.FullName()))

	return fmt.Sprintf("%x", id.Sum(nil))
}
