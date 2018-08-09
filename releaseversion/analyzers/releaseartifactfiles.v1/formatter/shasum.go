package formatter

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/result"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Shasum struct {
	Algorithm algorithm.Algorithm
}

func (f Shasum) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		for _, cs := range record.Result.Checksums {
			if cs.Algorithm().Name() == f.Algorithm.Name() {
				fmt.Fprintf(writer, "%x  %s\n", cs.Data(), fmt.Sprintf("%s!%s", record.Artifact, record.Result.Path))

				return nil
			}
		}

		return nil
	})
}
