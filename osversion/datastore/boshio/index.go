package boshio

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

	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/osversion/datastore/inmemory"
	"github.com/dpb587/metalink"

	"github.com/sirupsen/logrus"
)

const compileMatch = "bosh-stemcell-*-warden-boshlite-ubuntu-trusty-go_agent.tgz"

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

func (i *index) List() ([]osversion.Artifact, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref osversion.Reference) (osversion.Artifact, error) {
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

func (i *index) loader() ([]osversion.Artifact, error) {
	i.lastLoaded = time.Now()

	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

	var inmemory = []osversion.Artifact{}

	for _, meta4Path := range paths {
		ref := osversion.Reference{
			Name:    path.Base(path.Dir(path.Dir(meta4Path))),
			Version: path.Base(path.Dir(meta4Path)),
		}

		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading metalink: %v", err)
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling metalink: %v", err)
		}

		var meta4File *metalink.File

		for _, file := range meta4.Files {
			if matched, _ := filepath.Match(compileMatch, file.Name); matched {
				meta4File = &file

				break
			}
		}

		if meta4File == nil {
			// TODO log warning?
			continue
		}

		inmemory = append(
			inmemory,
			osversion.New(
				ref,
				*meta4File,
				map[string]interface{}{
					"uri": fmt.Sprintf("%s%s", i.metalinkRepository, strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.localPath)), "/")),
					"include_files": []string{
						compileMatch,
					},
				},
			),
		)
	}

	return inmemory, nil
}
