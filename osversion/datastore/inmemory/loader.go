package inmemory

import "github.com/dpb587/boshua/osversion"

type Loader func() ([]osversion.Artifact, error)
type Reloader func() (bool, error)
