package config

import (
	"time"

	"github.com/dpb587/boshua/config/types"
)

type Config struct {
	Global              GlobalConfig              `yaml:"global,omitempty"`
	Scheduler           *AbstractComponentConfig  `yaml:"scheduler,omitempty"`
	Stemcells           StemcellsConfig           `yaml:"stemcells,omitempty"`
	Releases            ReleasesConfig            `yaml:"releases,omitempty"`
	ReleaseCompilations ReleaseCompilationConfigs `yaml:"release_compilations,omitempty"`
	Analyses            AnalysesConfig            `yaml:"analyses,omitempty"`
	Server              ServerConfig              `yaml:"server,omitempty"`
}

type GlobalConfig struct {
	DefaultServer string         `yaml:"default_server"`
	DefaultWait   time.Duration  `yaml:"default_wait"`
	LogLevel      types.LogLevel `yaml:"log_level"`
	Quiet         bool           `yaml:"quiet"`
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

type ReleasesConfig struct {
	DefaultLabels []string                 `yaml:"default_labels"`
	Datastores    []ReleaseDatastoreConfig `yaml:"datastores"`
}

type ReleaseDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`
	Compilations            *ReleaseCompilationDatastoreConfig `yaml:"compilations"`
	Analyses                *AnalysisDatastoreConfig           `yaml:"analyses"`
}

type ReleaseCompilationConfigs struct {
	Datastores []ReleaseCompilationDatastoreConfig `yaml:"datastores"`
}

type ReleaseCompilationDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`
	Analyses                *AnalysisDatastoreConfig `yaml:"analyses"`
}

type StemcellsConfig struct {
	DefaultLabels []string                  `yaml:"default_labels"`
	Datastores    []StemcellDatastoreConfig `yaml:"datastores"`
}

type StemcellDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`
	Analyses                *AnalysisDatastoreConfig `yaml:"analyses"`
}

type AnalysesConfig struct {
	Datastores []AnalysisDatastoreConfig `yaml:"datastores"`
}

type AnalysisDatastoreConfig struct {
	AbstractComponentConfig `yaml:",inline"`
}

type AbstractComponentConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}
