package inmemory

import "github.com/dpb587/boshua/releaseversion"

type Loader func() ([]releaseversion.Subject, error)
type Reloader func() (bool, error)
