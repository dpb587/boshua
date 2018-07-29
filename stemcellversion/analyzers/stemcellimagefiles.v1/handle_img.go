package analyzer

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/pkg/errors"
)

func (a *analysisGenerator) handleIMG(results analysis.Writer, imageReader io.Reader) error {
	image, err := ioutil.TempFile("", "boshua-image-")
	if err != nil {
		return errors.Wrap(err, "making temp file")
	}

	defer os.Remove(image.Name())

	{ // extract
		cmd := exec.Command("tar", "-xzOf-", "root.img")
		cmd.Stdin = imageReader
		cmd.Stdout = image
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "extracting image")
		}
	}

	mountDir, err := ioutil.TempDir("", "boshua-mount-")
	if err != nil {
		return errors.Wrap(err, "making temp dir")
	}

	defer os.RemoveAll(mountDir)

	{ // mount
		cmd := exec.Command("mount", "-t", "ext4", "-o", "loop,offset=32256,ro", image.Name(), mountDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "mounting image")
		}
	}

	userMap, err := a.loadFileNameMap(filepath.Join(mountDir, "etc", "passwd"))
	if err != nil {
		return errors.Wrap(err, "loading /etc/passwd")
	}

	groupMap, err := a.loadFileNameMap(filepath.Join(mountDir, "etc", "group"))
	if err != nil {
		return errors.Wrap(err, "loading /etc/group")
	}

	err = filepath.Walk(mountDir, a.walkFS(results, mountDir, userMap, groupMap))
	if err != nil {
		return errors.Wrap(err, "walking image")
	}

	return nil
}
