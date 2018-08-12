package provider

import (
	"github.com/dpb587/boshua/task/scheduler"
)

func (c *Config) SetSchedulerFactory(f scheduler.Factory) {
	c.schedulerFactory = f
}

func (c *Config) GetScheduler() (scheduler.Scheduler, error) {
	if c.scheduler != nil {
		return c.scheduler, nil
	}

	sched, err := c.schedulerFactory.Create(c.Config.Scheduler.Type, c.Config.Scheduler.Options)

	c.scheduler = sched

	return sched, err
}

func (c *Config) HasScheduler() bool {
	return c.Config.Scheduler.Type != ""
}
