package clicommon

import "github.com/dpb587/boshua/artifact"

type ArtifactLoader func() (artifact.Artifact, error)
