package inmemory

import "github.com/dpb587/boshua/datastore/compiledreleaseversions"

type Loader func() ([]compiledreleaseversions.CompiledReleaseVersion, error)
type Reloader func() (bool, error)
