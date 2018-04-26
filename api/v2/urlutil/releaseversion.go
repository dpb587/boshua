package urlutil

import (
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/releaseversion"
)

func ApplyReleaseVersionRefToQuery(r *http.Request, ref releaseversion.Reference) {
	q := r.URL.Query()
	cs := ref.Checksums.Preferred()

	q.Add("release.name", ref.Name)
	q.Add("release.version", ref.Version)
	q.Add("release.checksum", cs.String())

	r.URL.RawQuery = q.Encode()
}

func ReleaseVersionRefFromParam(r *http.Request) (releaseversion.Reference, error) {
	releaseName, err := simpleQueryLookup(r, "release.name")
	if err != nil {
		return releaseversion.Reference{}, err
	}

	releaseVersion, err := simpleQueryLookup(r, "release.version")
	if err != nil {
		return releaseversion.Reference{}, err
	}

	releaseChecksumString, err := simpleQueryLookup(r, "release.checksum")
	if err != nil {
		return releaseversion.Reference{}, err
	}

	releaseChecksum, err := checksum.CreateFromString(releaseChecksumString)
	if err != nil {
		return releaseversion.Reference{}, fmt.Errorf("parameter 'release.checksum': %v", fmt.Errorf("parsing checksum: %v", err))
	}

	return releaseversion.Reference{
		Name:    releaseName,
		Version: releaseVersion,
		Checksums: checksum.ImmutableChecksums{
			releaseChecksum,
		},
	}, nil
}
