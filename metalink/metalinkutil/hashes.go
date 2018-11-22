package metalinkutil

import (
	"fmt"

	"github.com/dpb587/metalink"
)

func PreferredHash(hashes []metalink.Hash) metalink.Hash {
	for _, hashType := range []metalink.HashType{
		metalink.HashTypeSHA512,
		metalink.HashTypeSHA256,
		metalink.HashTypeSHA1,
		metalink.HashTypeMD5,
	} {
		for _, hash := range hashes {
			if hash.Type == hashType {
				return hash
			}
		}
	}

	panic(fmt.Errorf("preferred hash not found"))
}
