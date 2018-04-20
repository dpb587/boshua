package analysis

type Analyzer interface {
	Analyze(Writer) error
}
