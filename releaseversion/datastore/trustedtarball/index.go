package trustedtarball

import (
	"fmt"
	"reflect"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	name       string
	logger     logrus.FieldLogger
	config     Config
	repository *repository.Repository
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:   name,
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config: config,
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]releaseversion.Artifact, error) {
	if !f.NameExpected || !f.URIExpected || !f.VersionExpected {
		// name, uri, version are expected, required
		return nil, nil
	}

	var nameMatch bool

	for _, nameRE := range i.config.Names {
		if nameRE.MatchString(f.Name) {
			nameMatch = true

			break
		}
	}

	if !nameMatch {
		return nil, nil
	}

	var uriMatch bool

	for _, uriRE := range i.config.URIs {
		if uriRE.MatchString(f.URI) {
			uriMatch = true

			break
		}
	}

	if !uriMatch {
		return nil, nil
	}

	var hashes []metalink.Hash

	if f.ChecksumExpected {
		cs, err := checksum.CreateFromString(f.Checksum)
		if err == nil {
			hashes = append(hashes, metalinkutil.ChecksumToHash(cs))
		}
	}

	return []releaseversion.Artifact{
		{
			Datastore: i.name,
			Name:      f.Name,
			Version:   f.Version,
			SourceTarball: metalink.File{
				Name:    fmt.Sprintf("%s-%s.tgz", f.Name, f.Version),
				Version: f.Version,
				Hashes:  hashes,
				URLs: []metalink.URL{
					{
						URL: f.URI,
					},
				},
			},
		},
	}, nil
}

func (i *index) GetLabels() ([]string, error) {
	return nil, nil
}

func (i *index) FlushCache() error {
	return nil
}
