package output

import (
	"time"

	"github.com/dpb587/boshua/util/checksum"
)

type Result struct {
	Type       string                      `json:"type" yaml:"type"`
	Path       string                      `json:"path" yaml:"path"`
	Link       string                      `json:"link,omitempty" yaml:"link,omitempty"`
	Size       int64                       `json:"size,omitempty" yaml:"size,omitempty"`
	Mode       int64                       `json:"mode" yaml:"mode"`
	Uid        int64                       `json:"uid" yaml:"uid"`
	Gid        int64                       `json:"gid" yaml:"gid"`
	Uname      string                      `json:"uname" yaml:"uname"`
	Gname      string                      `json:"gname" yaml:"gname"`
	ModTime    time.Time                   `json:"modtime" yaml:"modtime"`
	AccessTime *time.Time                  `json:"accesstime,omitempty" yaml:"accesstime,omitempty"`
	ChangeTime *time.Time                  `json:"changetime,omitempty" yaml:"changetime,omitempty"`
	Checksums  checksum.ImmutableChecksums `json:"checksums,omitempty" yaml:"checksums,omitempty"`
}
