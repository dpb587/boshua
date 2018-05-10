package clicommon

import (
	"fmt"

	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/metalink"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd
}

func (c *ArtifactCmd) ExecuteAnalysis(loader AnalysisLoader) error {
	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		artifact, err := loader()
		if err != nil {
			return metalink.File{}, fmt.Errorf("finding artifact: %v", err)
		}

		return artifact.ArtifactMetalinkFile(), nil

	})
}
