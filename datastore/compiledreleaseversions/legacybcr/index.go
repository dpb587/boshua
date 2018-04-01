package legacybcr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/inmemory"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type index struct {
	metalinkRepository string
	localPath          string

	inmemory   compiledreleaseversions.Index
	lastLoaded time.Time
}

func New(releaseVersionIndex releaseversions.Index, metalinkRepository, localPath string) compiledreleaseversions.Index {
	idx := &index{
		metalinkRepository: metalinkRepository,
		localPath:          localPath,
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
	if time.Now().Sub(i.lastLoaded) > 5*time.Minute {
		// true
	} else if !strings.HasPrefix(i.metalinkRepository, "git+") {
		return false, nil
	}

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = i.localPath

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("pulling repository: %v", err)
	}

	if strings.Contains(outbuf.String(), "Already up to date.") {
		return false, nil
	}

	return true, nil
}

func (i *index) loader() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/data/**/**/**/bcr.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

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
