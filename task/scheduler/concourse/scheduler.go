package concourse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/dpb587/boshua/task"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Scheduler struct {
	Target   string
	Insecure bool
	URL      string
	Team     string
	Username string
	Password string

	SecretsPath string

	needsLogin bool
	needsSync  bool
}

var _ task.Scheduler = &Scheduler{}

func (s *Scheduler) Schedule(t task.Task) error {
	pipelineName := s.pipelineName(t)

	file, err := ioutil.TempFile("", "boshua-")
	if err != nil {
		return errors.Wrap(err, "creating temp file")
	}

	defer file.Close()

	pipeline, err := t.Config()
	if err != nil {
		return errors.Wrap(err, "building pipeline")
	}

	pipelineBytes, err := yaml.Marshal(pipeline)
	if err != nil {
		return errors.Wrap(err, "marshaling pipeline")
	}

	_, err = file.Write(pipelineBytes)
	if err != nil {
		return errors.Wrap(err, "writing pipeline")
	}

	_, _, err = s.runStdin(
		bytes.NewBufferString("y\n"),
		"set-pipeline",
		"--pipeline", pipelineName,
		"--config", file.Name(),
		"--load-vars-from", s.SecretsPath,
	)
	if err != nil {
		return errors.Wrap(err, "setting pipeline")
	}

	_, _, err = s.run(
		"unpause-pipeline",
		"--pipeline", pipelineName,
	)
	if err != nil {
		return errors.Wrap(err, "unpausing pipeline")
	}

	return nil
}

func (s *Scheduler) Status(t task.Task) (task.Status, error) {
	pipelineName := s.pipelineName(t)

	stdout, stderr, err := s.run("jobs", "--pipeline", pipelineName)
	if err != nil {
		if strings.Contains(string(stderr), "error: resource not found") {
			return task.StatusUnknown, nil
		}

		return task.StatusUnknown, errors.Wrap(err, "listing jobs")
	}

	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 1 {
		return task.StatusUnknown, errors.New("listing jobs: lines missing")
	}

	fields := strings.Fields(string(stdout))
	if len(fields) != 4 {
		return task.StatusUnknown, errors.New("listing jobs: columns incorrect")
	}

	if fields[2] == "succeeded" {
		return task.StatusSucceeded, nil
	} else if fields[3] == "started" {
		return task.StatusRunning, nil
	} else if fields[2] == "aborted" {
		return task.StatusFailed, nil
	} else if fields[2] == "failed" {
		return task.StatusFailed, nil
	} else if fields[2] == "errored" {
		return task.StatusFailed, nil
	} else if fields[3] == "pending" {
		return task.StatusPending, nil
	} else if fields[2] == "n/a" && fields[3] == "n/a" {
		return task.StatusPending, nil
	}

	return task.StatusUnknown, errors.New("unrecognized pipeline state")
}

func (s *Scheduler) prepare() error {
	if s.needsLogin {
		args := []string{
			"-c", s.URL,
			"-n", s.Team,
			"-u", s.Username,
			"-p", s.Password,
		}

		if s.Insecure {
			args = append(args, "-k")
		}

		_, _, err := s.run("login", args...)
		if err != nil {
			return errors.Wrap(err, "logging in")
		}

		s.needsLogin = false
	}

	if s.needsSync {
		args := []string{
			"-c", s.URL,
			"-n", s.Team,
			"-u", s.Username,
			"-p", s.Password,
		}

		if s.Insecure {
			args = append(args, "-k")
		}

		_, _, err := s.run("sync", args...)
		if err != nil {
			return errors.Wrap(err, "syncing")
		}

		s.needsSync = false
	}

	return nil
}

func (s *Scheduler) isPrepareCommand(command string) bool {
	return command == "login" || command == "sync"
}

func (s *Scheduler) runStdin(stdin io.Reader, command string, args ...string) ([]byte, []byte, error) {
	if !s.isPrepareCommand(command) {
		err := s.prepare()
		if err != nil {
			return nil, nil, errors.Wrap(err, "preparing to run")
		}
	}

	allArgs := append([]string{"-t", s.Target, command}, args...)
	cmd := exec.Command("/usr/local/bin/fly", allArgs...)

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

		if strings.Contains(string(errbuf.Bytes()), "not authorized.") {
			s.needsLogin = true
			retryable = true
		} else if strings.Contains(string(errbuf.Bytes()), "out of sync with the target") {
			s.needsSync = true
			retryable = true
		}

		retryable = retryable && !s.isPrepareCommand(command)

		if retryable {
			return s.runStdin(bytes.NewBuffer(stdinAll), command, args...)
		}

		return outbuf.Bytes(), errbuf.Bytes(), fmt.Errorf("cli: running %#+v: %v", allArgs, err)
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func (s *Scheduler) run(cmd string, args ...string) ([]byte, []byte, error) {
	return s.runStdin(bytes.NewBuffer(nil), cmd, args...)
}

func (s *Scheduler) pipelineName(t task.Task) string {
	return fmt.Sprintf("%s:%s:%s", t.Type(), t.ArtifactReference().Context, t.ArtifactReference().ID)
}
