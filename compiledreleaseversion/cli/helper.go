package cli

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion"
)

func (o *CmdOpts) getCompiledRelease() (compiledreleaseversion.Artifact, error) {
	datastore, err := o.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return compiledreleaseversion.Artifact{}, fmt.Errorf("loading compiled release index: %v", err)
	}

	res, err := datastore.Find(o.CompiledReleaseOpts.Reference())
	if err != nil {
		return compiledreleaseversion.Artifact{}, fmt.Errorf("finding compiled release: %v", err)
	}

	return res, err
}
