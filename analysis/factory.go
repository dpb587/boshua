package analysis

import (
	"github.com/dpb587/metalink"
)

func New(datastore string, ref Reference, meta4File metalink.File) Artifact {
	return Artifact{
		Datastore:    datastore,
		reference:    ref,
		metalinkFile: meta4File,
	}
}
