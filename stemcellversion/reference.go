package stemcellversion

import (
	"crypto/sha1"
	"fmt"
)

type Reference struct {
	IaaS       string `json:"iaas"`
	Hypervisor string `json:"hypervisor"`
	OS         string `json:"os"`
	Version    string `json:"version"`
	Light      bool   `json:"light"`
	// DiskFormat string
}

func (r Reference) Name() string {
	var prefix string

	if r.Light {
		prefix = "light-"
	}

	return fmt.Sprintf("%sbosh-%s-%s-%s-go_agent", prefix, r.IaaS, r.Hypervisor, r.OS)
}

func (r Reference) UniqueID() string {
	id := sha1.New()
	id.Write([]byte(r.Name()))

	return fmt.Sprintf("%x", id.Sum(nil))
}
