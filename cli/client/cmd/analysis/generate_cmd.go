package analysis

import (
	"fmt"
	"os"

	"github.com/dpb587/bosh-compiled-releases/analysis"
	releaseartifactchecksumsv1 "github.com/dpb587/bosh-compiled-releases/analysis/releaseartifactchecksums.v1/analyzer"
	releaseartifactfilestatv1 "github.com/dpb587/bosh-compiled-releases/analysis/releaseartifactfilestat.v1/analyzer"
	releasemanifestsv1 "github.com/dpb587/bosh-compiled-releases/analysis/releasemanifests.v1/analyzer"
)

type GenerateCmd struct {
	*CmdOpts `no-flag:"true"`

	Args GenerateArgs `positional-args:"true"`
}

type GenerateArgs struct {
	Artifact string `positional-arg-name:"ARTIFACT-PATH" description:"Artifact path to analyze"`
}

func (c *GenerateCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/generate")

	var analyzer analysis.Analyzer

	if c.AnalysisOpts.Analyzer == "releaseartifactchecksums.v1" {
		analyzer = releaseartifactchecksumsv1.New(c.Args.Artifact)
	} else if c.AnalysisOpts.Analyzer == "releaseartifactfilestat.v1" {
		analyzer = releaseartifactfilestatv1.New(c.Args.Artifact)
	} else if c.AnalysisOpts.Analyzer == "releasemanifests.v1" {
		analyzer = releasemanifestsv1.New(c.Args.Artifact)
	} else {
		return fmt.Errorf("invalid analyzer: %s", c.AnalysisOpts.Analyzer)
	}

	return analyzer.Analyze(analysis.NewJSONWriter(os.Stdout))
}
