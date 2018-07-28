package metalinkutil

import (
	"io/ioutil"

	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

func NewStaticArtifactLoader(path string) artifact.Loader {
	return func() (artifact.Artifact, error) {
		metalinkBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "reading metalink")
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(metalinkBytes, &meta4)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshaling metalink")
		}

		return artifact.StaticArtifact{
			StaticMetalinkFile: meta4.Files[0],
		}, nil
	}
}
