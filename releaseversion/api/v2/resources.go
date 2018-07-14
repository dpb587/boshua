package v2

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type GETIndexResponse struct {
	Data []Reference `json:"data"`
}

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

type InfoResponse struct {
	Data InfoResponseData `json:"data"`
}

type InfoResponseData struct {
	Reference Reference     `json:"reference"`
	Artifact  metalink.File `json:"file"`
}
