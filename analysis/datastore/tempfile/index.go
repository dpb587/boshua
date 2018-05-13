package tempfile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type index struct {
	inmemory map[analysis.Reference]metalink.File
}

var _ datastore.Index = &index{}

func New() datastore.Index {
	return &index{}
}

func (*index) Filter(analysis.Reference) ([]analysis.Artifact, error) {
	return nil, datastore.NoMatchErr
}

func (i *index) Store(analyzer analysis.Analyzer, subject analysis.Subject, reader io.Reader) error {
	tempfile, err := ioutil.TempFile("", "boshua-analysis-")
	if err != nil {
		return errors.Wrap(err, "creating tempfile")
	}

	fh, err := os.OpenFile(tempfile.Name(), os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "opening tempfile")
	}

	defer fh.Close()

	size, err := io.Copy(fh, reader)
	if err != nil {
		return errors.Wrap(err, "saving analysis")
	}

	i.inmemory[analysis.Reference{
		Subject:  subject,
		Analyzer: analyzer.Name(),
	}] = metalink.File{
		Name: fmt.Sprintf("%s.json", analyzer.Name()),
		Size: uint64(size),
		URLs: []metalink.URL{
			{
				URL: fmt.Sprintf("file://%s", tempfile.Name()),
			},
		},
	}

	return nil
}
