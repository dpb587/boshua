package releasemanifests

type Record struct {
	Path     string     `json:"path" yaml:"path"`
	Manifest RecordSpec `json:"spec" yaml:"spec"`
}

type RecordSpec interface{}
