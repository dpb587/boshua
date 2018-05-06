package datastore

import "github.com/dpb587/boshua/compiledreleaseversion"

type Index interface {
	Filter(compiledrelease compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error)
}
