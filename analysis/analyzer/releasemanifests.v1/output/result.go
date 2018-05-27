package output

type Result struct {
	Path   string     `json:"path" yaml:"path"`
	Raw    string     `json:"raw" yaml:"raw"`
	Parsed ResultSpec `json:"parsed" yaml:"parsed"`
}

type ResultSpec interface{}
