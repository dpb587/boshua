package concourse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type Fly struct {
	config FlyConfig

	needsLogin bool
	needsSync  bool
}

func NewFly(config FlyConfig) *Fly {
	return &Fly{
		config: config,
	}
}

func (f *Fly) RunWithStdin(stdin io.Reader, command string, args ...string) ([]byte, []byte, error) {
	if !f.isPrepareCommand(command) {
		err := f.prepare()
		if err != nil {
			return nil, nil, errors.Wrap(err, "preparing to run")
		}
	}

	allArgs := append([]string{"-t", f.config.Target, command}, args...)
	cmd := exec.Command(f.config.Exec, allArgs...)

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	stdinAll, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, nil, errors.Wrap(err, "buffering stdin")
	}

	cmd.Stdin = bytes.NewBuffer(stdinAll)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err = cmd.Start()
	if err != nil {
		return nil, nil, errors.Wrap(err, "starting command")
	}

	err = cmd.Wait()
	if err != nil {
		if len(errbuf.Bytes()) > 0 {
			err = fmt.Errorf("%v - %s", err, string(errbuf.Bytes()))
		}

		var retryable bool

		if strings.Contains(string(errbuf.Bytes()), "unknown target") {
			f.needsLogin = true
			retryable = true
		} else if strings.Contains(string(errbuf.Bytes()), "not authorized.") {
			f.needsLogin = true
			retryable = true
		} else if strings.Contains(string(errbuf.Bytes()), "out of sync with the target") {
			f.needsSync = true
			retryable = true
		}

		retryable = retryable && !f.isPrepareCommand(command)

		if retryable {
			return f.RunWithStdin(bytes.NewBuffer(stdinAll), command, args...)
		}

		return outbuf.Bytes(), errbuf.Bytes(), fmt.Errorf("cli: running %#+v: %v", allArgs, err)
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func (f *Fly) Run(cmd string, args ...string) ([]byte, []byte, error) {
	return f.RunWithStdin(bytes.NewBuffer(nil), cmd, args...)
}

func (f *Fly) prepare() error {
	if f.needsLogin {
		args := []string{
			"-c", f.config.URL,
			"-n", f.config.Team,
			"-u", f.config.Username,
			"-p", f.config.Password,
		}

		if f.config.Insecure {
			args = append(args, "-k")
		}

		_, _, err := f.Run("login", args...)
		if err != nil {
			return errors.Wrap(err, "logging in")
		}

		f.needsLogin = false
	}

	if f.needsSync {
		_, _, err := f.Run("sync")
		if err != nil {
			return errors.Wrap(err, "syncing")
		}

		f.needsSync = false
	}

	return nil
}

func (f *Fly) isPrepareCommand(command string) bool {
	return command == "login" || command == "sync"
}
