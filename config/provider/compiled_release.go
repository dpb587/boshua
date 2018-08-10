package provider

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/aggregate"
	"github.com/pkg/errors"
)

func (c *Config) SetReleaseCompilationFactory(f datastore.Factory) {
	c.releaseCompilationFactory = f
}

func (c *Config) GetCompiledReleaseIndex(name string) (datastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	releaseIndex, err := c.GetReleaseIndex("default")
	if err != nil {
		return nil, errors.Wrap(err, "loading release index")
	}

	var all []datastore.Index

	for _, cfg := range c.Config.CompiledReleases {
		var idx datastore.Index
		var err error

		idx, err = c.releaseCompilationFactory.Create(cfg.Type, cfg.Name, cfg.Options, releaseIndex)
		if err != nil {
			return nil, errors.Wrap(err, "creating compiled release version datastore")
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
