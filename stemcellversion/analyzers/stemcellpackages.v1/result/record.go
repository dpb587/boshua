package result

type Record struct {
	Line    string         `json:"line" yaml:"line"`
	Package *RecordPackage `json:"package,omitempty" yaml:"package,omitempty"`
}

type RecordPackage struct {
	Name         string `json:"name" yaml:"name"`
	Version      string `json:"version" yaml:"version"`
	Architecture string `json:"architecture" yaml:"architecture"`
	Description  string `json:"description" yaml:"description"`
}
