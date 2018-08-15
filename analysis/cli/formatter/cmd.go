package formatter

import (
	releaseartifactfilesv1 "github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/cli"
	releasemanifestsv1 "github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/cli"
	stemcellimagefilesv1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellimagefiles.v1/cli"
	stemcellmanifestv1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1/cli"
	stemcellpackagesv1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/cli"
)

type Cmd struct {
	ReleasemanifestsV1     releasemanifestsv1.Cmd     `command:"releasemanifests.v1" subcommands-optional:"true"`
	ReleaseartifactfilesV1 releaseartifactfilesv1.Cmd `command:"releaseartifactfiles.v1" subcommands-optional:"true"`
	StemcellimagefilesV1   stemcellimagefilesv1.Cmd   `command:"stemcellimagefiles.v1" subcommands-optional:"true"`
	StemcellpackagesV1     stemcellpackagesv1.Cmd     `command:"stemcellpackages.v1" subcommands-optional:"true"`
	StemcellmanifestV1     stemcellmanifestv1.Cmd     `command:"stemcellmanifest.v1" subcommands-optional:"true"`
}
