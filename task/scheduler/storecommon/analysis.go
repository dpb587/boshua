package storecommon

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/pivnetfile"
	pivnetfileoptsutil "github.com/dpb587/boshua/pivnetfile/cli/opts/optsutil"
	pivnetfiledatastore "github.com/dpb587/boshua/pivnetfile/datastore"
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
	var storeDatastore string

	switch analysisSubject := analysisRef.Subject.(type) {
	case pivnetfile.Artifact:
		storeDatastore = fmt.Sprintf("internal/pivnetfile/%s", analysisSubject.GetDatastoreName())
		storeArgs = append(
			[]string{"pivnet-file"},
			pivnetfileoptsutil.ArgsFromFilterParams(pivnetfiledatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case stemcellversion.Artifact:
		storeDatastore = fmt.Sprintf("internal/stemcell/%s", analysisSubject.GetDatastoreName())
		storeArgs = append(
			[]string{"stemcell"},
			stemcelloptsutil.ArgsFromFilterParams(stemcellversiondatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case releaseversion.Artifact:
		storeDatastore = fmt.Sprintf("internal/release/%s", analysisSubject.GetDatastoreName())
		storeArgs = append(
			[]string{"release"},
			releaseoptsutil.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case compilation.Artifact:
		storeDatastore = analysisSubject.GetDatastoreName()
		if storeDatastore[0:9] == "internal/" { // swap release -> compilation
			storeDatastore = fmt.Sprintf("internal/release.compilation/%s", strings.SplitN(storeDatastore, "/", 3)[2])
		}

		analysisSubjectRef := analysisSubject.Reference().(compilation.Reference)
		storeArgs = append(
			append(
				append(
					[]string{"release"},
					releaseoptsutil.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromReference(analysisSubjectRef.ReleaseVersion))...,
				),
				"compilation",
			),
			fmt.Sprintf("--stemcell-os=%s", analysisSubjectRef.OSVersion.Name),
			fmt.Sprintf("--stemcell-version=%s", analysisSubjectRef.OSVersion.Version),
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
			fmt.Sprintf("--datastore=%s", storeDatastore),
			filepath.Join("input", "results.jsonl.gz"),
		),
	})

	return tt
}
