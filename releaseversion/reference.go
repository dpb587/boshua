package releaseversion

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/dpb587/boshua/util/checksum"
)

type Reference struct {
	Name      string                      `json:"name"`
	Version   string                      `json:"version"`
	Checksums checksum.ImmutableChecksums `json:"checksums"`
	URLs      []string                    `json:"urls"`
}

func (r Reference) UniqueID() string {
	var tokens []string

	// always name, version
	tokens = append(tokens, r.Name, r.Version)

	// prefer checksum, when available
	if len(r.Checksums) > 0 {
		// assume sha256 or sha512
		tokens = append(tokens, r.Checksums.Preferred().String())
	} else if len(r.URLs) > 0 {
		// assume first url is canonical
		tokens = append(tokens, r.URLs[0])
	}

	id := sha1.New()
	id.Write([]byte(strings.Join(tokens, "\n")))

	return fmt.Sprintf("%x", id.Sum(nil))
}
