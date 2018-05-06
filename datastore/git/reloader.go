package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Reloader struct {
	logger     logrus.FieldLogger
	config     RepositoryConfig
	lastLoaded time.Time
}

func NewReloader(
	logger logrus.FieldLogger,
	config RepositoryConfig,
) *Reloader {
	return &Reloader{
		logger: logger,
		config: config,
	}
}

func (i *Reloader) Reload() (bool, error) {
	if i.config.SkipPull {
		return false, nil
	} else if time.Now().Sub(i.lastLoaded) < i.config.PullInterval {
		return false, nil
	} else if !strings.HasPrefix(i.config.Repository, "git+") {
		return false, nil
	}

	i.lastLoaded = time.Now()

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = i.config.LocalPath

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Run()
	if err != nil {
		i.logger.WithField("go.error", err).Errorf("pulling repository")

		return false, fmt.Errorf("pulling repository: %v", err)
	}

	if strings.Contains(outbuf.String(), "Already up to date.") {
		i.logger.Debugf("repository already up to date")

		return false, nil
	}

	i.logger.Debugf("repository updated")

	return true, nil
}
