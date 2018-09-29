package pivnet

import (
	"io"
	"path/filepath"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/metalink/file"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pkg/errors"
)

type Reference struct {
	client        pivnet.Client
	acceptEULA    bool
	productName   string
	releaseID     int
	productFileID int
	extract       string
}

var _ file.Reference = Reference{}

func NewReference(client pivnet.Client, acceptEULA bool, productName string, releaseID, productFileID int, extract string) Reference {
	return Reference{
		client:        client,
		acceptEULA:    acceptEULA,
		productName:   productName,
		releaseID:     releaseID,
		productFileID: productFileID,
		extract:       extract,
	}
}

func (o Reference) Name() (string, error) {
	// TODO determine some other way?

	name := filepath.Base(o.extract)

	if name == "" {
		name = "pivnet.download"
	}

	return name, nil
}

func (o Reference) Size() (uint64, error) {
	// TODO
	return 0, errors.New("Unsupported")
}

func (o Reference) Reader() (io.ReadCloser, error) {
	if o.acceptEULA {
		err := o.client.EULA.Accept(o.productName, o.releaseID)
		if err != nil {
			return nil, errors.Wrap(err, "accepting pivnet product eula")
		}
	}

	return &reader{
		client:        o.client,
		productName:   o.productName,
		releaseID:     o.releaseID,
		productFileID: o.productFileID,
		extract:       o.extract,
	}, nil
}

func (o Reference) ReaderURI() string {
	// TODO pass raw URI?
	panic("lazy implementation")

	return ""
}

func (o Reference) WriteFrom(from file.Reference, progress *pb.ProgressBar) error {
	return errors.New("unsupported write operation")
}
