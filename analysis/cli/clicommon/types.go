package clicommon

import "github.com/dpb587/boshua/analysis"

type AnalysisLoader func() (analysis.Artifact, error)
type SubjectLoader func() (analysis.Subject, error)
