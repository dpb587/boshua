package presentbcr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/inmemory"
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
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

func New(config Config, releaseVersionIndex releaseversiondatastore.Index, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: config.Repository,
		localPath:          config.LocalPath,
		pullInterval:       config.PullInterval,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.reloader, releaseVersionIndex)

	return idx
}

func (i *index) List() ([]compiledreleaseversion.Subject, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Subject, error) {
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

func (i *index) loader() ([]compiledreleaseversion.Subject, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/**/**/compiled-release.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

	var inmemory = []compiledreleaseversion.Subject{}

	for _, bcrJsonPath := range paths {
		bcrBytes, err := ioutil.ReadFile(bcrJsonPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrJsonPath, err)
		}

		var bcrJson Record

		err = json.Unmarshal(bcrBytes, &bcrJson)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", bcrJsonPath, err)
		}

		bcrMeta4Path := path.Join(path.Dir(bcrJsonPath), "compiled-release.meta4")

		meta4Bytes, err := ioutil.ReadFile(bcrMeta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrMeta4Path, err)
		}

		var bcrMeta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &bcrMeta4)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", bcrMeta4Path, err)
		}

		bcr := compiledreleaseversion.Subject{
			Reference: compiledreleaseversion.Reference{
				Release: releaseversion.Reference{
					Name:     bcrJson.Release.Name,
					Version:  bcrJson.Release.Version,
					Checksum: bcrJson.Release.Checksums[0],
				},
				Stemcell: stemcellversion.Reference{
					OS:      bcrJson.Stemcell.OS,
					Version: bcrJson.Stemcell.Version,
				},
			},
		}

		bcr.TarballPublished = bcrMeta4.Published
		bcr.TarballSize = &bcrMeta4.Files[0].Size

		for _, hash := range bcrMeta4.Files[0].Hashes {
			hashType, err := util.FromMetalinkHashType(hash.Type)
			if err != nil {
				continue
			}

			cs, err := checksum.CreateFromString(fmt.Sprintf("%s:%s", hashType, hash.Hash))
			if err != nil {
				continue
			}

			bcr.TarballChecksums = append(bcr.TarballChecksums, cs)
		}

		for _, url := range bcrMeta4.Files[0].URLs {
			bcr.TarballURL = url.URL
		}

		inmemory = append(inmemory, bcr)
	}

	return inmemory, nil
}
