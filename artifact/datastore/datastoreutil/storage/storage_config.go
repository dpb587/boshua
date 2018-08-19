package storage

// StorageConfig defines a list of storage locations to mirror artifacts.
type StorageConfig []StorageConfigEntry

// StorageConfigEntry defines a specific mirror location.
type StorageConfigEntry struct {
	// URI must define a metalink-friendly file path which supports writing.
	URI string `yaml:"uri"`

	// Location may define an ISO3166-1 alpha-2 country code for the geographical
	// location.
	Location string `yaml:"location"`

	// Priority may specify a priority for this mirror.
	Priority *uint `yaml:"priority"`

	// Options should specify credentials which are required for write-operations
	// for the URI-specific schema (e.g. AWS credentials for S3).
	Options map[string]string `yaml:"options"`
}
