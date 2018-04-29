package analyzer

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefilechecksums.v1/output"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/checksum/algorithm"
)

func (Analyzer) checksumFile(results analysis.Writer, path string, reader io.Reader) error {
	checksums := checksum.WritableChecksums{
		checksum.New(algorithm.MustLookupName(algorithm.MD5)),
		checksum.New(algorithm.MustLookupName(algorithm.SHA1)),
		checksum.New(algorithm.MustLookupName(algorithm.SHA256)),
		checksum.New(algorithm.MustLookupName(algorithm.SHA512)),
	}

	_, err := io.Copy(checksums, reader)
	if err != nil {
		return fmt.Errorf("creating checksum: %v", err)
	}

	err = results.Write(output.Result{
		Path:   path,
		Result: checksums.ImmutableChecksums(),
	})
	if err != nil {
		return fmt.Errorf("writing result: %v", err)
	}

	return nil
}
