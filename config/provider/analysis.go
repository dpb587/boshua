package provider

import (
	"fmt"

	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/scheduler"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func (c *Config) SetAnalysisFactory(f datastore.Factory) {
	c.analysisFactory = f
}

func (c *Config) GetAnalysisIndex(name string) (datastore.Index, error) {
	for _, cfg := range c.Config.Analyses.Datastores {
		if cfg.Name == name {
			return c.requireAnalysisIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
		}
	}

	return nil, fmt.Errorf("unrecognized analysis datastore (name: %s)", name)
}

func (c *Config) requireAnalysisIndex(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	if c.analysisIndices == nil {
		c.analysisIndices = map[string]datastore.Index{}
	}

	if _, found := c.analysisIndices[name]; !found {
		idx, err := c.analysisFactory.Create(provider, name, options)
		if err != nil {
			return nil, errors.Wrapf(err, "creating analysis datastore (name: %s)", name)
		}

		c.analysisIndices[name] = idx
	}

	return c.withScheduler(c.analysisIndices[name])
}

func (c *Config) withScheduler(index datastore.Index) (datastore.Index, error) {
	if !c.HasScheduler() {
		return index, nil
	}

	s, err := c.GetScheduler()
	if err != nil {
		return nil, errors.Wrap(err, "loading scheduler")
	}

	var callback schedulerpkg.StatusChangeCallback = nil

	if !c.Config.Global.Quiet {
		callback = schedulerpkg.DefaultStatusChangeCallback
	}

	return scheduler.New(index, s, callback), nil
}
