package compiledrelease

import (
	"fmt"
	"os"

	"github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

func (o *CmdOpts) getCompiledRelease() (*compiledreleaseversion.GETCompilationResponse, error) {
	client := o.AppOpts.GetClient()

	releaseVersionRef := releaseversion.Reference{
		Name:      o.CompiledReleaseOpts.Release.Name,
		Version:   o.CompiledReleaseOpts.Release.Version,
		Checksums: checksum.ImmutableChecksums{o.CompiledReleaseOpts.ReleaseChecksum.ImmutableChecksum},
	}
	osVersionRef := osversion.Reference{
		Name:    o.CompiledReleaseOpts.OS.Name,
		Version: o.CompiledReleaseOpts.OS.Version,
	}

	if o.CompiledReleaseOpts.NoWait {
		return client.GetCompiledReleaseVersionCompilation(releaseVersionRef, osVersionRef)
	}

	return client.RequireCompiledReleaseVersionCompilation(
		releaseVersionRef,
		osVersionRef,
		func(task scheduler.TaskStatus) {
			if o.AppOpts.Quiet {
				return
			}

			fmt.Fprintf(os.Stderr, "compilation status: %s\n", task.Status)
		},
	)
}
