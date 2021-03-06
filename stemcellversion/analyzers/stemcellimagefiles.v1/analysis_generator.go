package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const Name analysis.AnalyzerName = "stemcellimagefiles.v1"

type analysisGenerator struct {
	tarball string
}

var _ analysis.AnalysisGenerator = &analysisGenerator{}

func NewAnalysis(tarball string) analysis.AnalysisGenerator {
	return &analysisGenerator{
		tarball: tarball,
	}
}

func (a *analysisGenerator) Analyze(results analysis.Writer) error {
	fh, err := os.Open(a.tarball)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}

	defer fh.Close()

	gzReader, err := gzip.NewReader(fh)
	if err != nil {
		return errors.Wrap(err, "starting gzip")
	}

	tarReader := tar.NewReader(gzReader)

	var stemcellMF map[string]interface{}

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "advancing tar")
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if path == "stemcell.MF" {
			// TODO optimistically first
			stemcellMFBytes, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return errors.Wrap(err, "reading stemcell.MF")
			}

			err = yaml.Unmarshal(stemcellMFBytes, &stemcellMF)
			if err != nil {
				return errors.Wrap(err, "unmarshaling stemcell.MF")
			}

		} else if path == "image" {
			// TODO cleanly exit for light stemcells?
			// TODO other iaases
			if strings.HasPrefix(stemcellMF["name"].(string), "bosh-aws-") {
				err = a.handleIMG(results, tarReader)
				if err != nil {
					return errors.Wrap(err, "handling raw disk")
				}
			} else if strings.HasPrefix(stemcellMF["name"].(string), "bosh-warden-") {
				err = a.handleTGZ(results, tarReader)
				if err != nil {
					return errors.Wrap(err, "handling tarball")
				}
			} else {
				return errors.New("unknown image type to handle")
			}
		}
	}

	return nil
}
