package config

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Repository    string        `yaml:"repository"`
	LocalPath     string        `yaml:"local_path"`
	PullInterval_ Duration      `yaml:"pull_interval"`
	PullInterval  time.Duration `yaml:"-"`
}

func (c *Config) Load(options map[string]interface{}) error {
	optionsBytes, err := yaml.Marshal(options)
	if err != nil {
		return fmt.Errorf("remarshaling: %v", err)
	}

	err = yaml.Unmarshal(optionsBytes, c)
	if err != nil {
		return fmt.Errorf("unmarshaling: %v", err)
	}

	err = c.validate()
	if err != nil {
		return fmt.Errorf("validating: %v", err)
	}

	c.applyDefaults()

	c.PullInterval = time.Duration(c.PullInterval_)

	return nil
}

func (c *Config) validate() error {
	if c.Repository == "" {
		return errors.New("repository must not be empty")
	}

	return nil
}

func (c *Config) applyDefaults() {
	if c.LocalPath == "" {
		hasher := sha1.New()
		hasher.Write([]byte(c.Repository))
		c.LocalPath = filepath.Join(os.TempDir(), fmt.Sprintf("boshua-%x", hex.EncodeToString(hasher.Sum(nil))))
	}

	if c.PullInterval_ == 0 {
		c.PullInterval_ = Duration(5 * time.Minute)
	}
}
