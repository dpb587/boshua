package trustedtarball

import "github.com/dpb587/boshua/config/types"

type Config struct {
	// Names defines a list of releases to be trusted. At least one must match.
	Names []types.Regexp `yaml:"names"`

	// URIs defines a list of URI regexes to be trusted. At least one must match.
	URIs []types.Regexp `yaml:"uris"`
}
