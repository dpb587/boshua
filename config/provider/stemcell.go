package provider

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/aggregate"
	"github.com/pkg/errors"
)

func (c *Config) SetStemcellFactory(f datastore.Factory) {
	c.stemcellFactory = f
}

func (c *Config) GetStemcellIndex(name string) (datastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	if c.stemcellIndices == nil {
		c.stemcellIndices = map[string]datastore.Index{}
	}

	if idx, found := c.stemcellIndices[name]; found {
		return idx, nil
	}

	var all []datastore.Index

	for _, cfg := range c.Config.Stemcells {
		var idx datastore.Index
		var err error

		idx, err = c.stemcellFactory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, errors.Wrap(err, "creating stemcell version datastore")
		}

		// if cfg.Analysis != nil { // TODO configurable
		var analysisIdx analysisdatastore.Index

		// analysisIndex, err = o.GetAnalysisIndex(cfg.Analysis.Name)
		analysisIdx, err = c.GetAnalysisIndex(analysis.Reference{}) // TODO
		if err != nil {
			return nil, errors.Wrap(err, "loading stemcell analysis")
		}

		idx = datastore.NewAnalysisIndex(idx, analysisIdx)
		// }

		all = append(all, idx)
	}

	return aggregate.New(all...), nil
}
