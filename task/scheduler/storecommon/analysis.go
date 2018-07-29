package storecommon

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/releaseversion"
	releaseoptsutil "github.com/dpb587/boshua/releaseversion/cli/opts/optsutil"
	"github.com/dpb587/boshua/releaseversion/compilation"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	stemcelloptsutil "github.com/dpb587/boshua/stemcellversion/cli/opts/optsutil"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/task"
	"github.com/pkg/errors"
)

func AppendAnalysisStore(tt *task.Task, analysisRef analysis.Reference) *task.Task {
	var storeArgs []string

	switch analysisSubject := analysisRef.Subject.(type) {
	case stemcellversion.Artifact:
		storeArgs = append(
			[]string{"stemcell"},
			stemcelloptsutil.ArgsFromFilterParams(stemcellversiondatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case releaseversion.Artifact:
		storeArgs = append(
			[]string{"release"},
			releaseoptsutil.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case compilation.Artifact:
		analysisSubjectRef := analysisSubject.Reference().(compilation.Reference)
		storeArgs = append(
			append(
				append(
					[]string{"release"},
					releaseoptsutil.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromReference(analysisSubjectRef.ReleaseVersion))...,
				),
				"compilation",
			),
			fmt.Sprintf("--os=%s/%s", analysisSubjectRef.OSVersion.Name, analysisSubjectRef.OSVersion.Version),
		)
	default:
		panic(errors.New("unsupported analysis subject")) // TODO panic?
	}

	tt.Steps = append(tt.Steps, task.Step{
		Name: "storing",
		Args: append(
			storeArgs,
			"analysis",
			"store-results",
			fmt.Sprintf("--analyzer=%s", analysisRef.Analyzer),
			filepath.Join("input", "results.jsonl.gz"),
		),
	})

	return tt
}
