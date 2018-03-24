package releaseversions

type Index interface {
	Find(ref ReleaseVersionRef) (ReleaseVersion, error)
	List() ([]ReleaseVersion, error)
}
