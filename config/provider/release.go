package provider

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/aggregate"
	"github.com/pkg/errors"
)

func (c *Config) SetReleaseFactory(f datastore.Factory) {
	c.releaseFactory = f
}

func (c *Config) GetReleaseIndex(name string) (datastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	if c.releaseIndices == nil {
		c.releaseIndices = map[string]datastore.Index{}
	}

	if idx, found := c.releaseIndices[name]; found {
		return idx, nil
	}

	var all []datastore.Index

	for _, cfg := range c.Config.Releases {
		var idx datastore.Index
		var err error

		idx, err = c.releaseFactory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, errors.Wrap(err, "creating release version datastore")
		}

		// if cfg.Analysis != nil { // TODO configurable
		var analysisIdx analysisdatastore.Index

		// analysisIndex, err = o.GetAnalysisIndex(cfg.Analysis.Name)
		analysisIdx, err = c.GetAnalysisIndex(analysis.Reference{}) // TODO
		if err != nil {
			return nil, errors.Wrap(err, "loading release analysis")
		}

		idx = datastore.NewAnalysisIndex(idx, analysisIdx)
		// }

		all = append(all, idx)
	}

	return aggregate.New(all...), nil
}
