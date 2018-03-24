package inmemory

import "bcr-server/compiledreleaseversions"

type Loader func() ([]compiledreleaseversions.CompiledReleaseVersion, error)
