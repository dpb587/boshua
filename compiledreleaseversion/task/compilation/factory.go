package compilation

import "github.com/dpb587/boshua/compiledreleaseversion"

func New(subject compiledreleaseversion.ResolvedSubject) *Task {
	return &Task{
		subject: subject,
	}
}
