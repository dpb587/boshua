package scheduler

type Status string

var (
	StatusUnknown   Status = "unknown"
	StatusPending   Status = "pending"
	StatusFailed    Status = "failed"
	StatusRunning   Status = "running"
	StatusFinishing Status = "finishing"
	StatusSucceeded Status = "succeeded"
)
