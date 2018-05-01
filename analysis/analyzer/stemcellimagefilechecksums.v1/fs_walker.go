package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/analysis"
)

func (a Analyzer) walkFS(results analysis.Writer, baseDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.Mode()&os.ModeType != 0 {
			return nil
		}

		fh, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening path: %v", err)
		}

		defer fh.Close()

		return a.createResult(results, fmt.Sprintf("/%s", strings.TrimPrefix(strings.TrimPrefix(path, baseDir), "/")), fh)
	}
}
