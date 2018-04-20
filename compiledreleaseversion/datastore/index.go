package compiledreleaseversions

type Index interface {
	Find(compiledrelease CompiledReleaseVersionRef) (CompiledReleaseVersion, error)
	List() ([]CompiledReleaseVersion, error)
}
