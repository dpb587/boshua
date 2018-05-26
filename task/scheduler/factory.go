package scheduler

type Factory interface {
	Create(provider string, options map[string]interface{}) (Scheduler, error)
}
