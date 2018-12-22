package cli

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
)

type Cmd struct {
	UseCompiledReleasesCmd    UseCompiledReleasesCmd    `command:"use-compiled-releases" description:"For patching a manifest to refer to compiled releases"`
	UploadCompiledReleasesCmd UploadCompiledReleasesCmd `command:"upload-compiled-releases" description:"For uploading compiled releases referenced by a manifest"`
	UploadStemcellsCmd        UploadStemcellsCmd        `command:"upload-stemcells" description:"For uploading stemcells referenced by a manifest"`
}

func New(app *cmdopts.Opts) *Cmd {
	return &Cmd{}
}
