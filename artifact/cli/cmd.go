package cli

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
)

type Cmd struct {
	Download       DownloadCmd       `command:"download" description:"For downloading an artifact"`
	UploadRelease  UploadReleaseCmd  `command:"upload-release" description:"For uploading a release to BOSH"`
	UploadStemcell UploadStemcellCmd `command:"upload-stemcell" description:"For uploading a stemcell to BOSH"`
}

func New(app *cmdopts.Opts) *Cmd {
	return &Cmd{}
}
