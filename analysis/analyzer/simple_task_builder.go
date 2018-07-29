package task

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

func NewSimpleTask(subject analysis.Subject, analyzer analysis.AnalyzerName, privileged bool) (*task.Task, error) {
	file := subject.MetalinkFile()

	meta4Bytes, err := metalink.MarshalXML(metalink.Metalink{
		Files: []metalink.File{file},
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshaling metalink")
	}

	return &task.Task{
		Type: task.Type("analysis"),
		Steps: []task.Step{
			{
				Name: "downloading",
				Input: map[string][]byte{
					"metalink.meta4": meta4Bytes,
				},
				Args: []string{
					"artifact",
					"download",
					"input/metalink.meta4",
					"output",
				},
			},
			{
				Name: "analyzing",
				Args: []string{
					"analysis",
					"generate",
					fmt.Sprintf("--analyzer=%s", analyzer),
					filepath.Join("input", file.Name),
					filepath.Join("output", "results.jsonl.gz"),
				},
				Privileged: privileged,
			},
		},
	}, nil
}
