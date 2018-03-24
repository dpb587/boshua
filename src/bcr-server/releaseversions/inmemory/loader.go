package inmemory

import "bcr-server/releaseversions"

type Loader func() ([]releaseversions.ReleaseVersion, error)
