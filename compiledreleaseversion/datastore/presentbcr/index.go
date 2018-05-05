package presentbcr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/inmemory"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string
	inmemory           datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: config.Repository,
		localPath:          config.LocalPath,
	}

	reloader := git.NewReloader(logger, config.Repository, config.LocalPath, config.PullInterval)

	idx.inmemory = inmemory.New(idx.loader, reloader.Reload)

	return idx
}

func (i *index) List() ([]compiledreleaseversion.Artifact, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]compiledreleaseversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/compiledreleaseversion/**/**/**/reference.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []compiledreleaseversion.Artifact{}

	for _, bcrJsonPath := range paths {
		bcrBytes, err := ioutil.ReadFile(bcrJsonPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrJsonPath, err)
		}

		var bcrJson Record

		err = json.Unmarshal(bcrBytes, &bcrJson)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", bcrJsonPath, err)
		}

		meta4Path := path.Join(path.Dir(bcrJsonPath), "artifact.meta4")

		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", meta4Path, err)
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", meta4Path, err)
		}

		inmemory = append(
			inmemory,
			compiledreleaseversion.New(
				releaseversion.Reference{
					Name:      bcrJson.Release.Name,
					Version:   bcrJson.Release.Version,
					Checksums: bcrJson.Release.Checksums,
				},
				osversion.Reference{
					Name:    bcrJson.OS.Name,
					Version: bcrJson.OS.Version,
				},
				meta4.Files[0],
				map[string]interface{}{
					"uri":     fmt.Sprintf("%s//%s", i.metalinkRepository, strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.localPath)), "/")),
					"version": meta4.Files[0].Version,
					"options": map[string]interface{}{
						"private_key": "((index_private_key))",
					},
				},
			),
		)
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
