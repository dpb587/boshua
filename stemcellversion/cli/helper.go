package cli

import (
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/pkg/errors"
)

func (o *CmdOpts) getStemcell() (stemcellversion.Artifact, error) {
	datastore, err := o.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "loading stemcell index")
	}

	res, err := datastore.Find(o.StemcellOpts.Reference())
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	}

	return res, err
}
