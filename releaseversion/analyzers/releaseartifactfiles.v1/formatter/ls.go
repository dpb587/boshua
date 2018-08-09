package formatter

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/formatter"
	"github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/result"
	"github.com/pkg/errors"
)

type Ls struct{}

func (f Ls) Format(writer io.Writer, reader io.Reader) error {
	ff := formatter.NewLs(writer)

	err := result.NewProcessor(reader, func(record result.Record) error {
		record.Result.Path = fmt.Sprintf("%s!%s", record.Artifact, record.Result.Path)

		ff.Add(record.Result)

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "processing")
	}

	ff.Flush()

	return nil
}
