package provider

import (
	"github.com/dpb587/boshua/task/scheduler"
)

func (c *Config) SetSchedulerFactory(f scheduler.Factory) {
	c.schedulerFactory = f
}

func (c *Config) GetScheduler() (scheduler.Scheduler, error) {
	return c.schedulerFactory.Create(c.Config.Scheduler.Type, c.Config.Scheduler.Options)
}

func (c *Config) HasScheduler() bool {
	return c.Config.Scheduler.Type != ""
}
