package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/analysis"
	releaseartifactfilestat "github.com/dpb587/bosh-compiled-releases/analysis/releaseartifactfilestat.v1"
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

		filestat := releaseartifactfilestat.RecordFileStat{
			Type:    string(header.Typeflag),
			Path:    path,
			Link:    header.Linkname,
			Size:    header.Size,
			Mode:    header.Mode,
			Uid:     header.Uid,
			Gid:     header.Gid,
			Uname:   header.Uname,
			Gname:   header.Gname,
			ModTime: header.ModTime,
		}

		unknownTime := time.Time{}

		if header.Format == tar.FormatPAX || header.Format == tar.FormatGNU {
			if header.ChangeTime != unknownTime {
				filestat.ChangeTime = &header.ChangeTime
			}

			if header.AccessTime != unknownTime {
				filestat.AccessTime = &header.AccessTime
			}
		}

		err = results.Write(releaseartifactfilestat.Record{
			Artifact: artifact,
			Path:     path,
			Result:   filestat,
		})
		if err != nil {
			return fmt.Errorf("writing result: %v", err)
		}
	}

	return nil
}
