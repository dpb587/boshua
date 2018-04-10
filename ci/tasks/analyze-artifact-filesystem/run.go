package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dpb587/bosh-compiled-releases/analysis"
	"github.com/dpb587/bosh-compiled-releases/analysis/releaseartifactfilestat.v1/analyzer"
)

func main() {
	paths, err := filepath.Glob(os.Args[1])
	if err != nil {
		log.Fatalf("globbing: %v", err)
	} else if len(paths) != 1 {
		log.Fatalf("found %d matches for %s", len(paths), os.Args[1])
	}

	analyzer := analyzer.New(paths[0])

	err = analyzer.Analyze(analysis.NewJSONWriter(os.Stdout))
	if err != nil {
		log.Fatalf("analyzing: %v", err)
	}
}
