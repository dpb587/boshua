package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefilestat.v1/output"
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
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		result := output.Result{
			Type:    string(header.Typeflag),
			Path:    fmt.Sprintf("/%s", strings.TrimPrefix(header.Name, "./")),
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
				result.ChangeTime = &header.ChangeTime
			}

			if header.AccessTime != unknownTime {
				result.AccessTime = &header.AccessTime
			}
		}

		err = results.Write(result)
		if err != nil {
			return fmt.Errorf("writing result: %v", err)
		}
	}

	return nil
}
