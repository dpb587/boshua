package opts

import (
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/provider"
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	"github.com/dpb587/boshua/pivnetfile"
	"github.com/dpb587/boshua/pivnetfile/datastore"
	"github.com/pkg/errors"
)

type Opts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`

	PivnetProduct   string `long:"pivnet-product" description:"The product name/slug"`
	PivnetReleaseID int    `long:"pivnet-release-id" description:"The release ID"`
	PivnetFileID    int    `long:"pivnet-file-id" description:"The file ID"`
}

func (o *Opts) Artifact(cfg *provider.Config) (pivnetfile.Artifact, error) {
	index, err := cfg.GetPivnetFileIndex(config.DefaultName)
	if err != nil {
		return pivnetfile.Artifact{}, errors.Wrap(err, "loading pivnet file index")
	}

	result, err := datastore.GetArtifact(index, o.FilterParams())
	if err != nil {
		return pivnetfile.Artifact{}, errors.Wrap(err, "finding pivnet files")
	}

	return result, nil
}

func (o Opts) FilterParams() datastore.FilterParams {
	f := datastore.FilterParams{}

	f.ProductNameExpected = o.PivnetProduct != ""
	f.ProductName = o.PivnetProduct

	f.ReleaseIDExpected = o.PivnetReleaseID != 0
	f.ReleaseID = o.PivnetReleaseID

	f.FileIDExpected = o.PivnetFileID != 0
	f.FileID = o.PivnetFileID

	return f
}
