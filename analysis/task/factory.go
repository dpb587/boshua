package task

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/scheduler/task"
	"github.com/dpb587/metalink"
)

func New(subject analysis.Subject, analyzer string, privileged bool) (task.Task, error) {
	file := subject.ArtifactMetalinkFile()

	meta4Bytes, err := metalink.Marshal(metalink.Metalink{
		Files: []metalink.File{file},
	})
	if err != nil {
		return task.Task{}, fmt.Errorf("marshaling metalink: %v", err)
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
				fmt.Sprintf("--analyzer=input/%s", file.Name),
				"output/results.jsonl",
			},
		},
	}, nil
}
