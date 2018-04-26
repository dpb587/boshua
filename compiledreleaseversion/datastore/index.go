package datastore

import "github.com/dpb587/boshua/compiledreleaseversion"

type Index interface {
	Find(compiledrelease compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error)
	List() ([]compiledreleaseversion.Artifact, error)
}
