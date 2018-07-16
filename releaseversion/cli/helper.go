package cli

import (
	"github.com/dpb587/boshua/datastore"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/pkg/errors"
)

func (o *CmdOpts) getRelease() (releaseversion.Artifact, error) {
	index, err := o.AppOpts.GetReleaseIndex("default")
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "loading release index")
	}

	res, err := index.Filter(o.ReleaseOpts.FilterParams())
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "finding release")
	} else if err = datastore.RequireSingleResult(res); err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "finding release")
	}

	return res[0], err
}
