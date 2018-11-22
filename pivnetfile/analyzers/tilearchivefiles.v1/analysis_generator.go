package analyzer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dpb587/boshua/analysis"
	filescommonresult "github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/result"
	"github.com/dpb587/boshua/pivnetfile/analyzers/tilearchivefiles.v1/result"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/util/checksum/algorithm"
	"github.com/pkg/errors"
)

type analysisGenerator struct {
	subject string
}

var _ analysis.AnalysisGenerator = &analysisGenerator{}

func NewAnalysis(subject string) analysis.AnalysisGenerator {
	return &analysisGenerator{
		subject: subject,
	}
}

func (a *analysisGenerator) Analyze(records analysis.Writer) error {
	fh, err := zip.OpenReader(a.subject)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}

	defer fh.Close()

	for _, f := range fh.File {
		sfh, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "opening embedded file")
		}

		defer sfh.Close()

		filestat := filescommonresult.Record{
			Path:    f.Name,
			Size:    int64(f.UncompressedSize64), // TODO should switch filescommonresult to uint64
			ModTime: f.Modified,
		}

		if strings.HasSuffix(f.Name, "/") {
			filestat.Type = "d"
		} else {
			filestat.Type = "f"
		}

		checksums := checksum.WritableChecksums{
			checksum.New(algorithm.MustLookupName(algorithm.MD5)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA1)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA256)),
			checksum.New(algorithm.MustLookupName(algorithm.SHA512)),
		}

		teeReader := io.TeeReader(sfh, checksums)

		err = a.analyzeFile(records, []string{}, f.Name, teeReader)
		if err != nil {
			return errors.Wrapf(err, "analyzing embedded %s", f.Name)
		}

		_, err = io.Copy(ioutil.Discard, teeReader)
		if err != nil {
			return errors.Wrap(err, "finishing file")
		}

		filestat.Checksums = checksums.ImmutableChecksums()

		err = records.Write(result.Record{
			Result:   filestat,
		})
		if err != nil {
			return errors.Wrap(err, "writing result")
		}
	}

	return nil
}

func (a *analysisGenerator) analyzeFile(records analysis.Writer, parents []string, path string, reader io.Reader) error {
	var err error

	parents = append(parents, path)

	if len(parents) > 2 {
		// TODO configurable depth
		return nil
	}

	if strings.HasSuffix(path, ".gz") || strings.HasSuffix(path, ".tgz") {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return errors.Wrap(err, "starting gzip")
		}
	}

	if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tgz") || strings.HasSuffix(path, ".tar") {
		err = a.analyzeTarFile(records, parents, tar.NewReader(reader))
		if err != nil {
			return errors.Wrap(err, "analyzing tar")
		}
	} else if strings.HasSuffix(path, ".zip") {
		// TODO
	}

	return nil
}

func (a *analysisGenerator) analyzeTarFile(records analysis.Writer, parents []string, reader *tar.Reader) error {
	for {
		header, err := reader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "advancing tar")
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		filestat := filescommonresult.Record{
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

		teeReader := io.TeeReader(reader, checksums)

		err = a.analyzeFile(records, parents, path, teeReader)
		if err != nil {
			return errors.Wrapf(err, "analyzing file %s", path)
		}

		_, err = io.Copy(ioutil.Discard, teeReader)
		if err != nil {
			return errors.Wrap(err, "finishing file")
		}

		filestat.Checksums = checksums.ImmutableChecksums()

		err = records.Write(result.Record{
			Parents: parents,
			Result:  filestat,
		})
		if err != nil {
			return errors.Wrap(err, "writing result")
		}
	}

	return nil
}
