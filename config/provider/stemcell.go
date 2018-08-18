package provider

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	osversionstemcellversionindex "github.com/dpb587/boshua/osversion/datastore/stemcellversionindex"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/aggregate"
	"github.com/pkg/errors"
)

func (c *Config) SetStemcellFactory(f datastore.Factory) {
	c.stemcellFactory = f
}

func (c *Config) GetStemcellIndex(name string) (datastore.Index, error) {
	for _, cfg := range c.Config.Stemcells.Datastores {
		if cfg.Name == name {
			return c.requireStemcellIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
		}
	}

	if name == "default" {
		var all []datastore.Index

		for _, cfg := range c.Config.Stemcells.Datastores {
			idx, err := c.requireStemcellIndex(datastore.ProviderName(cfg.Type), cfg.Name, cfg.Options)
			if err != nil {
				return nil, err
			}

			all = append(all, idx)
		}

		return aggregate.New(all...), nil
	}

	return nil, fmt.Errorf("unrecognized stemcell datastore (name: %s)", name)
}

func (c *Config) requireStemcellIndex(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	if c.stemcellIndices == nil {
		c.stemcellIndices = map[string]datastore.Index{}
	}

	if _, found := c.stemcellIndices[name]; !found {
		idx, err := c.stemcellFactory.Create(provider, name, options)
		if err != nil {
			return nil, errors.Wrapf(err, "creating stemcell datastore (name: %s)", name)
		}

		c.stemcellIndices[name] = idx
	}

	return c.stemcellIndices[name], nil
}

func (c *Config) GetStemcellAnalysisIndex(name string) (analysisdatastore.Index, error) {
	for _, cfg := range c.Config.Stemcells.Datastores {
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

// TODO remove/move
func (c *Config) GetOSIndex(name string) (osversiondatastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	stemcellVersionIndex, err := c.GetStemcellIndex("default")
	if err != nil {
		return nil, errors.Wrap(err, "loading stemcell index")
	}

	return osversionstemcellversionindex.New(stemcellVersionIndex, c.GetLogger()), nil
}
