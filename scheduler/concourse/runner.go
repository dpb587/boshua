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

	err = c.login()
	if err != nil {
		return fmt.Errorf("logging in: %v", err)
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

	err := c.login()
	if err != nil {
		return scheduler.StatusUnknown, fmt.Errorf("logging in: %v", err)
	}

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

func (c *Runner) login() error {
	if !c.needsLogin {
		return nil
	}

	args := []string{
		"login",
		"-c", c.URL,
		"-n", c.Team,
		"-u", c.Username,
		"-p", c.Password,
	}

	if c.Insecure {
		args = append(args, "-k")
	}

	_, _, err := c.run(args...)
	if err != nil {
		return fmt.Errorf("logging in: %v", err)
	}

	c.needsLogin = false

	return nil
}

func (c *Runner) runStdin(stdin io.Reader, args ...string) ([]byte, []byte, error) {
	allArgs := append([]string{"-t", c.Target}, args...)
	cmd := exec.Command("/usr/local/bin/fly", allArgs...)

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	cmd.Stdin = stdin
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Start()
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
		}

		// TODO fly sync "out of sync with the target"

		return outbuf.Bytes(), errbuf.Bytes(), fmt.Errorf("cli: running %#+v: %v", allArgs, err)
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func (c *Runner) run(args ...string) ([]byte, []byte, error) {
	return c.runStdin(bytes.NewBuffer(nil), args...)
}

func (c *Runner) pipelineName(t task.Task) string {
	return fmt.Sprintf("%s:%s:%s", t.Type(), t.ArtifactReference().Context, t.ArtifactReference().ID)
}
