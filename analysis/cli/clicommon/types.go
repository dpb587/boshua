package clicommon

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/metalink/transfer"
)

type DownloaderGetter func () (transfer.Transfer, error)
type AnalysisLoader func() (analysis.Artifact, error)
type SubjectLoader func() (analysis.Subject, error)
