package task

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

func New(subject analysis.Subject, analyzer analysis.AnalyzerName) (task.Task, error) {
	file := subject.MetalinkFile()

	meta4Bytes, err := metalink.Marshal(metalink.Metalink{
		Files: []metalink.File{file},
	})
	if err != nil {
		return task.Task{}, errors.Wrap(err, "marshaling metalink")
	}

	var privileged bool

	if analyzer == "stemcellimagefiles.v1" { // TODO pass analyzer struct?
		privileged = true
	}

	return task.Task{
		task.Step{
			Name: "downloading",
			Input: map[string][]byte{
				"metalink.meta4": meta4Bytes,
			},
			Args: []string{
				"download-metalink",
				"input/metalink.meta4",
				"output",
			},
		},
		task.Step{
			Name: "analyzing",
			Args: []string{
				"analysis",
				"generate",
				fmt.Sprintf("--analyzer=%s", analyzer),
				filepath.Join("input", file.Name),
				filepath.Join("output", "results.jsonl"),
			},
			Privileged: privileged,
		},
	}, nil
}
