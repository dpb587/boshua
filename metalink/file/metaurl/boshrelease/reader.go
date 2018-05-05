package boshrelease

import (
	"io"
	"os"
)

type Reader struct {
	io.Reader

	path string
}

var _ io.ReadCloser = Reader{}

func (r Reader) Close() error {
	os.Remove(r.path)

	return nil
}
