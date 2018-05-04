package analysisutil

import (
	"fmt"
	"log"

	"github.com/dpb587/metalink"
)

type MetalinkCmd struct {
	Format string `long:"format" description:"Output format for metalink"`
}

func (c *MetalinkCmd) ExecuteAnalysis(loader AnalysisLoader) error {
	resInfo, err := loader()
	if err != nil {
		log.Fatal(err)
	} else if resInfo == nil {
		log.Fatalf("no analysis available")
	}

	meta4 := metalink.Metalink{
		Files: []metalink.File{
			resInfo.Data.Artifact,
		},
		Generator: "bosh-compiled-releases/0.0.0",
	}

	meta4Bytes, err := metalink.Marshal(meta4)
	if err != nil {
		log.Fatalf("marshalling response: %v", err)
	}

	fmt.Printf("%s\n", meta4Bytes)

	return nil
}
