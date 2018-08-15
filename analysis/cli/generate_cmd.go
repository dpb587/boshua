package cli

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
	"github.com/pkg/errors"
)

type GenerateCmd struct {
	Analyzer   analysis.AnalyzerName `long:"analyzer" description:"The analyzer to use"`
	NoCompress bool                  `long:"no-compress" description:"Skip gzip compression when writing results to file"`

	Args GenerateArgs `positional-args:"true"`
}

type GenerateArgs struct {
	Artifact string `positional-arg-name:"ARTIFACT-PATH" description:"Artifact path to analyze"`
	Output   string `positional-arg-name:"OUTPUT-PATH" description:"Path to output results (default: STDOUT)" optional:"true"`
}

func (c *GenerateCmd) Execute(_ []string) error {
	analyzer, err := factory.Factory{}.Create(c.Analyzer, c.Args.Artifact)
	if err != nil {
		return fmt.Errorf("finding analyzer: %s", c.Analyzer)
	}

	var fh io.WriteCloser = os.Stdout

	if c.Args.Output != "" && c.Args.Output != "-" {
		fullPath, err := filepath.Abs(c.Args.Output)
		if err != nil {
			return errors.Wrap(err, "finding output file")
		}

		fh, err = os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrap(err, "opening output file")
		}

		if !c.NoCompress {
			fh = gzip.NewWriter(fh)

			defer fh.Close()
		}
	}

	return analyzer.Analyze(analysis.NewJSONWriter(fh))
}
