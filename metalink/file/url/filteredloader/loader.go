package filteredloader

import (
	"regexp"

	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file"
	"github.com/dpb587/metalink/file/url"
)

type filteredLoader struct {
	loader url.Loader
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

func NewLoader(loader url.Loader, include []*regexp.Regexp, exclude []*regexp.Regexp) url.Loader {
	return &filteredLoader{
		loader: loader,
		include: include,
		exclude: exclude,
	}
}

func (l *filteredLoader) Load(source metalink.URL) (file.Reference, error) {
	for _, exclude := range l.exclude {
		if exclude.MatchString(source.URL) {
			return nil, url.UnsupportedURLError
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
		return nil, url.UnsupportedURLError
	}

	return l.loader.Load(source)
}

func (l *filteredLoader) Schemes() []string {
	return l.loader.Schemes()
}
