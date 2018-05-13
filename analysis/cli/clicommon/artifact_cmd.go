package clicommon

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd
}

func (c *ArtifactCmd) ExecuteAnalysis(loader AnalysisLoader) error {
	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		artifact, err := loader()
		if err != nil {
			return metalink.File{}, errors.Wrap(err, "finding artifact")
		}

		return artifact.MetalinkFile(), nil

	})
}
