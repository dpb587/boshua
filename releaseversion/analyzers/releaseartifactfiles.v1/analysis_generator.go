package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dpb587/boshua/analysis"
	filescommonoutput "github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/output"
	"github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/output"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/util/checksum/algorithm"
	"github.com/pkg/errors"
)

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

func (a *analysisGenerator) analyzeArtifact(results analysis.Writer, artifact string, reader io.Reader) error {
	gzReader, err := gzip.NewReader(reader)
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

		filestat := filescommonoutput.Result{
			Type:    string(header.Typeflag),
			Path:    path,
			Link:    header.Linkname,
			Size:    header.Size,
			Mode:    header.Mode,
			Uid:     int64(header.Uid),
			Gid:     int64(header.Gid),
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

		checksums := checksum.WritableChecksums{
			checksum.New(algorithm.MustLookupName(algorithm.MD5)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA1)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA256)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA512)),
		}

		_, err = io.Copy(checksums, tarReader)
		if err != nil {
			return errors.Wrap(err, "creating checksum")
		}

		filestat.Checksums = checksums.ImmutableChecksums()

		err = results.Write(output.Result{
			Artifact: artifact,
			Result:   filestat,
		})
		if err != nil {
			return errors.Wrap(err, "writing result")
		}
	}

	return nil
}
