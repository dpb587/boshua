package inmemory

import "github.com/dpb587/boshua/compiledreleaseversion"

type Loader func() ([]compiledreleaseversion.Artifact, error)
type Reloader func() (bool, error)
