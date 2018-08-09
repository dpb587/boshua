package formatter

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/result"
)

type Packages struct{}

func (f Packages) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		if record.Package == nil {
			return nil
		}

		fmt.Fprintf(writer, "%s\t%s\n", record.Package.Name, record.Package.Version)

		return nil
	})
}
