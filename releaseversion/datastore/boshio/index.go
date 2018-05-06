package boshio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type index struct {
	logger   logrus.FieldLogger
	config   Config
	inmemory datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config: config,
	}

	reloader := git.NewReloader(logger, config.RepositoryConfig)

	idx.inmemory = inmemory.New(idx.loader, reloader.Reload)

	return idx
}

func (i *index) Filter(ref releaseversion.Reference) ([]releaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]releaseversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/**/**/source.meta4", i.config.RepositoryConfig.LocalPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []releaseversion.Artifact{}

	for _, meta4Path := range paths {
		meta4Source := map[string]interface{}{
			"uri": fmt.Sprintf(
				"%s//%s",
				i.config.RepositoryConfig.Repository,
				strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.config.RepositoryConfig.LocalPath)), "/"),
			),
		}

		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", meta4Path, err)
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", meta4Path, err)
		}

		var metadataPath = fmt.Sprintf("%s/release.v1.yml", path.Dir(meta4Path))

		if _, err = os.Stat(metadataPath); err != nil {
			// TODO warn?
			continue
		}

		var metadataReleaseV1 MetadataReleaseV1

		metadataBytes, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", metadataPath, err)
		}

		err = yaml.Unmarshal(metadataBytes, &metadataReleaseV1)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", metadataPath, err)
		}

		meta4File := meta4.Files[0]

		ref := releaseversion.Reference{
			Name:      metadataReleaseV1.Name,
			Version:   metadataReleaseV1.Version,
			Checksums: metalinkutil.HashesToChecksums(meta4File.Hashes),
		}

		meta4Source["version"] = ref.Version

		inmemory = append(inmemory, releaseversion.New(
			ref,
			meta4File,
			meta4Source,
		))
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
