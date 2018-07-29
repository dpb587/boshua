package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/output"
	"github.com/dpb587/boshua/util/checksum/algorithm"
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

		for _, cs := range result.Result.Checksums {
			if cs.Algorithm().Name() == f.Algorithm.Name() {
				fmt.Fprintf(writer, "%x  %s\n", cs.Data(), fmt.Sprintf("%s!%s", result.Artifact, result.Result.Path))

				break
			}
		}
	}

	return nil
}
