package config

func (c *Config) ApplyDefaults() {
	if c.Server.Bind == "" {
		c.Server.Bind = "127.0.0.1:4508"
	}

	if c.General.DefaultServer != "" {
		defaultServer := AbstractComponentConfig{
			Name: "default",
			Type: "boshua.v2",
			Options: map[string]interface{}{
				"url": c.General.DefaultServer,
			},
		}

		if c.Scheduler == nil {
			c.Scheduler = &defaultServer
		}

		if len(c.Analyses) == 0 { // TODO check for name = default instead?
			c.Analyses = append(c.Analyses, AnalysisDatastore{
				AbstractComponentConfig: defaultServer,
			})
		}

		if len(c.Releases) == 0 { // TODO check for name = default instead?
			c.Releases = append(c.Releases, defaultServer)
		}

		if len(c.CompiledReleases) == 0 { // TODO check for name = default instead?
			c.CompiledReleases = append(c.CompiledReleases, defaultServer)
		}

		if len(c.Stemcells) == 0 { // TODO check for name = default instead?
			c.Stemcells = append(c.Stemcells, StemcellVersionDatastore{
				AbstractComponentConfig: defaultServer,
			})
		}
	}
}
