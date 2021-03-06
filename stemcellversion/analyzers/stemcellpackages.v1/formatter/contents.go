package formatter

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/result"
)

type Contents struct{}

func (f Contents) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		fmt.Fprintln(writer, record.Line)

		return nil
	})
}
