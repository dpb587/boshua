package datastore

import "github.com/dpb587/boshua/compiledreleaseversion"

type Index interface {
	Find(compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error)
	Filter(compiledrelease compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error)
	// Store(compiledreleaseversion.Artifact) error
}
