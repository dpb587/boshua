package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func (i *Repository) Reload() (bool, error) {
	if i.config.SkipPull {
		return false, nil
	} else if time.Now().Sub(i.lastReloaded) < i.config.PullInterval {
		return false, nil
	}

	i.lastReloaded = time.Now()

	err := i.requireClone()
	return true, err
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
		err := os.MkdirAll(filepath.Dir(filepath.Join(i.config.LocalPath, path)), 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir file dir")
		}

		err = ioutil.WriteFile(filepath.Join(i.config.LocalPath, path), data, 0644)
		if err != nil {
			return fmt.Errorf("writing file %s: %v", path, err)
		}

		cmd := exec.Command("git", "add", path)
		cmd.Dir = i.config.LocalPath

		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "adding file")
		}
	}

	{ // commit
		cmd := exec.Command("git", "commit", "--file", "-")
		cmd.Stdin = bytes.NewBuffer([]byte(message))
		cmd.Dir = i.config.LocalPath

		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "committing")
		}
	}

	if !i.config.SkipPush {
		cmd := exec.Command("git", "commit", "--file", "-")
		cmd.Stdin = bytes.NewBuffer([]byte(message))
		cmd.Dir = i.config.LocalPath

		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "pushing")
		}
	}

	return nil
}

func (r *Repository) requireClone() error {
	var execpath = "git"

	if r.config.PrivateKey != nil {
		execpath = fmt.Sprintf("%s.git", r.config.LocalPath)
		keyPath := fmt.Sprintf("%s.key", r.config.LocalPath)

		err := ioutil.WriteFile(keyPath, []byte(*r.config.PrivateKey), 0600)
		if err != nil {
			return errors.Wrap(err, "writing private key")
		}

		defer os.Remove(keyPath)

		err = ioutil.WriteFile(execpath, []byte(fmt.Sprintf(`#!/bin/bash
eval $(ssh-agent)
trap "kill $SSH_AGENT_PID" 0
set -eu
SSH_ASKPASS=false DISPLAY= ssh-add "%s"
git "$@"
`, keyPath)), 0700)
		if err != nil {
			return errors.Wrap(err, "writing git wrapper")
		}
	}

	if _, err := os.Stat(r.config.LocalPath); os.IsNotExist(err) {
		args := []string{
			"clone",
			"--single-branch",
		}

		if r.config.Branch != nil {
			args = append(args, "--branch", *r.config.Branch)
		}

		args = append(args, r.config.Repository, r.config.LocalPath)

		cmd := exec.Command(execpath, args...)
		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "cloning repository")
		}
	} else {
		args := []string{
			"pull",
			"--ff-only",
			r.config.Repository,
		}

		if r.config.Branch != nil {
			args = append(args, *r.config.Branch)
		}

		cmd := exec.Command(execpath, args...)
		cmd.Dir = r.config.LocalPath

		err = cmd.Run()

		if err != nil {
			return errors.Wrap(err, "pulling repository")
		}
	}

	return nil
}
