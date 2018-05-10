package cli

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/pkg/errors"
)

func (o *CmdOpts) getRelease() (releaseversion.Artifact, error) {
	datastore, err := o.AppOpts.GetReleaseIndex("default")
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "loading release index")
	}

	res, err := datastore.Find(o.ReleaseOpts.Reference())
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "finding release")
	}

	return res, err
}
