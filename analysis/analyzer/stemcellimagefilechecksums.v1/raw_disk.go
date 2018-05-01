package analyzer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/analysis"
)

func (a Analyzer) HandleRawDisk(results analysis.Writer, imageReader io.Reader) error {
	image, err := ioutil.TempFile("", "boshua-image-")
	if err != nil {
		return fmt.Errorf("making temp file: %v", err)
	}

	//	defer os.Remove(image.Name())

	{ // extract
		cmd := exec.Command("tar", "-xzOf-", "root.img")
		cmd.Stdin = imageReader
		cmd.Stdout = image
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("extracting image: %v", err)
		}
	}

	mountDir, err := ioutil.TempDir("", "boshua-mount-")
	if err != nil {
		return fmt.Errorf("making temp dir: %v", err)
	}

	//	defer os.RemoveAll(mountDir)

	{ // mount
		cmd := exec.Command("mount", "-t", "ext4", "-o", "loop,offset=32256,ro", image.Name(), mountDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("mounting image: %v", err)
		}
	}

	err = filepath.Walk(mountDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil
		}

		fh, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening path: %v", err)
		}

		defer fh.Close()

		return a.checksumFile(results, strings.TrimPrefix(strings.TrimPrefix(path, mountDir), "/"), fh)
	})
	if err != nil {
		return fmt.Errorf("walking disk: %v", err)
	}

	return nil
}
