package storecommon

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/releaseversion"
	releaseoptsutil "github.com/dpb587/boshua/releaseversion/cli/opts/optsutil"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
)

func AppendCompilationStore(tt *task.Task, release releaseversion.Artifact, stemcell stemcellversion.Artifact) *task.Task {
	tt.Steps = append(tt.Steps, task.Step{
		Name: "storing",
		Args: append(
			append(
				[]string{"release"},
				releaseoptsutil.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromArtifact(release))...,
			),
			"compilation",
			fmt.Sprintf("--os=%s/%s", stemcell.OS, stemcell.Version),
			"datastore",
			"store",
			filepath.Join("input", fmt.Sprintf("%s-%s-on-%s-stemcell-%s.tgz", release.Name, release.Version, stemcell.OS, stemcell.Version)),
		),
	})

	return tt
}
