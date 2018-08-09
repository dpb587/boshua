package formatter

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/result"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Shasum struct {
	Algorithm algorithm.Algorithm
}

func (f Shasum) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		for _, cs := range record.Checksums {
			if cs.Algorithm().Name() == f.Algorithm.Name() {
				fmt.Fprintf(writer, "%x  %s\n", cs.Data(), record.Path)

				break
			}
		}

		return nil
	})
}
