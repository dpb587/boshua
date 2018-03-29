package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type Status string

var (
	StatusUnknown   Status = "unknown"
	StatusPending   Status = "pending"
	StatusAborted   Status = "aborted"
	StatusFailed    Status = "failed"
	StatusStarted   Status = "started"
	StatusSucceeded Status = "succeeded"
)

type Compiler struct {
	Target   string
	Insecure bool
	URL      string
	Team     string
	Username string
	Password string

	PipelinePath string
	SecretsPath  string
}

func (c Compiler) Schedule(release releaseversions.ReleaseVersion, stemcell stemcellversions.StemcellVersion) error {
	pipelineName := c.pipelineName(release, stemcell)

	file, err := ioutil.TempFile("", "bcr-")
	if err != nil {
		return fmt.Errorf("creating temp file: %v", err)
	}

	defer file.Close()

	contextBytes, err := json.Marshal(map[string]interface{}{
		"release": map[string]interface{}{
			"name":      release.Name,
			"version":   release.Version,
			"checksums": release.Checksums,
		},
		"stemcell": map[string]interface{}{
			"os":      stemcell.OS,
			"version": stemcell.Version,
		},
	})
	if err != nil {
		return fmt.Errorf("marshalling context: %v", err)
	}

	varsBytes, err := yaml.Marshal(map[string]interface{}{
		"release_version": release.Version,
		"release_source":  release.MetalinkSource,
		"stemcell_source": stemcell.MetalinkSource,
		"index_storage":   fmt.Sprintf("%s/%s/%s-%s/%s", release.Name, release.Version, stemcell.OS, stemcell.Version, release.Checksums.Preferred()),
		"index_context":   string(contextBytes),
	})
	if err != nil {
		return fmt.Errorf("marshalling vars: %v", err)
	}

	_, err = file.Write(varsBytes)
	if err != nil {
		return fmt.Errorf("writing temp file: %v", err)
	}

	_, _, err = c.runStdin(
		bytes.NewBufferString("y\n"),
		"set-pipeline",
		"--pipeline", pipelineName,
		"--config", c.PipelinePath,
		"--load-vars-from", c.SecretsPath,
		"--load-vars-from", file.Name(),
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

func (c Compiler) Status(release releaseversions.ReleaseVersion, stemcell stemcellversions.StemcellVersion) (Status, error) {
	stdout, stderr, err := c.run(
		"jobs",
		"--pipeline", c.pipelineName(release, stemcell),
	)
	if err != nil {
		if strings.Contains(string(stderr), "error: resource not found") {
			return StatusUnknown, nil
		}

		return StatusUnknown, fmt.Errorf("listing jobs: %v", err)
	}

	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 1 {
		return StatusUnknown, errors.New("listing jobs: lines missing")
	}

	fields := strings.Fields(string(stdout))
	if len(fields) != 4 {
		return StatusUnknown, errors.New("listing jobs: columns incorrect")
	}

	if fields[2] == "succeeded" {
		return StatusSucceeded, nil
	} else if fields[2] == "aborted" {
		return StatusAborted, nil
	} else if fields[2] == "failed" {
		return StatusFailed, nil
	} else if fields[2] == "errored" {
		return StatusFailed, nil
	} else if fields[3] == "pending" {
		return StatusPending, nil
	} else if fields[2] == "n/a" && fields[3] == "n/a" {
		return StatusPending, nil
	} else if fields[3] == "started" {
		return StatusPending, nil
	}

	return StatusUnknown, nil
}

func (c Compiler) pipelineName(release releaseversions.ReleaseVersion, stemcell stemcellversions.StemcellVersion) string {
	return fmt.Sprintf(
		"bcr:%s:%s-%s-on-%s-stemcell-%s",
		release.Checksums.Preferred(),
		release.Name,
		release.Version,
		stemcell.OS,
		stemcell.Version,
	)
}

func (c Compiler) login() error {
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

	return nil
}

func (c Compiler) runStdin(stdin io.Reader, args ...string) ([]byte, []byte, error) {
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
		return outbuf.Bytes(), errbuf.Bytes(), fmt.Errorf("cli: running %#+v: %v", allArgs, err)
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func (c Compiler) run(args ...string) ([]byte, []byte, error) {
	return c.runStdin(bytes.NewBuffer(nil), args...)
}
