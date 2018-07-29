package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/formatter"
	"github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/output"
)

type Ls struct{}

func (f Ls) Format(writer io.Writer, reader io.Reader) error {
	ff := formatter.NewLs(writer)

	s := bufio.NewScanner(reader)
	for s.Scan() {
		var result output.Result

		err := json.Unmarshal(s.Bytes(), &result)
		if err != nil {
			return err
		}

		// TODO prefix instead?
		result.Result.Path = fmt.Sprintf("%s!%s", result.Artifact, result.Result.Path)

		ff.Add(result.Result)
	}

	ff.Flush()

	return nil
}
