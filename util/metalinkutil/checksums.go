package metalinkutil

import (
	"fmt"

	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/metalink"
)

func HashToChecksum(hash metalink.Hash) checksum.ImmutableChecksum {
	hashType, err := FromMetalinkHashType(hash.Type)
	if err != nil {
		// TODO no panic?
		panic(fmt.Errorf("converting hash type: %v", err))
	}

	cs, err := checksum.CreateFromString(fmt.Sprintf("%s:%s", hashType, hash.Hash))
	if err != nil {
		// TODO no panic?
		panic(fmt.Errorf("parsing checksum: %v", err))
	}

	return cs
}

func HashesToChecksums(hashes []metalink.Hash) checksum.ImmutableChecksums {
	var checksums checksum.ImmutableChecksums

	for _, hash := range hashes {
		checksums = append(checksums, HashToChecksum(hash))
	}

	return checksums
}
