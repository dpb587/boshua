package boshreleasesource

import (
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file"
	"github.com/dpb587/metalink/file/metaurl"
)

type Loader struct{}

var _ metaurl.Loader = &Loader{}

func (f Loader) MediaTypes() []string {
	return []string{
		DefaultMediaType,
	}
}

func (f Loader) Load(source metalink.MetaURL) (file.Reference, error) {
	return NewReference(
		source.URL,
		source.Name,
	), nil
}
