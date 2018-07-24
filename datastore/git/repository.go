package git

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
	return &Repository{
		logger: logger,
		config: config,
	}
}

func (i *Repository) Path(args ...string) string {
	return path.Join(append([]string{i.config.LocalPath}, args...)...)
}

func (i *Repository) Reload() error {
	if i.config.SkipPull {
		return nil
	} else if time.Now().Sub(i.lastReloaded) < i.config.PullInterval {
		return nil
	}

	return i.ForceReload()
}

func (i *Repository) ForceReload() error {
	i.lastReloaded = time.Now()
	var args []string

	if _, err := os.Stat(i.Path(".git")); os.IsNotExist(err) {
		args = []string{"clone", "--quiet", i.config.Repository}

		if i.config.Branch != "" {
			args = append(args, "--branch", i.config.Branch)
		}

		args = append(args, ".")

		err = os.MkdirAll(i.Path(), 0700)
		if err != nil {
			return errors.Wrap(err, "mkdir local repo")
		}
	} else {
		args = []string{"pull", "--ff-only", "--quiet", i.config.Repository}

		if i.config.Branch != "" {
			args = append(args, i.config.Branch)
		}
	}

	err := i.run(args...)

	if err != nil {
		return errors.Wrap(err, "fetching repository")
	}

	// TODO reset to handle force push?

	return nil
	// cmd := exec.Command("git")
	//
	// outbuf := bytes.NewBuffer(nil)
	// errbuf := bytes.NewBuffer(nil)
	//
	// cmd.Stdout = outbuf
	// cmd.Stderr = errbuf
	//
	// if _, err := os.Stat(i.config.LocalPath); os.IsNotExist(err) {
	// 	cmd.Args = []string{"clone", strings.TrimPrefix(i.config.Repository, "git+"), i.config.LocalPath}
	// } else {
	// 	cmd.Dir = i.config.LocalPath
	// 	cmd.Args = []string{"pull", "--ff-only"}
	// }

	// err := cmd.Run()
	// if err != nil {
	// 	return false, errors.Wrap(err, "pulling repository")
	// }
	//
	// if strings.Contains(outbuf.String(), "Already up to date.") {
	// 	i.logger.Debugf("repository already up to date")
	//
	// 	return false, nil
	// }
	//
	// i.logger.Debugf("repository updated")
	//
	// return true, nil
}

func (i *Repository) Commit(files map[string][]byte, message string) error {
	for path, data := range files {
		err := os.MkdirAll(filepath.Dir(i.Path(path)), 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir file dir")
		}

		err = ioutil.WriteFile(i.Path(path), data, 0644)
		if err != nil {
			return fmt.Errorf("writing file %s: %v", path, err)
		}

		err = i.run("add", path)
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
			err := i.run("config", k, v)
			if err != nil {
				return errors.Wrapf(err, "setting %s", k)
			}
		}

		err := i.run("commit", "-m", message)
		if err != nil {
			return errors.Wrap(err, "committing")
		}
	}

	if !i.config.SkipPush {
		err := i.run("push", i.config.Repository, fmt.Sprintf("HEAD:%s", i.config.Branch))
		if err != nil {
			return errors.Wrap(err, "pushing")
		}
	}

	return nil
}

func (r Repository) run(args ...string) error {
	return r.runRaw(os.Stderr, args...)
}

func (r Repository) runRaw(stdout io.Writer, args ...string) error {
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
