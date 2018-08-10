package loader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/configdef"
	"github.com/dpb587/boshua/config/provider"
	"github.com/pkg/errors"
)

const DefaultPath = "~/.config/boshua/config.yml"

func LoadFromFile(path string, cfg *config.Config) (*provider.Config, error) {
	configPath, isDefault := absPath(path)

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) || !isDefault {
			return nil, errors.Wrap(err, "reading config")
		}

		configBytes = []byte("--- {}")
	}

	err = configdef.UnmarshalYAML(configBytes, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading options")
	}

	return &provider.Config{
		Config:    cfg,
		RawConfig: configBytes,
	}, nil
}

func absPath(path string) (string, bool) {
	var isDefault = path == DefaultPath

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(os.Getenv("HOME"), path[1:])
	}

	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return path, isDefault
}
