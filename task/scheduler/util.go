package scheduler

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/stemcellversion"
)

func WaitForScheduledTask(t ScheduledTask, callback StatusChangeCallback) (Status, error) {
	currStatus := StatusUnknown

	if callback == nil {
		callback = func(_ ScheduledTask, _ Status) {}
	}

	for {
		status, err := t.Status()
		if status != currStatus {
			callback(t, status)

			currStatus = status
		}

		if status == StatusSucceeded || status == StatusFailed || err != nil {
			return status, err
		}

		time.Sleep(10 * time.Second)
	}
}

func DefaultStatusChangeCallback(task ScheduledTask, status Status) {
	ts := time.Now().Format("15:04:05")

	switch ref := task.Subject().(type) {
	case compilationdatastore.FilterParams: // TODO weird its FilterParams?
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s %s/%s] compilation is %s\n", ts, ref.OS.Name, ref.OS.Version, ref.Release.Name, ref.Release.Version, status))
	case analysis.Reference:
		switch ref := ref.Subject.Reference().(type) {
		case compilation.Reference:
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s %s/%s] analysis is %s\n", ts, ref.OSVersion.Name, ref.OSVersion.Version, ref.ReleaseVersion.Name, ref.ReleaseVersion.Version, status))
		case releaseversion.Reference:
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s] analysis is %s\n", ts, ref.Name, ref.Version, status))
		case stemcellversion.Reference:
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s] analysis is %s\n", ts, ref.FullName(), ref.Version, status))
		default:
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [unknown] analysis is %s\n", ts, status))
		}
	default:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [unknown] task is %s\n", ts, status))
	}
}
