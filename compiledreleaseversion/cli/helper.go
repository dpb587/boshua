package cli

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/pkg/errors"
)

func (o *CmdOpts) getCompiledRelease() (compiledreleaseversion.Artifact, error) {
	datastore, err := o.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return compiledreleaseversion.Artifact{}, errors.Wrap(err, "loading compiled release index")
	}

	res, err := datastore.Find(o.CompiledReleaseOpts.Reference())
	if err != nil {
		return compiledreleaseversion.Artifact{}, errors.Wrap(err, "finding compiled release")
	}

	return res, err
}
