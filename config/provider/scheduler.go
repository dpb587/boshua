package provider

import (
	"errors"

	"github.com/dpb587/boshua/task/scheduler"
)

func (c *Config) SetSchedulerFactory(f scheduler.Factory) {
	c.schedulerFactory = f
}

func (c *Config) GetScheduler() (scheduler.Scheduler, error) {
	if c.scheduler != nil {
		return c.scheduler, nil
	} else if c.Config.Scheduler == nil {
		return nil, errors.New("no scheduler configured")
	}

	sched, err := c.schedulerFactory.Create(c.Config.Scheduler.Type, c.Config.Scheduler.Options)

	c.scheduler = sched

	return sched, err
}

func (c *Config) HasScheduler() bool {
	return c.Config.Scheduler != nil && c.Config.Scheduler.Type != ""
}
