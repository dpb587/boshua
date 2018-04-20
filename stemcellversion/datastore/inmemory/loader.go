package inmemory

import "github.com/dpb587/boshua/stemcellversion"

type Loader func() ([]stemcellversion.Subject, error)
type Reloader func() (bool, error)
