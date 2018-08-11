package scheduler

import (
	"github.com/dpb587/boshua/analysis"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
)

// multiple schedulings of the same thing should not cause duplicates
type Scheduler interface {
	ScheduleCompilation(f compilationdatastore.FilterParams) (ScheduledTask, error)
	ScheduleAnalysis(ref analysis.Reference) (ScheduledTask, error)
}
