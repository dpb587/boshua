package pivnet

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/pivotal-cf/go-pivnet"
)

type reader struct {
	client        pivnet.Client
	productName   string
	releaseID     int
	productFileID int
	extract       string

	file *os.File
}

var _ io.ReadCloser = &reader{}

func (r *reader) Close() error {
	if r.file == nil {
		return nil
	}

	err := r.file.Close()
	if err != nil {
		return err
	}

	err = os.Remove(r.file.Name())
	if err != nil {
		return errors.Wrap(err, "removing temp file")
	}

	return nil
}

func (r *reader) Read(p []byte) (int, error) {
	var err error

	if r.file == nil {
		r.file, err = ioutil.TempFile("", "pivnet-url-")
		if err != nil {
			return 0, errors.Wrap(err, "creating temp file")
		}

		err := r.client.ProductFiles.DownloadForRelease(r.file, r.productName, r.releaseID, r.productFileID, ioutil.Discard)
		if err != nil {
			return 0, errors.Wrap(err, "downloading product file")
		}

		_, err = r.file.Seek(0, 0)
		if err != nil {
			return 0, errors.Wrap(err, "rewinding downloaded file")
		}
	}

	return r.file.Read(p)
}
