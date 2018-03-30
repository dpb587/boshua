package presentbcr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/inmemory"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"github.com/dpb587/metalink"
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
	if time.Now().Sub(i.lastLoaded) > time.Minute {
		return false, nil
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
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/**/**/compiled-release.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []compiledreleaseversions.CompiledReleaseVersion{}

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
		fmt.Printf("%s\n", bcrMeta4Path)

		meta4Bytes, err := ioutil.ReadFile(bcrMeta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrMeta4Path, err)
		}

		var bcrMeta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &bcrMeta4)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", bcrMeta4Path, err)
		}

		bcr := compiledreleaseversions.CompiledReleaseVersion{
			CompiledReleaseVersionRef: compiledreleaseversions.CompiledReleaseVersionRef{
				Release: releaseversions.ReleaseVersionRef{
					Name:    bcrJson.Release.Name,
					Version: bcrJson.Release.Version,
					Checksum: releaseversions.Checksum{
						Type:  bcrJson.Release.Checksums[0].Type,
						Value: bcrJson.Release.Checksums[0].Value,
					},
				},
				Stemcell: stemcellversions.StemcellVersionRef{
					OS:      bcrJson.Stemcell.OS,
					Version: bcrJson.Stemcell.Version,
				},
			},
		}

		for _, hash := range bcrMeta4.Files[0].Hashes {
			var hashType string

			if hash.Type == "sha-1" {
				hashType = "sha1"
			} else if hash.Type == "sha-256" {
				hashType = "sha256"
			} else {
				continue
			}

			bcr.TarballChecksums = append(bcr.TarballChecksums, releaseversions.Checksum{
				Type:  hashType,
				Value: hash.Hash,
			})
		}

		inmemory = append(inmemory, bcr)
	}

	return inmemory, nil
}
