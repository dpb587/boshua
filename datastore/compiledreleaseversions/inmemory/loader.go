package inmemory

import "github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"

type Loader func() ([]compiledreleaseversions.CompiledReleaseVersion, error)
