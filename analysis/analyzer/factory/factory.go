package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	tilearchivefilesv1 "github.com/dpb587/boshua/pivnetfile/analyzers/tilearchivefiles.v1"
	tilereleasemanifestsv1 "github.com/dpb587/boshua/pivnetfile/analyzers/tilereleasemanifests.v1"
	releaseartifactfilesv1 "github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1"
	releasemanifestsv1 "github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1"
	stemcellimagefilesv1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellimagefiles.v1"
	stemcellmanifestv1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1"
	stemcellpackagesv1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1"
	"github.com/dpb587/boshua/task"
)

type Factory struct{}

func (Factory) Create(analyzer analysis.AnalyzerName, path string) (analysis.AnalysisGenerator, error) {
	// TODO deprecate this to AppOpts service; should be configurable to support dynamic, external analysis
	switch analyzer {
	case releaseartifactfilesv1.AnalyzerName:
		return releaseartifactfilesv1.NewAnalysis(path), nil
	case releasemanifestsv1.AnalyzerName:
		return releasemanifestsv1.NewAnalysis(path), nil
	case stemcellimagefilesv1.AnalyzerName:
		return stemcellimagefilesv1.NewAnalysis(path), nil
	case stemcellmanifestv1.AnalyzerName:
		return stemcellmanifestv1.NewAnalysis(path), nil
	case stemcellpackagesv1.AnalyzerName:
		return stemcellpackagesv1.NewAnalysis(path), nil
	case tilereleasemanifestsv1.AnalyzerName:
		return tilereleasemanifestsv1.NewAnalysis(path), nil
	case tilearchivefilesv1.AnalyzerName:
		return tilearchivefilesv1.NewAnalysis(path), nil
	}

	return nil, fmt.Errorf("unknown analyzer: %s", analyzer)
}

func (Factory) BuildTask(analyzer analysis.AnalyzerName, subject analysis.Subject) (*task.Task, error) {
	// TODO deprecate this to AppOpts service; should be configurable to support dynamic, external analysis
	switch analyzer {
	case releaseartifactfilesv1.AnalyzerName:
		return releaseartifactfilesv1.Analyzer.BuildTask(subject)
	case releasemanifestsv1.AnalyzerName:
		return releasemanifestsv1.Analyzer.BuildTask(subject)
	case stemcellimagefilesv1.AnalyzerName:
		return stemcellimagefilesv1.Analyzer.BuildTask(subject)
	case stemcellmanifestv1.AnalyzerName:
		return stemcellmanifestv1.Analyzer.BuildTask(subject)
	case stemcellpackagesv1.AnalyzerName:
		return stemcellpackagesv1.Analyzer.BuildTask(subject)
	case tilereleasemanifestsv1.AnalyzerName:
		return tilereleasemanifestsv1.Analyzer.BuildTask(subject)
	case tilearchivefilesv1.AnalyzerName:
		return tilearchivefilesv1.Analyzer.BuildTask(subject)
	}

	return nil, fmt.Errorf("unknown analyzer: %s", analyzer)
}

var SoonToBeDeprecatedFactory = &Factory{}
