package provider

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/scheduler"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func (c *Config) SetAnalysisFactory(f datastore.Factory) {
	c.analysisFactory = f
}

func (c *Config) GetAnalysisIndex(_ analysis.Reference) (datastore.Index, error) {
	// TODO decide between name and analysis reference
	name := "default"

	for _, cfg := range c.Config.Analyses {
		if cfg.Name != name {
			continue
		}

		idx, err := c.analysisFactory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, errors.Wrap(err, "creating analysis datastore")
		}

		return idx, nil
	}

	return nil, fmt.Errorf("failed to find analysis index: %s", name)
}

func (c *Config) GetAnalysisIndexScheduler(ref analysis.Reference) (datastore.Index, error) {
	index, err := c.GetAnalysisIndex(ref)
	if err != nil {
		return nil, err
	}

	if !c.HasScheduler() {
		return index, nil
	}

	s, err := c.GetScheduler()
	if err != nil {
		return nil, errors.Wrap(err, "loading scheduler")
	}

	return scheduler.New(index, s, func(status schedulerpkg.Status) {
		fmt.Fprintf(os.Stderr, "%s [%s/%s] analysis is %s\n", time.Now().Format("15:04:05"), "TODO", "TODO", status)
	}), nil
}
