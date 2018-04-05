package algorithm

import (
	md5hash "crypto/md5"
	sha1hash "crypto/sha1"
	sha256hash "crypto/sha256"
	sha512hash "crypto/sha512"
	"fmt"
)

const (
	MD5    = "md5"
	SHA1   = "sha1"
	SHA256 = "sha256"
	SHA512 = "sha512"
)

var (
	md5    Algorithm = New(MD5, md5hash.New)
	sha1   Algorithm = New(SHA1, sha1hash.New)
	sha256 Algorithm = New(SHA256, sha256hash.New)
	sha512 Algorithm = New(SHA512, sha512hash.New)
)

func New(name string, hasher Hasher) Algorithm {
	return Algorithm{
		name:   name,
		hasher: hasher,
	}
}

func LookupName(name string) (Algorithm, error) {
	switch name {
	case md5.Name():
		return md5, nil
	case sha1.Name():
		return sha1, nil
	case sha256.Name():
		return sha256, nil
	case sha512.Name():
		return sha512, nil
	}

	return Algorithm{}, fmt.Errorf("unknown algorithm: %s", name)
}

func MustLookupName(name string) Algorithm {
	algorithm, err := LookupName(name)
	if err != nil {
		panic(err)
	}

	return algorithm
}
