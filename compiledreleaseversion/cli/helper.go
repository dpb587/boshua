package cli

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/pkg/errors"
)

func (o *CmdOpts) getCompiledRelease() (compiledreleaseversion.Artifact, error) {
	index, err := o.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return compiledreleaseversion.Artifact{}, errors.Wrap(err, "loading compiled release index")
	}

	results, err := index.Filter(o.CompiledReleaseOpts.FilterParams())
	if err != nil {
		return compiledreleaseversion.Artifact{}, errors.Wrap(err, "finding compiled release")
	}

	result, err := datastore.RequireSingleResult(results)
	if err != nil {
		return compiledreleaseversion.Artifact{}, errors.Wrap(err, "finding compiled release")
	}

	return result, err
}
