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

	"github.com/dpb587/boshua/datastore/compiledreleaseversions"
	"github.com/dpb587/boshua/datastore/compiledreleaseversions/inmemory"
	"github.com/dpb587/boshua/datastore/releaseversions"
	"github.com/dpb587/boshua/datastore/stemcellversions"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string
	pullInterval       time.Duration

	inmemory   compiledreleaseversions.Index
	lastLoaded time.Time
}

var _ compiledreleaseversions.Index = &index{}

func New(config Config, releaseVersionIndex releaseversions.Index, logger logrus.FieldLogger) compiledreleaseversions.Index {
	idx := &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: config.Repository,
		localPath:          config.LocalPath,
		pullInterval:       config.PullInterval,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.reloader, releaseVersionIndex)

	return idx
}

func (i *index) List() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref compiledreleaseversions.CompiledReleaseVersionRef) (compiledreleaseversions.CompiledReleaseVersion, error) {
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

func (i *index) loader() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/data/**/**/**/bcr.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

	var inmemory = []compiledreleaseversions.CompiledReleaseVersion{}

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

			inmemory = append(inmemory, compiledreleaseversions.CompiledReleaseVersion{
				CompiledReleaseVersionRef: compiledreleaseversions.CompiledReleaseVersionRef{
					Release: releaseversions.ReleaseVersionRef{
						Name:     record.Name,
						Version:  record.Version,
						Checksum: releaseversions.Checksum(fmt.Sprintf("sha1:%s", record.Source.Digest)),
					},
					Stemcell: stemcellversions.StemcellVersionRef{
						OS:      record.Stemcell.OS,
						Version: record.Stemcell.Version,
					},
				},
				TarballChecksums: releaseversions.Checksums{
					releaseversions.Checksum(fmt.Sprintf("sha1:%s", record.Tarball.Digest)),
				},
				TarballURL: record.Tarball.URL,
			})
		}
	}

	return inmemory, nil
}
