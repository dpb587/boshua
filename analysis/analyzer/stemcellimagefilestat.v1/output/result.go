package output

import "time"

type Result struct {
	Type       string     `json:"type" yaml:"type"`
	Path       string     `json:"path" yaml:"path"`
	Link       string     `json:"link,omitempty" yaml:"link,omitempty"`
	Size       int64      `json:"size,omitempty" yaml:"size,omitempty"`
	Mode       int64      `json:"mode" yaml:"mode"`
	Uid        int        `json:"uid" yaml:"uid"`
	Gid        int        `json:"gid" yaml:"gid"`
	Uname      string     `json:"uname" yaml:"uname"`
	Gname      string     `json:"gname" yaml:"gname"`
	ModTime    time.Time  `json:"modtime" yaml:"modtime"`
	AccessTime *time.Time `json:"accesstime,omitempty" yaml:"accesstime,omitempty"`
	ChangeTime *time.Time `json:"changetime,omitempty" yaml:"changetime,omitempty"`
}
