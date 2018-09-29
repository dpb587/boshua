package config

import (
	"time"

	"github.com/dpb587/boshua/config/types"
)

// Config represents the standard YAML configuration file.
type Config struct {
	RawConfig func() ([]byte, error) `yaml:"-"`

	// Global defines common options. These can usually be specified via global
	// CLI options which take precedent.
	Global GlobalConfig `yaml:"global,omitempty"`

	// Scheduler defines how to execute tasks.
	Scheduler *SchedulerConfig `yaml:"scheduler,omitempty"`

	// Stemcells defines stemcell datastores and settings.
	Stemcells StemcellsConfig `yaml:"stemcells,omitempty"`

	// Releases defines release datastores and settings.
	Releases ReleasesConfig `yaml:"releases,omitempty"`

	// ReleaseCompilations defines release compilation datastores and settings.
	ReleaseCompilations ReleaseCompilationConfigs `yaml:"release_compilations,omitempty"`

	// PivnetFiles defines pivnet file datastores and settings.
	PivnetFiles PivnetFilesConfig `yaml:"pivnet_files,omitempty"`

	// Analyses defines analysis datastores and settings.
	Analyses AnalysesConfig `yaml:"analyses,omitempty"`

	// Server defines how a local API server should run.
	Server ServerConfig `yaml:"server,omitempty"`

	// Downloaders defines how artifacts should be downloaded.
	Downloaders DownloadersConfig `yaml:"downloaders,omitempty"`
}

// GlobalConfig defines common options.
type GlobalConfig struct {
	// DefaultServer defines a default remote API server which will be injected
	// for empty stemcells, releases, release compilations, analyses, and
	// scheduler if no other providers are configured.
	DefaultServer string `yaml:"default_server"`

	// DefaultWait defines how long to wait for asynchronous operations.
	DefaultWait time.Duration `yaml:"default_wait"`

	// LogLevel defines the log level for messages sent to STDERR.
	LogLevel types.LogLevel `yaml:"log_level"`

	// Quiet defines whether informational or progress information should be
	// suppressed.
	Quiet bool `yaml:"quiet"`
}

type DownloadersConfig struct {
	DisableDefaultHandlers bool                      `yaml:"disable_default_handlers"`
	URLHandlers            []DownloaderHandlerConfig `yaml:"url_handlers"`
	// MetaURLHandlers        []DownloaderHandlerConfig `yaml:"meta_url_handlers"`
}

// ServerConfig defines local API server options.
type ServerConfig struct {
	// Bind defines the host and port to run on.
	Bind string `yaml:"bind"`

	// Mount defines paths for customized endpoints.
	Mount ServerMountConfig `yaml:"mount"`

	// Redirect defines redirect rules.
	Redirect ServerRedirectConfig `yaml:"redirect"`
}

// ServerMountConfig defines paths for customzied endpoints.
type ServerMountConfig struct {
	// CLI defines a directory to serve binaries from `/cli/`.
	CLI string `yaml:"cli"`

	// UI defines a local path to serve from `/ui/`.
	UI string `yaml:"ui"`
}

// ServerRedirectConfig defines redirect rules.
type ServerRedirectConfig struct {
	// Root defines a target path for clients accessing `/`.
	Root string `yaml:"root"`
}

// PivnetFilesConfig defines pivnet file datastores and settings.
type PivnetFilesConfig struct {
	// Datastores defines a list of release datastores.
	Datastores []PivnetFileDatastoreConfig `yaml:"datastores"`
}

// ReleaseDatastoreConfig defines a release datastore.
type PivnetFileDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`

	// Analyses defines an explicit inline or reference to an analysis datastore.
	AnalysisDatastore *AnalysisDatastoreConfig `yaml:"analysis_datastore"`
}

// ReleasesConfig defines release datastores and settings.
type ReleasesConfig struct {
	// Datastores defines a list of release datastores.
	Datastores []ReleaseDatastoreConfig `yaml:"datastores"`
}

// ReleaseDatastoreConfig defines a release datastore.
type ReleaseDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`

	// Compilations defines an explicit inline or reference to a compilation
	// datastore.
	CompilationDatastore *ReleaseCompilationDatastoreConfig `yaml:"compilation_datastore"`

	// Analyses defines an explicit inline or reference to an analysis datastore.
	AnalysisDatastore *AnalysisDatastoreConfig `yaml:"analysis_datastore"`
}

// ReleaseCompilationConfigs defines release compilation datastores and
// settings.
type ReleaseCompilationConfigs struct {
	// Datastores defines a list of release compilation datastores.
	Datastores []ReleaseCompilationDatastoreConfig `yaml:"datastores"`
}

// ReleaseCompilationDatastoreConfig defines a release compilation datastore.
type ReleaseCompilationDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`

	// Analyses defines an explicit inline or reference to an analysis datastore.
	AnalysisDatastore *AnalysisDatastoreConfig `yaml:"analysis_datastore"`
}

// StemcellsConfig defines stemcell datastores and settings.
type StemcellsConfig struct {
	// Datastores defines a list of stemcell datastores.
	Datastores []StemcellDatastoreConfig `yaml:"datastores"`
}

type StemcellDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`
	AnalysisDatastore       *AnalysisDatastoreConfig `yaml:"analysis_datastore"`
}

type AnalysesConfig struct {
	Datastores []AnalysisDatastoreConfig `yaml:"datastores"`
}

type AnalysisDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`
}

type SchedulerConfig struct {
	AbstractComponentConfig `yaml:",inline"`
	NoWait                  bool `yaml:"no_wait"`
}

type DownloaderHandlerConfig struct {
	AbstractComponentConfig `yaml:",inline"`
	Include                 types.RegexpList `yaml:"include"`
	Exclude                 types.RegexpList `yaml:"exclude"`
}

type AbstractComponentConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}
