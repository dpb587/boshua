package result

type Record struct {
	Path   string     `json:"path" yaml:"path"`
	Raw    string     `json:"raw" yaml:"raw"`
	Parsed RecordSpec `json:"parsed" yaml:"parsed"`
}

type RecordSpec interface{}
