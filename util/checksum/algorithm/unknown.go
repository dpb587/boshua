package algorithm

import (
	"fmt"
	"hash"
)

func NewUnknown(name string) Algorithm {
	return Algorithm{
		name: name,
		hasher: func() hash.Hash {
			return erroringHash{
				name: name,
			}
		},
	}
}

type erroringHash struct {
	name string
}

var _ hash.Hash = erroringHash{}

func (h erroringHash) Write([]byte) (int, error) {
	return 0, fmt.Errorf(`algorithm "%s" was not known and cannot create a hash`, h.name)
}

func (erroringHash) Sum([]byte) []byte {
	return nil
}

func (erroringHash) Reset() {
	return
}

func (erroringHash) Size() int {
	return 0
}

func (erroringHash) BlockSize() int {
	return 256
}
