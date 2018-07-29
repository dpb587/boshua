package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1/output"
)

type Contents struct{}

func (f Contents) Format(writer io.Writer, reader io.Reader) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		var result output.Result

		err := json.Unmarshal(s.Bytes(), &result)
		if err != nil {
			return err
		}

		fmt.Fprintf(writer, result.Raw)
	}

	return nil
}
