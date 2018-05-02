package osversion

import "github.com/dpb587/metalink"

func New(ref Reference, meta4File metalink.File, meta4Source map[string]interface{}) Artifact {
	return Artifact{
		Reference:      ref,
		metalinkFile:   meta4File,
		metalinkSource: meta4Source,
	}
}
