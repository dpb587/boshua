package concourse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/dpb587/boshua/scheduler"
	"github.com/dpb587/boshua/scheduler/task"
)

type Runner struct {
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

func (c *Runner) Schedule(t task.Task) error {
	pipelineName := c.pipelineName(t)

	file, err := ioutil.TempFile("", "boshua-")
	if err != nil {
		return fmt.Errorf("creating temp file: %v", err)
	}

	defer file.Close()

	pipeline, err := t.Config()
	if err != nil {
		return fmt.Errorf("building pipeline: %v", err)
	}

	pipelineBytes, err := yaml.Marshal(pipeline)
	if err != nil {
		return fmt.Errorf("marshaling pipeline: %v", err)
	}

	_, err = file.Write(pipelineBytes)
	if err != nil {
		return fmt.Errorf("writing pipeline: %v", err)
	}

	_, _, err = c.runStdin(
		bytes.NewBufferString("y\n"),
		"set-pipeline",
		"--pipeline", pipelineName,
		"--config", file.Name(),
		"--load-vars-from", c.SecretsPath,
	)
	if err != nil {
		return fmt.Errorf("setting pipeline: %v", err)
	}

	_, _, err = c.run(
		"unpause-pipeline",
		"--pipeline", pipelineName,
	)
	if err != nil {
		return fmt.Errorf("unpausing pipeline: %v", err)
	}

	return nil
}

func (c *Runner) Status(t task.Task) (scheduler.Status, error) {
	pipelineName := c.pipelineName(t)

	stdout, stderr, err := c.run("jobs", "--pipeline", pipelineName)
	if err != nil {
		if strings.Contains(string(stderr), "error: resource not found") {
			return scheduler.StatusUnknown, nil
		}

		return scheduler.StatusUnknown, fmt.Errorf("listing jobs: %v", err)
	}

	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 1 {
		return scheduler.StatusUnknown, errors.New("listing jobs: lines missing")
	}

	fields := strings.Fields(string(stdout))
	if len(fields) != 4 {
		return scheduler.StatusUnknown, errors.New("listing jobs: columns incorrect")
	}

	if fields[2] == "succeeded" {
		return scheduler.StatusSucceeded, nil
	} else if fields[3] == "started" {
		return scheduler.StatusRunning, nil
	} else if fields[2] == "aborted" {
		return scheduler.StatusFailed, nil
	} else if fields[2] == "failed" {
		return scheduler.StatusFailed, nil
	} else if fields[2] == "errored" {
		return scheduler.StatusFailed, nil
	} else if fields[3] == "pending" {
		return scheduler.StatusPending, nil
	} else if fields[2] == "n/a" && fields[3] == "n/a" {
		return scheduler.StatusPending, nil
	}

	return scheduler.StatusUnknown, errors.New("unrecognized pipeline state")
}

func (c *Runner) prepare() error {
	if c.needsLogin {
		args := []string{
			"-c", c.URL,
			"-n", c.Team,
			"-u", c.Username,
			"-p", c.Password,
		}

		if c.Insecure {
			args = append(args, "-k")
		}

		_, _, err := c.run("login", args...)
		if err != nil {
			return fmt.Errorf("logging in: %v", err)
		}

		c.needsLogin = false
	}

	if c.needsSync {
		args := []string{
			"-c", c.URL,
			"-n", c.Team,
			"-u", c.Username,
			"-p", c.Password,
		}

		if c.Insecure {
			args = append(args, "-k")
		}

		_, _, err := c.run("sync", args...)
		if err != nil {
			return fmt.Errorf("syncing: %v", err)
		}

		c.needsSync = false
	}

	return nil
}

func (c *Runner) isPrepareCommand(command string) bool {
	return command == "login" || command == "sync"
}

func (c *Runner) runStdin(stdin io.Reader, command string, args ...string) ([]byte, []byte, error) {
	if !c.isPrepareCommand(command) {
		err := c.prepare()
		if err != nil {
			return nil, nil, fmt.Errorf("preparing to run: %v", err)
		}
	}

	allArgs := append([]string{"-t", c.Target, command}, args...)
	cmd := exec.Command("/usr/local/bin/fly", allArgs...)

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	stdinAll, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, nil, fmt.Errorf("buffering stdin: %v", err)
	}

	cmd.Stdin = bytes.NewBuffer(stdinAll)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err = cmd.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("starting command: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		if len(errbuf.Bytes()) > 0 {
			err = fmt.Errorf("%v - %s", err, string(errbuf.Bytes()))
		}

		if strings.Contains(string(errbuf.Bytes()), "not authorized.") {
			c.needsLogin = true
		} else if strings.Contains(string(errbuf.Bytes()), "out of sync with the target") {
			c.needsSync = true
		}

		if !c.isPrepareCommand(command) {
			return c.runStdin(bytes.NewBuffer(stdinAll), command, args...)
		}

		return outbuf.Bytes(), errbuf.Bytes(), fmt.Errorf("cli: running %#+v: %v", allArgs, err)
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func (c *Runner) run(cmd string, args ...string) ([]byte, []byte, error) {
	return c.runStdin(bytes.NewBuffer(nil), cmd, args...)
}

func (c *Runner) pipelineName(t task.Task) string {
	return fmt.Sprintf("%s:%s:%s", t.Type(), t.ArtifactReference().Context, t.ArtifactReference().ID)
}
