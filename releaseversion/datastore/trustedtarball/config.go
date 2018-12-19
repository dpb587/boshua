package trustedtarball

import "github.com/dpb587/boshua/config/types"

type Config struct {
	// Names defines a list of releases to be trusted. At least one must match.
	Names types.RegexpList `yaml:"names"`

	// URIs defines a list of URI regexes to be trusted. At least one must match.
	URIs types.RegexpList `yaml:"uris"`

	// Labels defines a list of labels to be assigned to releases.
	Labels []string `yaml:"labels"`
}
