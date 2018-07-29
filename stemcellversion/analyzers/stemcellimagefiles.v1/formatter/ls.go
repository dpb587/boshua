package formatter

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/formatter"
	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/output"
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

		ff.Add(result)
	}

	ff.Flush()

	return nil
}
