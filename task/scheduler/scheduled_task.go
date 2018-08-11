package scheduler

type ScheduledTask interface {
	Status() (Status, error)
	Subject() interface{}
}
