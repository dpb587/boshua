package formatter

import (
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/formatter"
	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/result"
	"github.com/pkg/errors"
)

type Ls struct{}

func (f Ls) Format(writer io.Writer, reader io.Reader) error {
	ff := formatter.NewLs(writer)

	err := result.NewProcessor(reader, func(record result.Record) error {
		ff.Add(record)

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "processing")
	}

	ff.Flush()

	return nil
}
