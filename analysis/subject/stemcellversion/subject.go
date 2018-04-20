package stemcellversion

import "github.com/dpb587/boshua/stemcellversion/datastore"

type Subject struct {
	input stemcellversions.StemcellVersion
}

func (s Subject) SupportedAnalyzers() []string {
	return []string{
		"stemcellimagechecksums.v1",
		"stemcellimagefilestat.v1",
		"stemcellmanifest.v1",
	}
}

func (s Subject) Input() map[string]interface{} {
	return s.input.MetalinkSource
}
