package inmemory

import "github.com/dpb587/boshua/compiledreleaseversion/datastore"

type Loader func() ([]compiledreleaseversions.CompiledReleaseVersion, error)
type Reloader func() (bool, error)
