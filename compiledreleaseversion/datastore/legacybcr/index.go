package legacybcr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
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
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string
	pullInterval       time.Duration

	releaseVersionIndex releaseversiondatastore.Index
	inmemory            datastore.Index
	lastLoaded          time.Time
}

var _ datastore.Index = &index{}

func New(config Config, releaseVersionIndex releaseversiondatastore.Index, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:              logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository:  config.Repository,
		localPath:           config.LocalPath,
		pullInterval:        config.PullInterval,
		releaseVersionIndex: releaseVersionIndex,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.reloader)

	return idx
}

func (i *index) List() ([]compiledreleaseversion.Artifact, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
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

func (i *index) loader() ([]compiledreleaseversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/data/**/**/**/bcr.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

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

			releaseArtifact, err := i.releaseVersionIndex.Find(releaseRef)
			if err == nil {
				releaseRef = releaseArtifact.Reference
			}

			inmemory = append(inmemory, compiledreleaseversion.New(
				releaseRef,
				osversion.Reference{
					OS:      record.Stemcell.OS,
					Version: record.Stemcell.Version,
				},
				meta4File,
				map[string]interface{}{},
			))
		}
	}

	return inmemory, nil
}
