package git

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dpb587/boshua/util/marshaltypes"
)

type RepositoryConfig struct {
	Repository    string                `yaml:"repository"`
	Branch        *string               `yaml:"branch"`
	LocalPath     string                `yaml:"local_path"`
	PrivateKey    *string               `yaml:"private_key"`
	PullInterval_ marshaltypes.Duration `yaml:"pull_interval"`
	PullInterval  time.Duration         `yaml:"-"`
	SkipPull      bool                  `yaml:"skip_pull"`
	SkipPush      bool                  `yaml:"skip_push"`
}

func (c *RepositoryConfig) ApplyDefaults() {
	c.PullInterval = time.Duration(c.PullInterval_)

	if c.LocalPath == "" {
		hasher := sha1.New()
		hasher.Write([]byte(c.Repository))
		c.LocalPath = filepath.Join(os.TempDir(), fmt.Sprintf("boshua-%x", hasher.Sum(nil)))
	}

	if c.PullInterval_ == 0 {
		c.PullInterval_ = marshaltypes.Duration(5 * time.Minute)
	}
}
