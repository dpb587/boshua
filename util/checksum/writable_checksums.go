package checksum

import "github.com/pkg/errors"

type WritableChecksums []*WritableChecksum

func (cs WritableChecksums) Write(p []byte) (int, error) {
	for _, c := range cs {
		l, err := c.Write(p)
		if err != nil {
			return l, errors.Wrap(err, "writing checksum")
		}
	}

	return len(p), nil // TODO optimistic
}

func (cs WritableChecksums) ImmutableChecksums() ImmutableChecksums {
	var res ImmutableChecksums

	for _, c := range cs {
		res = append(res, c.ImmutableChecksum())
	}

	return res
}
