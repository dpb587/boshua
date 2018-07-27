package clicommon

import (
	"bytes"
	"log"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/transfer"
	"github.com/dpb587/metalink/verification/hash"
	"github.com/pkg/errors"
)

type DownloadCmd struct {
	Args DownloadCmdArgs `positional-args:"true"`
}

type DownloadCmdArgs struct {
	TargetDir *string `positional-arg-name:"TARGET-DIR" description:"Directory to download files (default: .)"`
}

func (c *DownloadCmd) ExecuteArtifact(loader ArtifactLoader) error {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)
	metaurlLoader := metaurl.NewLoaderFactory()
	metaurlLoader.Add(boshreleasesource.Loader{})

	artifact, err := loader()
	if err != nil {
		log.Fatal(err)
	}

	file := artifact.MetalinkFile()

	localPath := file.Name

	if c.Args.TargetDir != nil {
		localPath = filepath.Join(*c.Args.TargetDir, localPath)
	}

	fullPath, err := filepath.Abs(localPath)
	if err != nil {
		return errors.Wrap(err, "finding output file")
	}

	local, err := urlLoader.Load(metalink.URL{URL: fullPath})
	if err != nil {
		return errors.Wrap(err, "loading download destination")
	}

	progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
	if file.Size == 0 {
		progress.SetWriter(bytes.NewBuffer(nil))
	}

	err = transfer.NewVerifiedTransfer(metaurlLoader, urlLoader, hash.StrongestVerification).TransferFile(file, local, progress)
	if err != nil {
		return errors.Wrap(err, "failed transferring")
	}

	return nil
}
