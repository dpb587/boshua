package output

type Result struct {
	Path     string     `json:"path" yaml:"path"`
	Manifest ResultSpec `json:"spec" yaml:"spec"`
}

type ResultSpec interface{}
