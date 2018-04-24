package manager

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
)

type Manager struct {
	releaseVersionIndex  releaseversiondatastore.Index
	stemcellVersionIndex stemcellversiondatastore.Index
}

func NewManager(
	releaseVersionIndex releaseversiondatastore.Index,
	stemcellVersionIndex stemcellversiondatastore.Index,
) *Manager {
	return &Manager{
		releaseVersionIndex:  releaseVersionIndex,
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (rsr *Manager) Resolve(subject compiledreleaseversion.Subject) (compiledreleaseversion.ResolvedSubject, error) {
	release, err := rsr.releaseVersionIndex.Find(subject.Reference.Release)
	if err != nil {
		return compiledreleaseversion.ResolvedSubject{}, err
	}

	stemcell, err := rsr.stemcellVersionIndex.Find(subject.Reference.Stemcell)
	if err != nil {
		return compiledreleaseversion.ResolvedSubject{}, err
	}

	return compiledreleaseversion.ResolvedSubject{
		Subject:                 subject,
		ResolvedReleaseVersion:  release,
		ResolvedStemcellVersion: stemcell,
	}, nil
}
