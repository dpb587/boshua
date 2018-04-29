package stemcellversion

type Reference struct {
	IaaS       string
	Hypervisor string
	OS         string
	Version    string
	Light      bool
	// DiskFormat string
}
