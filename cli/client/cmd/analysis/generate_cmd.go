package analysis

import (
	"fmt"
	"os"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
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

	analyzer, err := factory.Factory{}.Create(c.AnalysisOpts.Analyzer, c.Args.Artifact)
	if err != nil {
		return fmt.Errorf("finding analyzer: %s", c.AnalysisOpts.Analyzer)
	}

	return analyzer.Analyze(analysis.NewJSONWriter(os.Stdout))
}
