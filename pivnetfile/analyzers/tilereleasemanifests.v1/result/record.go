package result

type Record struct {
	Release string     `json:"release" yaml:"release"`
	Path    string     `json:"path" yaml:"path"`
	Raw     string     `json:"raw" yaml:"raw"`
	Parsed  RecordSpec `json:"parsed" yaml:"parsed"`
}

type RecordSpec interface{}
