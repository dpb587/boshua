package task

import "time"

type ScheduledTask interface {
	Status() (Status, error)
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
