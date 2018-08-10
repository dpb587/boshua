package boshreleasesource

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/metalink/file"
	"github.com/pkg/errors"
)

type Reference struct {
	url string
}

var _ file.Reference = Reference{}

func NewReference(url, _ string) Reference {
	return Reference{
		url: url,
	}
}

func (o Reference) Name() (string, error) {
	// TODO generated from release yml?
	return filepath.Base(o.url), nil
}

func (o Reference) Size() (uint64, error) {
	// TODO possible if this creates the release?
	return 0, fmt.Errorf("unsupported")
}

func (o Reference) Reader() (io.ReadCloser, error) {
	schemeSplit := strings.SplitN(o.url, "://", 2)

	uriSplit := strings.SplitN(schemeSplit[1], "//", 2)

	if !strings.HasPrefix(schemeSplit[0], "git+") {
		schemeSplit[0] = fmt.Sprintf("git+%s", schemeSplit[0]) // TODO normalize repository uris
		// return nil, fmt.Errorf("unsupported scheme: %s", schemeSplit[0])
	}

	repoURI := fmt.Sprintf("%s://%s", schemeSplit[0][4:], uriSplit[0])

	tmpdir, err := ioutil.TempDir("", "boshrelease-")
	if err != nil {
		return nil, errors.Wrap(err, "creating tempdir")
	}

	defer os.RemoveAll(tmpdir)

	{ // clone
		cmd := exec.Command("git", "clone", "--depth=1", repoURI, tmpdir)
		stderr := bytes.NewBuffer(nil)
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("cloning repository: %v: %s", err, stderr.Bytes())
		}
	}

	tmptar, err := ioutil.TempFile("", "boshrelease-")
	if err != nil {
		return nil, errors.Wrap(err, "creating tempfile")
	}

	{ // build release
		cmd := exec.Command("bosh", "create-release", fmt.Sprintf("--dir=%s", tmpdir), fmt.Sprintf("--tarball=%s", tmptar.Name()), fmt.Sprintf("%s/%s", tmpdir, uriSplit[1]))
		stderr := bytes.NewBuffer(nil)
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("creating release: %v: %s", err, stderr.Bytes())
		}
	}

	fh, err := os.Open(tmptar.Name())
	if err != nil {
		return nil, errors.Wrap(err, "opening release tarball")
	}

	return Reader{
		Reader: fh,
		path:   tmptar.Name(),
	}, nil
}

func (o Reference) ReaderURI() string {
	return o.url
}

func (o Reference) WriteFrom(r file.Reference, _ *pb.ProgressBar) error {
	return fmt.Errorf("unsupported")
}
