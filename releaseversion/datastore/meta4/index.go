package meta4

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/boshua/util"

	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string
	pullInterval       time.Duration

	inmemory   datastore.Index
	lastLoaded time.Time
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: config.Repository,
		localPath:          config.LocalPath,
		pullInterval:       config.PullInterval,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.reloader)

	return idx
}

func (i *index) List() ([]releaseversion.Subject, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Subject, error) {
	return i.inmemory.Find(ref)
}

func (i *index) reloader() (bool, error) {
	if time.Now().Sub(i.lastLoaded) < i.pullInterval {
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
		i.logger.WithField("go.error", err).Errorf("pulling repository")

		return false, fmt.Errorf("pulling repository: %v", err)
	}

	if strings.Contains(outbuf.String(), "Already up to date.") {
		i.logger.Debugf("repository already up to date")

		return false, nil
	}

	i.logger.Debugf("repository updated")

	return true, nil
}

func (i *index) loader() ([]releaseversion.Subject, error) {
	paths, err := filepath.Glob(filepath.Join(i.localPath, "*.meta4"))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

	var inmemory = []releaseversion.Subject{}

	for _, meta4Path := range paths {
		releaseversion := releaseversion.Subject{
			Reference: releaseversion.Reference{},
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

		for _, hash := range meta4.Files[0].Hashes {
			hashType, err := util.FromMetalinkHashType(hash.Type)
			if err != nil {
				continue
			}

			cs, err := checksum.CreateFromString(fmt.Sprintf("%s:%s", hashType, hash.Hash))
			if err != nil {
				continue
			}

			releaseversion.Checksums = append(releaseversion.Checksums, cs)
		}

		releaseversion.Reference.Name = path.Base(path.Dir(meta4Path))
		releaseversion.Reference.Version = meta4.Files[0].Version
		releaseversion.MetalinkSource["version"] = releaseversion.Reference.Version

		inmemory = append(inmemory, releaseversion)
	}

	return inmemory, nil
}
