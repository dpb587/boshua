package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dpb587/boshua/util/marshaltypes"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type RepositoryConfig struct {
	Repository    string                `yaml:"repository"`
	LocalPath     string                `yaml:"local_path"`
	PullInterval_ marshaltypes.Duration `yaml:"pull_interval"`
	PullInterval  time.Duration         `yaml:"-"`
	SkipPull      bool                  `yaml:"skip_pull"`
	SkipPush      bool                  `yaml:"skip_push"`
}

func (c *RepositoryConfig) Load(options map[string]interface{}) error {
	optionsBytes, err := yaml.Marshal(options)
	if err != nil {
		return errors.Wrap(err, "remarshaling")
	}

	err = yaml.Unmarshal(optionsBytes, c)
	if err != nil {
		return errors.Wrap(err, "unmarshaling")
	}

	err = c.validate()
	if err != nil {
		return errors.Wrap(err, "validating")
	}

	c.applyDefaults()

	c.PullInterval = time.Duration(c.PullInterval_)

	// TODO forceful dev
	c.SkipPush = true

	return nil
}

func (c *RepositoryConfig) validate() error {
	if c.Repository == "" {
		return errors.New("repository must not be empty")
	}

	return nil
}

func (c *RepositoryConfig) applyDefaults() {
	if c.LocalPath == "" {
		hasher := sha1.New()
		hasher.Write([]byte(c.Repository))
		c.LocalPath = filepath.Join(os.TempDir(), fmt.Sprintf("boshua-%x", hex.EncodeToString(hasher.Sum(nil))))
	}

	if c.PullInterval_ == 0 {
		c.PullInterval_ = marshaltypes.Duration(5 * time.Minute)
	}
}
