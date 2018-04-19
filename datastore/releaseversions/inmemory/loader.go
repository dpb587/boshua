package inmemory

import "github.com/dpb587/boshua/datastore/releaseversions"

type Loader func() ([]releaseversions.ReleaseVersion, error)
type Reloader func() (bool, error)
