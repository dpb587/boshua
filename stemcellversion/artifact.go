package stemcellversion

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/dpb587/boshua"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Reference

	MetalinkFile   metalink.File
	MetalinkSource map[string]interface{}
}

var _ boshua.Artifact = &Artifact{}

func (s Artifact) ArtifactReference() boshua.Reference {
	return boshua.Reference{
		Context: "stemcellversion",
		ID:      s.id(),
	}
}

func (s Artifact) ArtifactStorageDir() string {
	ref := s.ArtifactReference()

	return fmt.Sprintf(
		"%s/%s/%s/%s",
		ref.Context,
		ref.ID[0:2],
		ref.ID[2:4],
		ref.ID[4:],
	)
}

func (s Artifact) ArtifactMetalink() metalink.Metalink {
	return metalink.Metalink{
		Files: []metalink.File{
			s.MetalinkFile,
		},
	}
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return s.MetalinkSource
}

func (s Artifact) id() string {
	lightString := "false"
	if s.Reference.Light {
		lightString = "true"
	}

	h := sha1.New()
	h.Write([]byte(strings.Join(
		[]string{
			s.Reference.IaaS,
			s.Reference.Hypervisor,
			s.Reference.OS,
			s.Reference.Version,
			lightString,
		},
		"/",
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
