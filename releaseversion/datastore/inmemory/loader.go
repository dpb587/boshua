package inmemory

import "github.com/dpb587/boshua/releaseversion/datastore"

type Loader func() ([]releaseversions.ReleaseVersion, error)
type Reloader func() (bool, error)
