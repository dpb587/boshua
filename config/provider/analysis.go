package provider

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
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
