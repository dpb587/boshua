package stemcellversion

import "github.com/dpb587/metalink"

func New(ref Reference, meta4File metalink.File) Artifact {
	// TODO deprecated
	return Artifact{
		IaaS:       ref.IaaS,
		Hypervisor: ref.Hypervisor,
		OS:         ref.OS,
		Version:    ref.Version,
		Light:      ref.Light,
		Tarball:    meta4File,
	}
}
