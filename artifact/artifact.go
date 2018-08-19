package artifact

import (
	"github.com/dpb587/metalink"
)

type Artifact interface {
	Reference() interface{}
	MetalinkFile() metalink.File
	GetLabels() []string
	GetDatastoreName() string
}
