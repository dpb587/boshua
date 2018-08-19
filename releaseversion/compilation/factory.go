package compilation

import (
	"github.com/dpb587/metalink"
)

func New(datastore string, ref Reference, meta4File metalink.File) Artifact {
	return Artifact{
		Datastore: datastore,
		Release:   ref.ReleaseVersion,
		OS:        ref.OSVersion,
		Tarball:   meta4File,
	}
}
