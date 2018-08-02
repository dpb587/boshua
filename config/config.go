package config

type Config struct {
	General   GeneralConfig              `yaml:"general,omitempty"`
	Scheduler *AbstractComponentConfig   `yaml:"scheduler,omitempty"`
	Stemcells []StemcellVersionDatastore `yaml:"stemcell_versions"`
	Releases  []AbstractComponentConfig  `yaml:"release_versions"` // TODO ReleaseDatastore
	// TODO release-specific indices for compiled release datastores
	CompiledReleases []AbstractComponentConfig `yaml:"compiled_release_versions"` // TODO ReleaseCompilationDatastore
	Analyses         []AnalysisDatastore       `yaml:"analyses"`
	Server           ServerConfig              `yaml:"server"`
}

type GeneralConfig struct {
	DefaultServer string `yaml:"default_server"`
}

type ServerConfig struct {
	Bind     string               `yaml:"bind"`
	Mount    ServerMountConfig    `yaml:"mount"`
	Redirect ServerRedirectConfig `yaml:"redirect"`
}

type ServerMountConfig struct {
	UI  string `yaml:"ui"`
	CLI string `yaml:"cli"`
}

type ServerRedirectConfig struct {
	Root string `yaml:"root"`
}

type ServerTLSConfig struct {
	CA          string `yaml:"ca"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
}

//
// type ReleaseVersionDatastore struct {
// 	AbstractComponentConfig `yaml:",inline"`
// 	Compilation             *ReleaseVersionCompilationDatastore `yaml:"compilation"`
// 	Analysis                *AnalysisDatastore                  `yaml:"analysis"`
// }
//
// type ReleaseVersionCompilationDatastore struct {
// 	AbstractComponentConfig `yaml:",inline"`
// 	Analysis                *AnalysisDatastore `yaml:"analysis"`
// }

type StemcellVersionDatastore struct {
	AbstractComponentConfig `yaml:",inline"`
	Analyses                []AnalysisDatastore `yaml:"analyses"`
}

type AnalysisDatastore struct {
	AbstractComponentConfig `yaml:",inline"`
}

type AbstractComponentConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}

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
