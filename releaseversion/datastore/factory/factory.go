package factory

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/datastore"
)

type Factory struct {
	providers map[string]datastore.Factory
}

func New() *Factory {
	return &Factory{
		providers: map[string]datastore.Factory{},
	}
}

func (f *Factory) Add(provider string, factory datastore.Factory) {
	f.providers[provider] = factory
}

func (f *Factory) Create(provider, name string, options map[string]interface{}) (datastore.Index, error) {
	factory, found := f.providers[provider]
	if !found {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return factory.Create(provider, name, options)
}
