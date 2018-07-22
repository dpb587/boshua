package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

type Cmd struct {
	Download       DownloadCmd       `command:"download" description:"For downloading an artifact"`
	UploadRelease  UploadReleaseCmd  `command:"upload-release" description:"For uploading a release to BOSH"`
	UploadStemcell UploadStemcellCmd `command:"upload-stemcell" description:"For uploading a stemcell to BOSH"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.Download.CmdOpts = cmdOpts
	cmd.UploadRelease.CmdOpts = cmdOpts
	cmd.UploadStemcell.CmdOpts = cmdOpts

	return cmd
}
