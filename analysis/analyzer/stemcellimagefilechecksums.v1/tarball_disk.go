package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/boshua/analysis"
)

func (a Analyzer) HandleTarballDisk(results analysis.Writer, imageReader io.Reader) error {
	gzReader, err := gzip.NewReader(imageReader)
	if err != nil {
		return fmt.Errorf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("advancing tar: %v", err)
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		err = a.checksumFile(results, strings.TrimPrefix(header.Name, "./"), tarReader)
		if err != nil {
			return fmt.Errorf("checksum file: %v", err)
		}
	}

	return nil
}
