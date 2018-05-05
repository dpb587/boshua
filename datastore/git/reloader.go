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
	logger       logrus.FieldLogger
	repository   string
	localPath    string
	pullInterval time.Duration
	lastLoaded   time.Time
}

func NewReloader(
	logger logrus.FieldLogger,
	repository string,
	localPath string,
	pullInterval time.Duration,
) *Reloader {
	return &Reloader{
		logger:       logger,
		repository:   repository,
		localPath:    localPath,
		pullInterval: pullInterval,
	}
}

func (i *Reloader) Reload() (bool, error) {
	if time.Now().Sub(i.lastLoaded) < i.pullInterval {
		return false, nil
	} else if !strings.HasPrefix(i.repository, "git+") {
		return false, nil
	}

	i.lastLoaded = time.Now()

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = i.localPath

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
