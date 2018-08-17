package provider

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/aggregate"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/scheduler"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func (c *Config) SetReleaseCompilationFactory(f datastore.Factory) {
	c.releaseCompilationFactory = f
}

func (c *Config) GetReleaseCompilationIndex(name string) (datastore.Index, error) {
	for _, cfg := range c.Config.ReleaseCompilations {
		if cfg.Name == name {
			return c.requireReleaseCompilationIndex(cfg.Name, cfg.Type, cfg.Options)
		}
	}

	if name == "default" {
		var all []datastore.Index

		for _, cfg := range c.Config.ReleaseCompilations {
			idx, err := c.requireReleaseCompilationIndex(cfg.Name, cfg.Type, cfg.Options)
			if err != nil {
				return nil, err
			}

			all = append(all, idx)
		}

		return aggregate.New(all...), nil
	}

	return nil, fmt.Errorf("unrecognized release compilation datastore (name: %s)", name)
}

func (c *Config) requireReleaseCompilationIndex(name, provider string, options map[string]interface{}) (datastore.Index, error) {
	if c.releaseCompilationIndices == nil {
		c.releaseCompilationIndices = map[string]datastore.Index{}
	}

	if _, found := c.releaseCompilationIndices[name]; !found {
		idx, err := c.releaseCompilationFactory.Create(datastore.ProviderName(provider), name, options)
		if err != nil {
			return nil, errors.Wrapf(err, "creating release compilation datastore (name: %s)", name)
		}

		c.releaseCompilationIndices[name] = idx
	}

	return c.withReleaseCompilationScheduler(c.releaseCompilationIndices[name])
}

func (c *Config) withReleaseCompilationScheduler(index datastore.Index) (datastore.Index, error) {
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

func (c *Config) GetReleaseCompilationAnalysisIndex(name string) (analysisdatastore.Index, error) {
	for _, cfg := range c.Config.ReleaseCompilations {
		if cfg.Name != name {
			continue
		}

		if cfg.Analysis != nil {
			if cfg.Analysis.Type == "" {
				return c.GetAnalysisIndex(cfg.Analysis.Name)
			}

			return c.requireAnalysisIndex(
				analysisdatastore.ProviderName(cfg.Analysis.Type),
				fmt.Sprintf("release-compilation/%s/%s", name, cfg.Analysis.Name),
				cfg.Analysis.Options,
			)
		}

		return c.GetAnalysisIndex("default")
	}

	if name != "default" {
		return nil, fmt.Errorf("unrecognized release compilation datastore (name: %s)", name)
	}

	return c.GetAnalysisIndex("default")
}
