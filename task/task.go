package task

type Task []Step

type Step struct {
	Name       string
	Input      map[string][]byte
	Args       []string
	Privileged bool
}
