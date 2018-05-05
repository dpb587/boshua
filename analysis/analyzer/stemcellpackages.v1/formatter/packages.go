package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellpackages.v1/output"
)

type Packages struct{}

func (f Packages) Format(writer io.Writer, reader io.Reader) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		var result output.Result

		err := json.Unmarshal(s.Bytes(), &result)
		if err != nil {
			return err
		}

		if result.Package == nil {
			continue
		}

		fmt.Fprintf(writer, "%s\t%s\n", result.Package.Name, result.Package.Version)
	}

	return nil
}
