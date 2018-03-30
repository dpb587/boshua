package inmemory

import "github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"

type Loader func() ([]releaseversions.ReleaseVersion, error)
type Reloader func() (bool, error)
