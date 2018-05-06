package releaseversion

import (
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/releaseversion"
)

type Reference struct {
	Name      string                      `json:"name"`
	Version   string                      `json:"version"`
	Checksum  *checksum.ImmutableChecksum `json:"checksum,omitempty"`
	Checksums checksum.ImmutableChecksums `json:"checksums,omitempty"`
}

func FromReference(ref releaseversion.Reference) Reference {
	return Reference{
		Name:      ref.Name,
		Version:   ref.Version,
		Checksums: ref.Checksums,
	}
}
