package boshmeta4releaseindex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/inmemory"

	"github.com/dpb587/metalink"
)

type index struct {
	metalinkRepository string
	localPath          string

	inmemory   releaseversions.Index
	lastLoaded time.Time
}

func New(metalinkRepository, localPath string) releaseversions.Index {
	idx := &index{
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

func (i *index) loader() ([]releaseversions.ReleaseVersion, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/releases/**/*.meta4", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

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

		for _, hash := range meta4.Files[0].Hashes {
			var hashType string

			if hash.Type == "sha-1" {
				hashType = "sha1"
			} else if hash.Type == "sha-256" {
				hashType = "sha256"
			} else {
				continue
			}

			releaseversion.Checksums = append(releaseversion.Checksums, releaseversions.Checksum{
				Type:  hashType,
				Value: hash.Hash,
			})
		}

		releaseversion.ReleaseVersionRef.Name = path.Base(path.Dir(meta4Path))
		releaseversion.ReleaseVersionRef.Version = meta4.Files[0].Version
		releaseversion.MetalinkSource["version"] = releaseversion.ReleaseVersionRef.Version

		inmemory = append(inmemory, releaseversion)
	}

	return inmemory, nil
}
