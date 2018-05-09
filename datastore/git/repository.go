package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Repository struct {
	logger logrus.FieldLogger
	config RepositoryConfig

	lastReloaded time.Time
}

func NewRepository(
	logger logrus.FieldLogger,
	config RepositoryConfig,
) *Repository {
	return &Repository{
		logger: logger,
		config: config,
	}
}

func (i *Repository) Reload() (bool, error) {
	if i.config.SkipPull {
		return false, nil
	} else if time.Now().Sub(i.lastReloaded) < i.config.PullInterval {
		return false, nil
	} else if !strings.HasPrefix(i.config.Repository, "git+") {
		return false, nil
	}

	i.lastReloaded = time.Now()

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = i.config.LocalPath

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("pulling repository: %v", err)
	}

	if strings.Contains(outbuf.String(), "Already up to date.") {
		i.logger.Debugf("repository already up to date")

		return false, nil
	}

	i.logger.Debugf("repository updated")

	return true, nil
}

func (i *Repository) Commit(files map[string][]byte, message string) error {
	for path, data := range files {
		err := os.MkdirAll(filepath.Dir(filepath.Join(i.config.LocalPath, path)), 0644)
		if err != nil {
			return fmt.Errorf("mkdir file dir: %v", err)
		}

		err = ioutil.WriteFile(filepath.Join(i.config.LocalPath, path), data, 0644)
		if err != nil {
			return fmt.Errorf("writing file %s: %v", path, err)
		}

		cmd := exec.Command("git", "add", path)
		cmd.Dir = i.config.LocalPath

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("adding file: %v", err)
		}
	}

	{ // commit
		cmd := exec.Command("git", "commit", "--file", "-")
		cmd.Stdin = bytes.NewBuffer([]byte(message))
		cmd.Dir = i.config.LocalPath

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("committing: %v", err)
		}
	}

	if !i.config.SkipPush {
		cmd := exec.Command("git", "commit", "--file", "-")
		cmd.Stdin = bytes.NewBuffer([]byte(message))
		cmd.Dir = i.config.LocalPath

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("pushing: %v", err)
		}
	}

	return nil
}
