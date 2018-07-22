package task

type Type string

type Task struct {
	Type  Type
	Steps []Step
}

type Step struct {
	Name       string
	Input      map[string][]byte
	Args       []string
	Privileged bool
}
