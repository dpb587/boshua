package cli

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
)

func (o *CmdOpts) getStemcell() (stemcellversion.Artifact, error) {
	datastore, err := o.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return stemcellversion.Artifact{}, fmt.Errorf("loading stemcell index: %v", err)
	}

	res, err := datastore.Find(o.StemcellOpts.Reference())
	if err != nil {
		return stemcellversion.Artifact{}, fmt.Errorf("finding stemcell: %v", err)
	}

	return res, err
}
