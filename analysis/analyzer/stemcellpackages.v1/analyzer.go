package analyzer

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/stemcellpackages.v1/output"
	"github.com/pkg/errors"
)

const Name analysis.AnalyzerName = "stemcellpackages.v1"

var dpkgList = regexp.MustCompile(`^ii\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+(.+)$`)

type Analyzer struct {
	tarball string
}

var _ analysis.Analyzer = &Analyzer{}

func New(tarball string) Analyzer {
	return Analyzer{
		tarball: tarball,
	}
}

func (a Analyzer) Name() analysis.AnalyzerName {
	return Name
}

func (a Analyzer) Analyze(results analysis.Writer) error {
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

		if path == "packages.txt" || path == "stemcell_rpm_qa.txt" || path == "stemcell_dpkg_l.txt" {
			err = a.analyzePackages(results, path, tarReader)
		} else {
			continue
		}

		if err != nil {
			return fmt.Errorf("analyzing artifact %s: %v", path, err)
		}
	}

	return nil
}

func (a Analyzer) analyzePackages(results analysis.Writer, artifact string, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		result := output.Result{
			Line: scanner.Text(),
		}

		parsed := dpkgList.FindStringSubmatch(result.Line)
		if len(parsed) > 0 {
			result.Package = &output.ResultPackage{
				Name:         parsed[1],
				Version:      parsed[2],
				Architecture: parsed[3],
				Description:  parsed[4],
			}
		}

		err := results.Write(result)
		if err != nil {
			return errors.Wrap(err, "writing result")
		}
	}

	if scanner.Err() != nil {
		return errors.Wrap(scanner.Err(), "scanning packages")
	}

	return nil
}
