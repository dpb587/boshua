package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/output"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/checksum/algorithm"
)

func (a Analyzer) handleTGZ(results analysis.Writer, imageReader io.Reader) error {
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
		}

		result := output.Result{
			Type:    string(header.Typeflag),
			Path:    fmt.Sprintf("/%s", strings.TrimPrefix(header.Name, "./")),
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
				result.ChangeTime = &header.ChangeTime
			}

			if header.AccessTime != unknownTime {
				result.AccessTime = &header.AccessTime
			}
		}

		if header.Typeflag == tar.TypeReg {
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

			result.Checksums = checksums.ImmutableChecksums()
		}

		err = results.Write(result)
		if err != nil {
			return fmt.Errorf("writing result: %v", err)
		}
	}

	return nil
}
