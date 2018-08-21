package repository

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dpb587/boshua/util/marshaltypes"
)

// RepositoryConfig defines pull and push operations of a git repository.
type RepositoryConfig struct {
	// URI must define a git remote origin where the repository lives.
	URI string `yaml:"uri"`

	// Branch may define a specific branch to use. Branch must be configured if
	// push operations will be used.
	Branch string `yaml:"branch"`

	// LocalPath may
	LocalPath string `yaml:"local_path"`

	// PrivateKey is an SSH key when using private or push-enabled repositories.
	PrivateKey string `yaml:"private_key"`

	// PullInterval may define a specific interval at which the repository is
	// pulled for new data.
	PullInterval *marshaltypes.Duration `yaml:"pull_interval"`

	// SkipPull may be enabled to prevent all pull operations. This is typically
	// useful for debugging when also configuring LocalPath.
	SkipPull bool `yaml:"skip_pull"`

	// SkipPush may be enabled to prevent all push operations.
	SkipPush bool `yaml:"skip_push"`

	// AuthorName may define a custom author name for use in commits.
	AuthorName string `yaml:"author_name"`

	// AuthorEmail may define a custom author email for use in commits.
	AuthorEmail string `yaml:"author_email"`
}

func (c *RepositoryConfig) ApplyDefaults() {
	if c.PullInterval == nil {
		d := marshaltypes.Duration(5 * time.Minute)
		c.PullInterval = &d
	}

	if c.LocalPath == "" {
		hasher := sha1.New()
		hasher.Write([]byte(c.URI))
		c.LocalPath = filepath.Join(os.TempDir(), fmt.Sprintf("boshua-repository-%x", hasher.Sum(nil)))
	}

	if c.AuthorName == "" {
		c.AuthorName = "boshua" // TODO w/ version?
	}

	if c.AuthorEmail == "" {
		c.AuthorEmail = "boshua@localhost"
	}
}
