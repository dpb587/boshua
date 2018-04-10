package releaseartifactfilestat

import "time"

type Record struct {
	Artifact string         `json:"artifact"`
	Path     string         `json:"path"`
	Result   RecordFileStat `json:"result"`
}

type RecordFileStat struct {
	Type       string     `json:"type"`
	Path       string     `json:"path"`
	Link       string     `json:"link,omitempty"`
	Size       int64      `json:"size,omitempty"`
	Mode       int64      `json:"mode"`
	Uid        int        `json:"uid"`
	Gid        int        `json:"gid"`
	Uname      string     `json:"uname"`
	Gname      string     `json:"gname"`
	ModTime    time.Time  `json:"modtime"`
	AccessTime *time.Time `json:"accesstime,omitempty"`
	ChangeTime *time.Time `json:"changetime,omitempty"`
}
