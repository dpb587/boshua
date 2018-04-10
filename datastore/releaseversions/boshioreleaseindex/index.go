package boshioreleaseindex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/inmemory"
	"github.com/dpb587/bosh-compiled-releases/util"

	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string

	inmemory   releaseversions.Index
	lastLoaded time.Time
}

func New(logger logrus.FieldLogger, metalinkRepository, localPath string) releaseversions.Index {
	idx := &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: metalinkRepository,
		localPath:          localPath,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.reloader)

	return idx
}

func (i *index) List() ([]releaseversions.ReleaseVersion, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref releaseversions.ReleaseVersionRef) (releaseversions.ReleaseVersion, error) {
	return i.inmemory.Find(ref)
}

func (i *index) reloader() (bool, error) {
	if time.Now().Sub(i.lastLoaded) < 5*time.Minute {
		return false, nil
	} else if !strings.HasPrefix(i.metalinkRepository, "git+") {
		return false, nil
	}

	i.lastLoaded = time.Now()

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = i.localPath

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Run()
	if err != nil {
		i.logger.WithField("error", err).Errorf("pulling repository")

		return false, fmt.Errorf("pulling repository: %v", err)
	}

	if strings.Contains(outbuf.String(), "Already up to date.") {
		i.logger.Debugf("repository already up to date")

		return false, nil
	}

	i.logger.Debugf("repository updated")

	return true, nil
}

func (i *index) loader() ([]releaseversions.ReleaseVersion, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/**/**/source.meta4", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

	var inmemory = []releaseversions.ReleaseVersion{}

	for _, meta4Path := range paths {
		releaseversion := releaseversions.ReleaseVersion{
			ReleaseVersionRef: releaseversions.ReleaseVersionRef{},
			Checksums:         releaseversions.Checksums{},
			MetalinkSource: map[string]interface{}{
				"uri": fmt.Sprintf("%s%s", i.metalinkRepository, strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.localPath)), "/")),
			},
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

		for _, files := range meta4.Files {
			for _, hash := range files.Hashes {
				hashType, err := util.FromMetalinkHashType(hash.Type)
				if err != nil {
					continue
				}

				releaseversion.Checksums = append(releaseversion.Checksums, releaseversions.Checksum(fmt.Sprintf("%s:%s", hashType, hash.Hash)))
			}
		}

		var metadataPath = fmt.Sprintf("%s/release.v1.yml", path.Dir(meta4Path))

		if _, err = os.Stat(metadataPath); err == nil {
			var metadataReleaseV1 MetadataReleaseV1

			metadataBytes, err := ioutil.ReadFile(metadataPath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %v", metadataPath, err)
			}

			err = yaml.Unmarshal(metadataBytes, &metadataReleaseV1)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", metadataPath, err)
			}

			releaseversion.ReleaseVersionRef.Name = metadataReleaseV1.Name
			releaseversion.ReleaseVersionRef.Version = metadataReleaseV1.Version
			releaseversion.MetalinkSource["version"] = releaseversion.ReleaseVersionRef.Version
		}

		inmemory = append(inmemory, releaseversion)
	}

	return inmemory, nil
}
