package artifactutil

import "github.com/dpb587/metalink"

type ArtifactLoader func() (metalink.File, error)
