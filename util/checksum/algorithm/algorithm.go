package algorithm

import (
	"hash"
	"io"

	"github.com/pkg/errors"
)

type Algorithm struct {
	name   string
	hasher Hasher
}

func (a Algorithm) Name() string {
	return a.name
}

func (a Algorithm) NewHash() hash.Hash {
	return a.hasher()
}

func (a Algorithm) Hash(reader io.Reader) ([]byte, error) {
	h := a.NewHash()

	_, err := io.Copy(h, reader)
	if err != nil {
		return nil, errors.Wrap(err, "copying data")
	}

	return h.Sum(nil), nil
}
