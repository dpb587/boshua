package provider

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/pivnetfile/datastore"
	"github.com/pkg/errors"
)

func (c *Config) SetPivnetFileFactory(f datastore.Factory) {
	c.pivnetFileFactory = f
}

func (c *Config) GetPivnetFileIndex(name string) (datastore.Index, error) {
	for _, cfg := range c.Config.PivnetFiles.Datastores {
		if cfg.Name == name {
			return c.requirePivnetFileIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
		}
	}

	if name == config.DefaultName {
		var all []datastore.Index

		for _, cfg := range c.Config.PivnetFiles.Datastores {
			idx, err := c.requirePivnetFileIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
			if err != nil {
				return nil, err
			}

			all = append(all, idx)
		}

		if len(all) == 0 {
			return nil, errors.New("no pivnet file datastores configured")
		}

		// TODO theoretically should support aggregating multiple datasources; current code structure kept for parity when that happens
		return all[0], nil
	}

	return nil, fmt.Errorf("unrecognized pivnet file datastore (name: %s)", name)
}

func (c *Config) requirePivnetFileIndex(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	if c.pivnetFileIndices == nil {
		c.pivnetFileIndices = map[string]datastore.Index{}
	}

	if _, found := c.pivnetFileIndices[name]; !found {
		idx, err := c.pivnetFileFactory.Create(provider, name, options)
		if err != nil {
			return nil, errors.Wrapf(err, "creating pivnet file datastore (name: %s)", name)
		}

		c.pivnetFileIndices[name] = idx
	}

	return c.pivnetFileIndices[name], nil
}

func (c *Config) GetPivnetFileAnalysisIndex(name string) (analysisdatastore.Index, error) {
	for _, cfg := range c.Config.PivnetFiles.Datastores {
		if cfg.Name != name {
			continue
		}

		if cfg.AnalysisDatastore != nil {
			if cfg.AnalysisDatastore.Type == "" {
				return c.getAnalysisIndex(cfg.AnalysisDatastore.Name)
			}

			return c.requireAnalysisIndex(
				analysisdatastore.ProviderName(cfg.AnalysisDatastore.Type),
				fmt.Sprintf("pivnetfile/%s/%s", name, cfg.AnalysisDatastore.Name),
				cfg.AnalysisDatastore.Options,
			)
		}
	}

	return c.getAnalysisIndex(config.DefaultName)
}
