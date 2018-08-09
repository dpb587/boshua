package file

import (
	"io"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/metalink/file"
	"github.com/pkg/errors"
)

type writer struct {
	w io.WriteCloser
}

func NewWriter(w io.WriteCloser) file.Reference {
	return &writer{
		w: w,
	}
}

var _ file.Reference = &writer{}

func (f *writer) Name() (string, error) {
	return "", errors.New("unsupported")
}
func (f *writer) Size() (uint64, error) {
	return 0, errors.New("unsupported")
}

func (f *writer) Reader() (io.ReadCloser, error) {
	return nil, errors.New("unsupported")
}

func (f *writer) ReaderURI() string {
	return "unknown://unsupported"
}

func (f *writer) WriteFrom(from file.Reference, _ *pb.ProgressBar) error {
	defer f.w.Close()

	reader, err := from.Reader()
	if err != nil {
		return errors.Wrap(err, "opening from")
	}

	defer reader.Close()

	_, err = io.Copy(f.w, reader)

	return err
}
