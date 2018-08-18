package storage

type StorageConfig []StorageConfigEntry

type StorageConfigEntry struct {
	URI      string            `yaml:"uri"`
	Location string            `yaml:"location"`
	Priority *uint             `yaml:"priority"`
	Options  map[string]string `yaml:"options"`
}
