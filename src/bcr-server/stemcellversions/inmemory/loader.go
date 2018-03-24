package inmemory

import "bcr-server/stemcellversions"

type Loader func() ([]stemcellversions.StemcellVersion, error)
