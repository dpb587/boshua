package boshiostemcellindex

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/inmemory"

	"github.com/sirupsen/logrus"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string

	inmemory   stemcellversions.Index
	lastLoaded time.Time
}

func New(logger logrus.FieldLogger, metalinkRepository, localPath string) stemcellversions.Index {
	idx := &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: metalinkRepository,
		localPath:          localPath,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.reloader)

	return idx
}

func (i *index) List() ([]stemcellversions.StemcellVersion, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref stemcellversions.StemcellVersionRef) (stemcellversions.StemcellVersion, error) {
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

func (i *index) loader() ([]stemcellversions.StemcellVersion, error) {
	i.lastLoaded = time.Now()

	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	i.logger.Infof("found %d entries", len(paths))

	var inmemory = []stemcellversions.StemcellVersion{}

	for _, meta4Path := range paths {
		stemcellversion := stemcellversions.StemcellVersion{
			StemcellVersionRef: stemcellversions.StemcellVersionRef{},
			MetalinkSource: map[string]interface{}{
				"uri": fmt.Sprintf("%s%s", i.metalinkRepository, strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.localPath)), "/")),
				"include_files": []string{
					"bosh-stemcell-*-warden-boshlite-ubuntu-trusty-go_agent.tgz",
				},
			},
		}

		stemcellversion.StemcellVersionRef.OS = path.Base(path.Dir(path.Dir(meta4Path)))
		stemcellversion.StemcellVersionRef.Version = path.Base(path.Dir(meta4Path))
		// stemcells are not currently recording their version :(
		// stemcellversion.MetalinkSource["version"] = stemcellversion.StemcellVersionRef.Version

		inmemory = append(inmemory, stemcellversion)
	}

	return inmemory, nil
}
