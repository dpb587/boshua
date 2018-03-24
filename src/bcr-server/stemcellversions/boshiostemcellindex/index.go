package boshiostemcellindex

import (
	"bcr-server/stemcellversions"
	"bcr-server/stemcellversions/inmemory"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type index struct {
	metalinkRepository string
	localPath          string

	inmemory stemcellversions.Index
}

func New(metalinkRepository, localPath string) stemcellversions.Index {
	idx := &index{
		metalinkRepository: metalinkRepository,
		localPath:          localPath,
	}

	idx.inmemory = inmemory.New(idx.loader)

	return idx
}

func (i *index) List() ([]stemcellversions.StemcellVersion, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref stemcellversions.StemcellVersionRef) (stemcellversions.StemcellVersion, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]stemcellversions.StemcellVersion, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []stemcellversions.StemcellVersion{}

	for _, meta4Path := range paths {
		stemcellversion := stemcellversions.StemcellVersion{
			StemcellVersionRef: stemcellversions.StemcellVersionRef{},
			MetalinkSource: map[string]interface{}{
				"source": fmt.Sprintf("%s%s", i.metalinkRepository, strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.localPath)), "/")),
				"include_files": []string{
					"bosh-stemcell-*-warden-boshlite-ubuntu-trusty-go_agent.tgz",
				},
			},
		}

		stemcellversion.StemcellVersionRef.OS = path.Base(path.Dir(path.Dir(meta4Path)))
		stemcellversion.StemcellVersionRef.Version = path.Base(path.Dir(meta4Path))

		inmemory = append(inmemory, stemcellversion)
	}

	return inmemory, nil
}
