package boshuaV2

import (
	boshuaV2 "github.com/dpb587/boshua/datastore/boshua.v2"
)

type Config struct {
	boshuaV2.BoshuaConfig `yaml:"-,inline"`
}
