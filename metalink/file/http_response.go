package file

import (
	"io"
	"net/http"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/metalink/file"
	"github.com/pkg/errors"
)

type httpResponse struct {
	w http.ResponseWriter
}

func NewHTTPResponse(w http.ResponseWriter) file.Reference {
	return &httpResponse{
		w: w,
	}
}

var _ file.Reference = &httpResponse{}

func (f *httpResponse) Name() (string, error) {
	return "", errors.New("unsupported")
}
func (f *httpResponse) Size() (uint64, error) {
	return 0, errors.New("unsupported")
}

func (f *httpResponse) Reader() (io.ReadCloser, error) {
	return nil, errors.New("unsupported")
}

func (f *httpResponse) ReaderURI() string {
	return "unknown://unsupported"
}

func (f *httpResponse) WriteFrom(from file.Reference, _ *pb.ProgressBar) error {
	reader, err := from.Reader()
	if err != nil {
		return errors.Wrap(err, "opening from")
	}

	defer reader.Close()

	_, err = io.Copy(f.w, reader)

	return err
}
