package analysisutil

import (
	"log"

	"github.com/dpb587/boshua/cli/client/cmd/artifactutil"
	"github.com/dpb587/metalink"
)

type ArtifactCmd struct {
	artifactutil.ArtifactCmd
}

func (c *ArtifactCmd) ExecuteAnalysis(loader AnalysisLoader) error {
	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		resInfo, err := loader()
		if err != nil {
			log.Fatal(err)
		} else if resInfo == nil {
			log.Fatalf("no analysis available")
		}

		return resInfo.Data.Artifact, nil
	})
}
