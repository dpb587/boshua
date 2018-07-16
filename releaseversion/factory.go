package releaseversion

import "github.com/dpb587/metalink"

// TODO backcompat remove
func New(ref Reference, meta4File metalink.File) Artifact {
	return Artifact{
		Name:          ref.Name,
		Version:       ref.Version,
		SourceTarball: meta4File,
	}
}
