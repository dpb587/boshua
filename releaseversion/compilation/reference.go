package compilation

import (
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

type Reference struct {
	ReleaseVersion releaseversion.Reference `json:"release"`
	OSVersion      osversion.Reference      `json:"os"`
}
