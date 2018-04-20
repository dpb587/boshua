package inmemory

import "github.com/dpb587/boshua/stemcellversion/datastore"

type Loader func() ([]stemcellversions.StemcellVersion, error)
type Reloader func() (bool, error)
