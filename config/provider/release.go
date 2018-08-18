package provider

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/aggregate"
	"github.com/pkg/errors"
)

func (c *Config) SetReleaseFactory(f datastore.Factory) {
	c.releaseFactory = f
}

func (c *Config) GetReleaseIndex(name string) (datastore.Index, error) {
	for _, cfg := range c.Config.Releases.Datastores {
		if cfg.Name == name {
			return c.requireReleaseIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
		}
	}

	if name == "default" {
		var all []datastore.Index

		for _, cfg := range c.Config.Releases.Datastores {
			idx, err := c.requireReleaseIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
			if err != nil {
				return nil, err
			}

			all = append(all, idx)
		}

		return aggregate.New(all...), nil
	}

	return nil, fmt.Errorf("unrecognized release datastore (name: %s)", name)
}

func (c *Config) requireReleaseIndex(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	if c.releaseIndices == nil {
		c.releaseIndices = map[string]datastore.Index{}
	}

	if _, found := c.releaseIndices[name]; !found {
		idx, err := c.releaseFactory.Create(provider, name, options)
		if err != nil {
			return nil, errors.Wrapf(err, "creating release datastore (name: %s)", name)
		}

		c.releaseIndices[name] = idx
	}

	return c.releaseIndices[name], nil
}

func (c *Config) GetReleaseAnalysisIndex(name string) (analysisdatastore.Index, error) {
	for _, cfg := range c.Config.Releases.Datastores {
		if cfg.Name != name {
			continue
		}

		if cfg.Analyses != nil {
			if cfg.Analyses.Type == "" {
				return c.GetAnalysisIndex(cfg.Analyses.Name)
			}

			return c.requireAnalysisIndex(
				analysisdatastore.ProviderName(cfg.Analyses.Type),
				fmt.Sprintf("release/%s/%s", name, cfg.Analyses.Name),
				cfg.Analyses.Options,
			)
		}
	}

	return c.GetAnalysisIndex("default")
}
