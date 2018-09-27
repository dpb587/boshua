package url

import (
	"github.com/dpb587/metalink/file/url"
)

type ProviderName string

type Factory interface {
	Create(provider ProviderName, name string, options map[string]interface{}) (url.Loader, error)
}
