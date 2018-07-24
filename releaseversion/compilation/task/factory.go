package task

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

func New(release releaseversion.Artifact, stemcell stemcellversion.Artifact) (task.Task, error) {
	releaseFile := release.MetalinkFile()

	meta4ReleaseBytes, err := metalink.MarshalXML(metalink.Metalink{
		Files: []metalink.File{releaseFile},
	})
	if err != nil {
		return task.Task{}, errors.Wrap(err, "marshaling release metalink")
	}

	stemcellFile := stemcell.MetalinkFile()

	meta4StemcellBytes, err := metalink.MarshalXML(metalink.Metalink{
		Files: []metalink.File{stemcellFile},
	})
	if err != nil {
		return task.Task{}, errors.Wrap(err, "marshaling stemcell metalink")
	}

	return task.Task{
		Type: task.Type("compilation"),
		Steps: []task.Step{
			{
				Name: "uploading-release",
				Input: map[string][]byte{
					"metalink.meta4": meta4ReleaseBytes,
				},
				Args: []string{
					"artifact",
					"upload-release",
					"input/metalink.meta4",
				},
			},
			{
				Name: "uploading-stemcell",
				Input: map[string][]byte{
					"metalink.meta4": meta4StemcellBytes,
				},
				Args: []string{
					"artifact",
					"upload-stemcell",
					"input/metalink.meta4",
				},
			},
			{
				Name: "compiling",
				Args: []string{
					"release",
					fmt.Sprintf("--release-name=%s", release.Name),
					fmt.Sprintf("--release-version=%s", release.Version),
					"compilation",
					fmt.Sprintf("--os=%s/%s", stemcell.OS, stemcell.Version),
					"export-release",
					filepath.Join("output", fmt.Sprintf("%s-%s-on-%s-stemcell-%s.tgz", release.Name, release.Version, stemcell.OS, stemcell.Version)),
				},
			},
		},
	}, nil
}
