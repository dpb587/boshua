package compiledrelease

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
)

func (o *CmdOpts) getCompiledRelease() (*compiledreleaseversion.GETCompilationInfoResponse, error) {
	client := o.AppOpts.GetClient()

	ref := o.CompiledReleaseOpts.Reference()

	if o.CompiledReleaseOpts.NoWait {
		return client.GetCompiledReleaseVersionCompilation(ref.ReleaseVersion, ref.OSVersion)
	}

	return client.RequireCompiledReleaseVersionCompilation(
		ref.ReleaseVersion,
		ref.OSVersion,
		func(task scheduler.TaskStatus) {
			if !o.AppOpts.Quiet {
				fmt.Fprintf(
					os.Stderr,
					"boshua | %s | fetching compiled release: %s/%s: %s/%s: compilation %s\n",
					time.Now().Format("15:04:05"),
					ref.OSVersion.Name,
					ref.OSVersion.Version,
					ref.ReleaseVersion.Name,
					ref.ReleaseVersion.Version,
					task.Status,
				)
			}
		},
	)
}
