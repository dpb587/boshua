package output

type Result struct {
	Line    string         `json:"line" yaml:"line"`
	Package *ResultPackage `json:"package,omitempty" yaml:"package,omitempty"`
}

type ResultPackage struct {
	Name         string `json:"name" yaml:"name"`
	Version      string `json:"version" yaml:"version"`
	Architecture string `json:"architecture" yaml:"architecture"`
	Description  string `json:"description" yaml:"description"`
}
