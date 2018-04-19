package inmemory

import "github.com/dpb587/boshua/datastore/stemcellversions"

type Loader func() ([]stemcellversions.StemcellVersion, error)
type Reloader func() (bool, error)
