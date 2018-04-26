package releaseversion

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/checksum"
)

type Reference struct {
	Name      string
	Version   string
	Checksums checksum.ImmutableChecksums
}

func (r Reference) ArtifactReference() boshua.Reference {
	return boshua.Reference{
		Context: "releaseversion",
		ID:      r.id(),
	}
}

func (r Reference) id() string {
	cs := r.Checksums.Preferred()

	h := sha1.New()
	h.Write([]byte(strings.Join(
		[]string{
			r.Name,
			r.Version,
			cs.String(),
		},
		"/",
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
