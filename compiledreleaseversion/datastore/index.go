package datastore

import "github.com/dpb587/boshua/compiledreleaseversion"

type Index interface {
	Find(compiledrelease compiledreleaseversion.Reference) (compiledreleaseversion.Subject, error)
	List() ([]compiledreleaseversion.Subject, error)
}
