package inmemory

import "github.com/dpb587/boshua/stemcellversion"

type Loader func() ([]stemcellversion.Artifact, error)
type Reloader func() (bool, error)
