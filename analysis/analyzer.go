package analysis

type AnalyzerName string

type Analyzer interface {
	Name() AnalyzerName
	Analyze(Writer) error
}
