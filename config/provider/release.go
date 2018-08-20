package provider

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/config"
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

	if name == config.DefaultName {
		var all []datastore.Index

		for _, cfg := range c.Config.Releases.Datastores {
			idx, err := c.requireReleaseIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
			if err != nil {
				return nil, err
			}

			all = append(all, idx)
		}

		if len(all) == 0 {
			return nil, errors.New("no release datastores configured")
		}

		return aggregate.New(name, all...), nil
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

		if cfg.AnalysisDatastore != nil {
			if cfg.AnalysisDatastore.Type == "" {
				return c.getAnalysisIndex(cfg.AnalysisDatastore.Name)
			}

			return c.requireAnalysisIndex(
				analysisdatastore.ProviderName(cfg.AnalysisDatastore.Type),
				fmt.Sprintf("release/%s/%s", name, cfg.AnalysisDatastore.Name),
				cfg.AnalysisDatastore.Options,
			)
		}
	}

	return c.getAnalysisIndex(config.DefaultName)
}
