package filteredloader

import (
	"regexp"

	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file"
	"github.com/dpb587/metalink/file/metaurl"
)

type filteredLoader struct {
	loader metaurl.Loader
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

func NewLoader(loader metaurl.Loader, include []*regexp.Regexp, exclude []*regexp.Regexp) metaurl.Loader {
	return &filteredLoader{
		loader: loader,
		include: include,
		exclude: exclude,
	}
}

func (l *filteredLoader) Load(source metalink.MetaURL) (file.Reference, error) {
	for _, exclude := range l.exclude {
		if exclude.MatchString(source.URL) {
			return nil, metaurl.UnsupportedMetaURLError
		}
	}

	var included bool

	for _, include := range l.include {
		if include.MatchString(source.URL) {
			included = true

			break
		}
	}

	if !included {
		return nil, metaurl.UnsupportedMetaURLError
	}

	return l.loader.Load(source)
}

func (l *filteredLoader) MediaTypes() []string {
	return l.loader.MediaTypes()
}
