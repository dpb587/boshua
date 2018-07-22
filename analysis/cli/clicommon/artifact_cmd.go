package clicommon

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/pkg/errors"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd
}

func (c *ArtifactCmd) ExecuteAnalysis(loader AnalysisLoader) error {
	return c.ArtifactCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := loader()
		if err != nil {
			return nil, errors.Wrap(err, "finding artifact")
		}

		return artifact, nil
	})
}
