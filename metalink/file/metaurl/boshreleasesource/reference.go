package boshreleasesource

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
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
	// TODO privateKey support
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
	uri, err := url.Parse(o.url)
	if err != nil {
		return nil, errors.Wrap(err, "parsing url")
	}

	scheme := uri.Scheme
	pathSplit := strings.SplitN(uri.Path, "//", 2)

	if !strings.HasPrefix(scheme, "git+") {
		scheme = fmt.Sprintf("git+%s", scheme) // TODO normalize repository uris
		// return nil, fmt.Errorf("unsupported scheme: %s", scheme)
	}

	// TODO port missing
	repoURI := fmt.Sprintf("%s://%s%s", scheme[4:], uri.Hostname(), pathSplit[0])

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

	if uri.Query().Get("dev_release") == "true" {
		checkout := uri.Query().Get("checkout")
		devName := uri.Query().Get("name")
		devVersion := uri.Query().Get("version")

		if checkout != "" {
			// checkout-specific
			cmd := exec.Command(
				"git",
				fmt.Sprintf("--git-dir=%s/.git", tmpdir),
				fmt.Sprintf("--git-dir=%s", tmpdir),
				"checkout",
				checkout,
			)
			stderr := bytes.NewBuffer(nil)
			cmd.Stderr = stderr

			err := cmd.Run()
			if err != nil {
				os.RemoveAll(tmptar.Name()) // TODO ignored err

				return nil, fmt.Errorf("checkout out %s: %v: %s", checkout, err, stderr.Bytes())
			}
		}

		version := "--timestamp-version"
		if devVersion != "" {
			version = fmt.Sprintf("--version=%s", devVersion)
		}

		// build release
		cmd := exec.Command(
			"bosh",
			"create-release",
			fmt.Sprintf("--dir=%s", tmpdir),
			fmt.Sprintf("--tarball=%s", tmptar.Name()),
			"--force",
			fmt.Sprintf("--name=%s", devName),
			version,
		)
		stderr := bytes.NewBuffer(nil)
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			os.RemoveAll(tmptar.Name()) // TODO ignored err

			return nil, fmt.Errorf("creating release: %v: %s", err, stderr.Bytes())
		}
	} else {
		// build release
		cmd := exec.Command(
			"bosh",
			"create-release",
			fmt.Sprintf("--dir=%s", tmpdir),
			fmt.Sprintf("--tarball=%s", tmptar.Name()),
			fmt.Sprintf("%s/%s", tmpdir, pathSplit[1]),
		)
		stderr := bytes.NewBuffer(nil)
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			os.RemoveAll(tmptar.Name()) // TODO ignored err

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
