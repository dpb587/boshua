package inmemory

import "github.com/dpb587/boshua/compiledreleaseversion"

type Loader func() ([]compiledreleaseversion.Subject, error)
type Reloader func() (bool, error)
