package cli

import (
	"github.com/dpb587/boshua/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/pkg/errors"
)

func (o *CmdOpts) getStemcell() (stemcellversion.Artifact, error) {
	index, err := o.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "loading stemcell index")
	}

	res, err := index.Filter(o.StemcellOpts.FilterParams())
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	} else if err = datastore.RequireSingleResult(res); err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	}

	return res[0], err
}
