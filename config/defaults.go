package config

func (c *Config) ApplyDefaults() {
	if c.Server.Bind == "" {
		c.Server.Bind = "127.0.0.1:4508"
	}

	if c.Global.DefaultServer != "" {
		defaultServer := AbstractComponentConfig{
			Name: "default",
			Type: "boshua.v2",
			Options: map[string]interface{}{
				"url": c.Global.DefaultServer,
			},
		}

		if c.Scheduler == nil {
			c.Scheduler = &defaultServer
		}

		if len(c.Analyses.Datastores) == 0 { // TODO check for name = default instead?
			c.Analyses.Datastores = append(c.Analyses.Datastores, AnalysisDatastoreConfig{
				AbstractComponentConfig: defaultServer,
			})
		}

		if len(c.Releases.Datastores) == 0 { // TODO check for name = default instead?
			c.Releases.Datastores = append(c.Releases.Datastores, ReleaseDatastoreConfig{
				AbstractComponentConfig: defaultServer,
			})
		}

		if len(c.ReleaseCompilations.Datastores) == 0 { // TODO check for name = default instead?
			c.ReleaseCompilations.Datastores = append(c.ReleaseCompilations.Datastores, ReleaseCompilationDatastoreConfig{
				AbstractComponentConfig: defaultServer,
			})
		}

		if len(c.Stemcells.Datastores) == 0 { // TODO check for name = default instead?
			c.Stemcells.Datastores = append(c.Stemcells.Datastores, StemcellDatastoreConfig{
				AbstractComponentConfig: defaultServer,
			})
		}
	}
}
