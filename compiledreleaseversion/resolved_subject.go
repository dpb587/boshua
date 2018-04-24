package compiledreleaseversion

import (
	"crypto/sha1"
	"fmt"

	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
)

type ResolvedSubject struct {
	Subject

	ResolvedReleaseVersion  releaseversion.Subject
	ResolvedStemcellVersion stemcellversion.Subject
}

func (s ResolvedSubject) SubjectReference() boshua.Reference {
	return boshua.Reference{
		Context: "compiledreleaseversion",
		ID:      s.id(),
	}
}

func (s ResolvedSubject) id() string {
	cs := s.ResolvedReleaseVersion.Checksums.Preferred()

	h := sha1.New()
	h.Write([]byte(fmt.Sprintf(
		"compiledreleaseversion:v1:%s:%s:%s:%s",
		s.ResolvedStemcellVersion.OS,
		s.ResolvedStemcellVersion.Version,
		s.ResolvedReleaseVersion.Name,
		cs.String(),
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
