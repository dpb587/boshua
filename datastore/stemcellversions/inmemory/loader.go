package inmemory

import "github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"

type Loader func() ([]stemcellversions.StemcellVersion, error)
type Reloader func() (bool, error)
