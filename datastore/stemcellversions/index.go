package stemcellversions

type Index interface {
	Find(ref StemcellVersionRef) (StemcellVersion, error)
	List() ([]StemcellVersion, error)
}
