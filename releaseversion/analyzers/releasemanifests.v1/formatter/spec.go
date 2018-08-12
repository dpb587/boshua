package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/result"
)

type Spec struct {
	ReleaseOnly bool
	Jobs        []string
}

func (f Spec) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		trimmedPath := strings.TrimPrefix(record.Path, "./")

		if f.ReleaseOnly {
			if trimmedPath != "release.MF" {
				return nil
			}
		}

		if len(f.Jobs) > 0 {
			var found bool

			for _, job := range f.Jobs {
				if trimmedPath == fmt.Sprintf("jobs/%s.tgz", job) {
					found = true

					break
				}
			}

			if !found {
				return nil
			}
		}

		raw := record.Raw

		// TODO fix logic or flags
		if len(f.Jobs) > 1 || (len(f.Jobs) == 0 && !f.ReleaseOnly) {
			// only adjust raw if multiple results may be shown
			raw = strings.TrimSpace(strings.TrimPrefix(raw, "---\n"))
			raw = fmt.Sprintf("---\n# %s\n%s\n", record.Path, raw)
		}

		fmt.Fprintf(writer, "%s\n", raw)

		return nil
	})
}
