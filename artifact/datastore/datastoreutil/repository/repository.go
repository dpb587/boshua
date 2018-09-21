package repository

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
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
	logger.Infof("initialized (repository: %s; branch: %s; local: %s)", config.URI, config.Branch, config.LocalPath)

	return &Repository{
		logger: logger,
		config: config,
	}
}

func (i *Repository) Path(args ...string) string {
	return path.Join(append([]string{i.config.LocalPath}, args...)...)
}

func (i *Repository) WarmCache() bool {
	return time.Now().Sub(i.lastReloaded) < time.Duration(*i.config.PullInterval)
}

func (i *Repository) Reload() error {
	if i.config.SkipPull {
		return nil
	} else if i.WarmCache() {
		return nil
	}

	return i.ForceReload()
}

func (i *Repository) ForceReload() error {
	i.lastReloaded = time.Now()

	var args []string

	i.logger.Debug("reloading repository")

	if _, err := os.Stat(i.Path(".git")); os.IsNotExist(err) {
		args = []string{"clone", "--quiet", i.config.URI}

		if i.config.Branch != "" {
			args = append(args, "--branch", i.config.Branch)
		}

		args = append(args, ".")

		err = os.MkdirAll(i.Path(), 0700)
		if err != nil {
			return errors.Wrap(err, "mkdir local repo")
		}
	} else {
		args = []string{"pull", "--ff-only", "--quiet", i.config.URI}

		if i.config.Branch != "" {
			args = append(args, i.config.Branch)
		}
	}

	err := i.Exec(args...)
	if err != nil {
		return errors.Wrap(err, "fetching repository")
	}

	// TODO reset to handle force push?

	i.logger.Info("reloaded repository")

	return nil
}

func (i *Repository) Commit(files map[string][]byte, message string) error {
	i.logger.Debug("committing changes")

	err := i.ForceReload()
	if err != nil {
		return errors.Wrap(err, "reloading")
	}

	for path, data := range files {
		err = os.MkdirAll(filepath.Dir(i.Path(path)), 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir file dir")
		}

		err = ioutil.WriteFile(i.Path(path), data, 0644)
		if err != nil {
			return fmt.Errorf("writing file %s: %v", path, err)
		}

		err = i.Exec("add", path)
		if err != nil {
			return errors.Wrap(err, "adding file")
		}
	}

	{ // commit
		configs := map[string]string{
			"user.name":  i.config.AuthorName,
			"user.email": i.config.AuthorEmail,
		}

		for k, v := range configs {
			err := i.Exec("config", k, v)
			if err != nil {
				return errors.Wrapf(err, "setting %s", k)
			}
		}

		err := i.Exec("commit", "-m", message)
		if err != nil {
			return errors.Wrap(err, "committing")
		}
	}

	i.logger.Info("committed changes")

	if !i.config.SkipPush {
		i.logger.Debug("pushing repository")

		err := i.Exec("push", i.config.URI, fmt.Sprintf("HEAD:%s", i.config.Branch))
		if err != nil {
			return errors.Wrap(err, "pushing")
		}

		i.logger.Info("pushed repository")
	}

	return nil
}

func (r Repository) Exec(args ...string) error {
	return r.ExecCapture(os.Stderr, args...)
}

func (r Repository) ExecCapture(stdout io.Writer, args ...string) error {
	var executable = "git"

	if r.config.PrivateKey != "" && (args[0] == "clone" || args[0] == "pull" || args[0] == "push") {
		privateKey, err := ioutil.TempFile("", "git-privateKey")
		if err != nil {
			return errors.Wrap(err, "tempfile for id_rsa")
		}

		defer os.RemoveAll(privateKey.Name())

		err = os.Chmod(privateKey.Name(), 0600)
		if err != nil {
			return errors.Wrap(err, "chmod git wrapper")
		}

		err = ioutil.WriteFile(privateKey.Name(), []byte(r.config.PrivateKey), 0600)
		if err != nil {
			return errors.Wrap(err, "writing id_rsa")
		}

		executableWrapper, err := ioutil.TempFile("", "git-executable")
		if err != nil {
			return errors.Wrap(err, "tempfile for git wrapper")
		}

		defer os.RemoveAll(executableWrapper.Name())

		_, err = executableWrapper.WriteString(fmt.Sprintf(`#!/bin/bash

set -eu

mkdir -p ~/.ssh

cat > ~/.ssh/config <<EOF
StrictHostKeyChecking no
LogLevel quiet
EOF

chmod 0600 ~/.ssh/config

eval $(ssh-agent) > /dev/null

trap "kill $SSH_AGENT_PID" 0

SSH_ASKPASS=false DISPLAY= ssh-add "%s" 2>/dev/null # TODO suppresses real errors?

exec git "$@"`, privateKey.Name()))
		if err != nil {
			return errors.Wrap(err, "writing git wrapper")
		}

		err = executableWrapper.Close()
		if err != nil {
			return errors.Wrap(err, "closing tempfile")
		}

		err = os.Chmod(executableWrapper.Name(), 0500)
		if err != nil {
			return errors.Wrap(err, "chmod git wrapper")
		}

		executable = executableWrapper.Name()
	}

	// fmt.Fprintf(os.Stderr, "> %s %s\n", executable, strings.Join(args, " "))

	cmd := exec.Command(executable, args...)
	cmd.Dir = r.Path()
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
