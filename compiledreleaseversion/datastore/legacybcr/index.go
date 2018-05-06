package legacybcr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/inmemory"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger              logrus.FieldLogger
	config              Config
	repository          *git.Repository
	releaseVersionIndex releaseversiondatastore.Index
	inmemory            datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, releaseVersionIndex releaseversiondatastore.Index, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:              logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:              config,
		repository:          git.NewRepository(logger, config.RepositoryConfig),
		releaseVersionIndex: releaseVersionIndex,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.repository.Reload)

	return idx
}

func (i *index) Filter(ref compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) Store(artifact compiledreleaseversion.Artifact) error {
	return errors.New("TODO")
}

func (i *index) loader() ([]compiledreleaseversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/data/**/**/**/bcr.json", i.config.RepositoryConfig.LocalPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []compiledreleaseversion.Artifact{}

	for _, bcrPath := range paths {
		bcrBytes, err := ioutil.ReadFile(bcrPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrPath, err)
		}

		scanner := bufio.NewScanner(bytes.NewBuffer(bcrBytes))
		for scanner.Scan() {
			var record Record

			err = json.Unmarshal(scanner.Bytes(), &record)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", bcrPath, err)
			}

			meta4File := metalink.File{
				Name: fmt.Sprintf("%s-%s-on-%s-%s.tgz", record.Name, record.Version, record.Stemcell.OS, record.Stemcell.Version),
				URLs: []metalink.URL{
					{
						URL: record.Tarball.URL,
					},
				},
				Hashes: []metalink.Hash{
					{
						Type: "sha-1",
						Hash: fmt.Sprintf("%x", record.Tarball.Digest.Data()),
					},
				},
			}

			releaseRef := releaseversion.Reference{
				Name:      record.Name,
				Version:   record.Version,
				Checksums: checksum.ImmutableChecksums{record.Source.Digest},
			}

			releaseArtifacts, err := i.releaseVersionIndex.Filter(releaseRef)
			if err != nil || len(releaseArtifacts) != 1 {
				// TODO log?
				continue
			}

			releaseRef = releaseArtifacts[0].Reference

			inmemory = append(inmemory, compiledreleaseversion.New(
				releaseRef,
				osversion.Reference{
					Name:    record.Stemcell.OS,
					Version: record.Stemcell.Version,
				},
				meta4File,
				map[string]interface{}{},
			))
		}
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
