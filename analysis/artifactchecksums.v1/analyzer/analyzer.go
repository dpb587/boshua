package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dpb587/bosh-compiled-releases/analysis"
	artifactchecksums "github.com/dpb587/bosh-compiled-releases/analysis/artifactchecksums.v1"
	"github.com/dpb587/bosh-compiled-releases/checksum"
	"github.com/dpb587/bosh-compiled-releases/checksum/algorithm"
)

type Analyzer struct {
	tarball string
}

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
		fmt.Errorf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Errorf("advancing tar: %v", err)
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if !strings.HasSuffix(path, ".tgz") {
			continue
		}

		err = a.analyzeArtifact(results, path, tarReader)
		if err != nil {
			return fmt.Errorf("analyzing artifact %s: %v", path, err)
		}
	}

	return nil
}

func (a Analyzer) analyzeArtifact(results analysis.Writer, artifact string, reader io.Reader) error {
	gzReader, err := gzip.NewReader(reader)
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

		path := strings.TrimPrefix(header.Name, "./")

		checksums := checksum.WritableChecksums{
			checksum.New(algorithm.MustLookupName(algorithm.MD5)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA1)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA256)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA512)),
		}

		_, err = io.Copy(checksums, tarReader)
		if err != nil {
			return fmt.Errorf("creating checksum: %v", err)
		}

		err = results.Write(artifactchecksums.Record{
			Artifact: artifact,
			Path:     path,
			Result:   checksums.ImmutableChecksums(),
		})
		if err != nil {
			return fmt.Errorf("writing result: %v", err)
		}
	}

	return nil
}
