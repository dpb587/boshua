package releaseversion

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/checksum"
)

type Subject struct {
	Reference

	Checksums      checksum.ImmutableChecksums
	MetalinkSource map[string]interface{}
}

var _ boshua.Subject = &Subject{}

func (s Subject) SubjectReference() boshua.Reference {
	return boshua.Reference{
		Context: "releaseversion",
		ID:      s.id(),
	}
}

func (s Subject) id() string {
	cs := s.Checksums.Preferred()

	h := sha1.New()
	h.Write([]byte(strings.Join(
		[]string{
			s.Reference.Name,
			s.Reference.Version,
			cs.String(),
		},
		"/",
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
