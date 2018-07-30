package scheduler

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
)

// multiple schedulings of the same thing should not cause duplicates
type Scheduler interface {
	ScheduleCompilation(release releaseversion.Artifact, stemcell stemcellversion.Artifact) (task.ScheduledTask, error)
	ScheduleAnalysis(analysisRef analysis.Reference) (task.ScheduledTask, error)
}
