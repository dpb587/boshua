package template

import (
	"github.com/dpb587/metalink"
)

type templateFile metalink.File

func (tf templateFile) SHA1() string {
	for _, hash := range tf.Hashes {
		if hash.Type == "sha-1" {
			return hash.Hash
		}
	}

	panic("sha-1 missing")
}

func (tf templateFile) SHA256() string {
	for _, hash := range tf.Hashes {
		if hash.Type == "sha-256" {
			return hash.Hash
		}
	}

	panic("sha-256 missing")
}

func (tf templateFile) SHA512() string {
	for _, hash := range tf.Hashes {
		if hash.Type == "sha-256" {
			return hash.Hash
		}
	}

	panic("sha-256 missing")
}
