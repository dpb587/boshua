package releaseversion

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/util/checksum"
)

type Reference struct {
	Name      string                      `json:"name"`
	Version   string                      `json:"version"`
	Checksums checksum.ImmutableChecksums `json:"checksums"`
	URLs      []string                    `json:"urls"`
}

var _ artifact.ArtifactReference = &Reference{}

func (r Reference) ArtifactReference() artifact.Reference {
	return artifact.Reference{
		Context: "releaseversion",
		ID:      r.id(),
	}
}

func (r Reference) ArtifactStorageDir() string {
	ref := r.ArtifactReference()

	return fmt.Sprintf(
		"%s/%s/%s/%s",
		ref.Context,
		ref.ID[0:2],
		ref.ID[2:4],
		ref.ID[4:],
	)
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
