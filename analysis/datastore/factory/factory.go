package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis/datastore"
)

type Factory struct {
	providers map[datastore.ProviderName]datastore.Factory
}

func New() *Factory {
	return &Factory{
		providers: map[datastore.ProviderName]datastore.Factory{},
	}
}

func (f *Factory) Add(provider datastore.ProviderName, factory datastore.Factory) {
	f.providers[provider] = factory
}

func (f *Factory) Create(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	factory, found := f.providers[provider]
	if !found {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return factory.Create(provider, name, options)
}
