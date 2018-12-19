package factory

import (
	"fmt"

	"github.com/dpb587/boshua/artifact/downloader/url"
	fileurl "github.com/dpb587/metalink/file/url"
)

type Factory struct {
	handlers map[url.ProviderName]url.Factory
}

func New() *Factory {
	return &Factory{
		handlers: map[url.ProviderName]url.Factory{},
	}
}

func (f *Factory) Add(handler url.ProviderName, factory url.Factory) {
	f.handlers[handler] = factory
}

func (f *Factory) Create(handler url.ProviderName, name string, options map[string]interface{}) (fileurl.Loader, error) {
	factory, found := f.handlers[handler]
	if !found {
		return nil, fmt.Errorf("unsupported handler: %s", handler)
	}

	return factory.Create(handler, name, options)
}
