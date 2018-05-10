package manager

import (
	"errors"

	"github.com/dpb587/boshua/analysis"
)

type Manager struct{}

func (Manager) Require(subject analysis.Subject, analyzer string) (analysis.Artifact, error) {
	return analysis.Artifact{}, errors.New("TODO")
}
