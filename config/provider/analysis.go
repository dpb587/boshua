package provider

import (
	"fmt"

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

	if c.analysisIndices == nil {
		c.analysisIndices = map[string]datastore.Index{}
	}

	if idx, found := c.analysisIndices[name]; found {
		return idx, nil
	}

	for _, cfg := range c.Config.Analyses {
		if cfg.Name != name {
			continue
		}

		idx, err := c.analysisFactory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, errors.Wrap(err, "creating analysis datastore")
		}

		c.analysisIndices[name] = idx

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

	var callback schedulerpkg.StatusChangeCallback = nil

	if !c.Config.General.Quiet {
		callback = schedulerpkg.DefaultStatusChangeCallback
	}

	return scheduler.New(index, s, callback), nil
}
