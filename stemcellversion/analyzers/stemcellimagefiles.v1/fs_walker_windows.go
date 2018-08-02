package analyzer

import (
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
)

func (a *analysisGenerator) walkFS(results analysis.Writer, baseDir string, userMap map[int]string, groupMap map[int]string) filepath.WalkFunc {
	panic("not supported on windows")
}

func (a *analysisGenerator) loadFileNameMap(path string) (map[int]string, error) {
	panic("not supported on windows")
}
