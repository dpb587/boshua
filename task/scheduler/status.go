package scheduler

type Status string

var (
	StatusUnknown   Status = "unknown"
	StatusPending   Status = "pending"
	StatusFailed    Status = "failed"
	StatusRunning   Status = "running"
	StatusFinishing Status = "finishing"
	StatusSucceeded Status = "finished"
)

type StatusChangeCallback func(ScheduledTask, Status) // TODO pass only ScheduledTask.Subject()?
