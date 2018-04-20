package compiledreleaseversions

type Factory interface {
	Create(provider, name string, options map[string]interface{}) (Index, error)
}
