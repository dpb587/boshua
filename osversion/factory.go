package osversion

import "github.com/dpb587/metalink"

func New(ref Reference, meta4File metalink.File, meta4Source map[string]interface{}) Artifact {
	return Artifact{
		Reference:      ref,
		MetalinkFile:   meta4File,
		MetalinkSource: meta4Source,
	}
}
