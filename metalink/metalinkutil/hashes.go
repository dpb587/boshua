package metalinkutil

import (
	"fmt"

	"github.com/dpb587/metalink"
)

func PreferredHash(hashes []metalink.Hash) metalink.Hash {
	for _, hashType := range []string{"sha-512", "sha-256", "sha-1", "md5"} {
		for _, hash := range hashes {
			if hash.Type == hashType {
				return hash
			}
		}
	}

	panic(fmt.Errorf("preferred hash not found"))
}
