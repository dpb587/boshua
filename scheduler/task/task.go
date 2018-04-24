package task

import "github.com/dpb587/boshua"
import "github.com/concourse/atc"

type Task interface {
	Type() string
	SubjectReference() boshua.Reference
	Config() (atc.Config, error)
}
