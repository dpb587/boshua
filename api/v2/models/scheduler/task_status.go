package scheduler

type TaskStatus struct {
	Complete bool   `json:"complete"`
	Status   string `json:"status"`
}
