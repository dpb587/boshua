package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/boshua/pivnetfile/analyzers/tilearchivefiles.v1/result"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Shasum struct {
	Algorithm algorithm.Algorithm
}

func (f Shasum) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		for _, cs := range record.Result.Checksums {
			if cs.Algorithm().Name() == f.Algorithm.Name() {
				fmt.Fprintf(writer, "%x  %s\n", cs.Data(), strings.TrimPrefix(fmt.Sprintf("%s!%s", strings.Join(record.Parents, "!"), record.Result.Path), "!"))

				return nil
			}
		}

		return nil
	})
}
