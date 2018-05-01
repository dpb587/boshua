package analyzer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
)

func (a Analyzer) handleIMG(results analysis.Writer, imageReader io.Reader) error {
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

	err = filepath.Walk(mountDir, a.walkFS(results, mountDir))
	if err != nil {
		return fmt.Errorf("walking disk: %v", err)
	}

	return nil
}
