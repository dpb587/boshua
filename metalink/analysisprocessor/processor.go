package analysisprocessor

import (
	"compress/gzip"
	"io"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/metalink"
	"github.com/pkg/errors"
)

func Process(artifact analysis.Artifact, callback func(io.Reader) error) error {
	r, w := io.Pipe()

	go func() {
		defer w.Close()

		err := metalink.StreamFile(artifact.MetalinkFile(), w)
		if err != nil {
			panic(errors.Wrap(err, "streaming")) // TODO panic
		}
	}()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "starting gzip")
	}

	err = callback(gzr)
	if err != nil {
		return errors.Wrap(err, "processing")
	}

	return nil
}
