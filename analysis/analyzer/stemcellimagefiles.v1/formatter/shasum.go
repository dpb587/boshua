package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/output"
	"github.com/dpb587/boshua/checksum/algorithm"
)

type Shasum struct {
	Algorithm algorithm.Algorithm
}

func (f Shasum) Format(writer io.Writer, reader io.Reader) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		var result output.Result

		err := json.Unmarshal(s.Bytes(), &result)
		if err != nil {
			return err
		}

		for _, cs := range result.Checksums {
			if cs.Algorithm().Name() == f.Algorithm.Name() {
				fmt.Fprintf(writer, "%x  %s\n", cs.Data(), result.Path)

				break
			}
		}
	}

	return nil
}
