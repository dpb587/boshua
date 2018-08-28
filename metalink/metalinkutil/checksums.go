package metalinkutil

import (
	"fmt"

	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

func ChecksumToHash(cs checksum.ImmutableChecksum) metalink.Hash {
	t, err := ToMetalinkHashType(cs.Algorithm().Name())
	if err != nil {
		// TODO no panic?
		panic(errors.Wrap(err, "converting hash type"))
	}

	return metalink.Hash{
		Type: t,
		Hash: fmt.Sprintf("%x", cs.Data()),
	}
}

func HashToChecksum(hash metalink.Hash) checksum.ImmutableChecksum {
	hashType, err := FromMetalinkHashType(hash.Type)
	if err != nil {
		// TODO no panic?
		panic(errors.Wrap(err, "converting hash type"))
	}

	cs, err := checksum.CreateFromString(fmt.Sprintf("%s:%s", hashType, hash.Hash))
	if err != nil {
		// TODO no panic?
		panic(errors.Wrap(err, "parsing checksum"))
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
