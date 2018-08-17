package config

import (
	"time"

	"github.com/dpb587/boshua/config/types"
)

type Config struct {
	General             GeneralConfig                 `yaml:"general,omitempty"`
	Scheduler           *AbstractComponentConfig      `yaml:"scheduler,omitempty"`
	Stemcells           []StemcellDatastore           `yaml:"stemcells"`
	Releases            []ReleaseDatastore            `yaml:"releases"`
	ReleaseCompilations []ReleaseCompilationDatastore `yaml:"release_compilations"`
	Analyses            []AnalysisDatastore           `yaml:"analyses"`
	Server              ServerConfig                  `yaml:"server"`
}

type GeneralConfig struct {
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

type ReleaseDatastore struct {
	AbstractComponentConfig `yaml:",inline"`
	Compilation             *ReleaseCompilationDatastore `yaml:"compilation"`
	Analysis                *AnalysisDatastore           `yaml:"analysis"`
}

type ReleaseCompilationDatastore struct {
	AbstractComponentConfig `yaml:",inline"`
	Analysis                *AnalysisDatastore `yaml:"analysis"`
}

type StemcellDatastore struct {
	AbstractComponentConfig `yaml:",inline"`
	Analysis                *AnalysisDatastore `yaml:"analysis"`
}

type AnalysisDatastore struct {
	AbstractComponentConfig `yaml:",inline"`
}

type AbstractComponentConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}
