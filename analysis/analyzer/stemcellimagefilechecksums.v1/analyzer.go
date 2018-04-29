package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"gopkg.in/yaml.v2"
)

type Analyzer struct {
	tarball string
}

var _ analysis.Analyzer = &Analyzer{}

func New(tarball string) Analyzer {
	return Analyzer{
		tarball: tarball,
	}
}

func (a Analyzer) Analyze(results analysis.Writer) error {
	fh, err := os.Open(a.tarball)
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}

	defer fh.Close()

	gzReader, err := gzip.NewReader(fh)
	if err != nil {
		return fmt.Errorf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	var stemcellMF map[string]interface{}

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("advancing tar: %v", err)
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if path == "stemcell.MF" {
			// TODO optimistically first
			stemcellMFBytes, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("reading stemcell.MF: %v", err)
			}

			err = yaml.Unmarshal(stemcellMFBytes, &stemcellMF)
			if err != nil {
				return fmt.Errorf("unmarshaling stemcell.MF: %v", err)
			}

		} else if path == "image" {
			if strings.HasPrefix(stemcellMF["name"].(string), "bosh-aws-") {
				err = a.HandleRawDisk(results, tarReader)
				if err != nil {
					return fmt.Errorf("handling raw disk: %v", err)
				}
			} else if strings.HasPrefix(stemcellMF["name"].(string), "bosh-warden-") {
				err = a.HandleTarballDisk(results, tarReader)
				if err != nil {
					return fmt.Errorf("handling raw disk: %v", err)
				}
			} else {
				return errors.New("unknown image type to handle")
			}
		}
	}

	return nil
}
