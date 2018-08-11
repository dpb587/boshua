package scheduler

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/releaseversion/compilation"
)

type ScheduledTask interface {
	Status() (Status, error)
	Subject() interface{}
}

func WaitForScheduledTask(t ScheduledTask, callback StatusChangeCallback) (Status, error) {
	currStatus := StatusUnknown

	if callback == nil {
		callback = func(_ Status) {}
	}

	for {
		status, err := t.Status()
		if status != currStatus {
			callback(status)

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
	case compilation.Reference:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s %s/%s] compilation is %s\n", ts, ref.OSVersion.Name, ref.OSVersion.Version, ref.ReleaseVersion.Name, ref.ReleaseVersion.Version, status))
	default:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [unknown] status is %s\n", ts, status))
	}
}
