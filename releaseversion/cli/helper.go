package cli

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion"
)

func (o *CmdOpts) getRelease() (releaseversion.Artifact, error) {
	datastore, err := o.AppOpts.GetReleaseIndex("default")
	if err != nil {
		return releaseversion.Artifact{}, fmt.Errorf("loading release index: %v", err)
	}

	res, err := datastore.Find(o.ReleaseOpts.Reference())
	if err != nil {
		return releaseversion.Artifact{}, fmt.Errorf("finding release: %v", err)
	}

	return res, err
}
