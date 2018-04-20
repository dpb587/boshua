package scheduler

type Status string

var (
	StatusUnknown   Status = "unknown"
	StatusPending   Status = "pending"
	StatusFailed    Status = "failed"
	StatusCompiling Status = "compiling"
	StatusFinishing Status = "finishing"
	StatusSucceeded Status = "succeeded"
)
