package cmd

import (
	"bytes"
	"io/ioutil"
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

type DownloadMetalinkCmd struct {
	*CmdOpts `no-flag:"true"`

	Args DownloadMetalinkCmdArgs `positional-args:"true" required:"true"`
}

type DownloadMetalinkCmdArgs struct {
	Metalink  string  `positional-arg-name:"PATH" description:"Path to the metalink file"`
	TargetDir *string `positional-arg-name:"TARGET-DIR" description:"Directory to download files"`
}

func (c *DownloadMetalinkCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("download-metalink")

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)
	metaurlLoader := metaurl.NewLoaderFactory()
	metaurlLoader.Add(boshreleasesource.Loader{})

	meta4Bytes, err := ioutil.ReadFile(c.Args.Metalink)
	if err != nil {
		return errors.Wrap(err, "reading metalink")
	}

	var meta4 metalink.Metalink

	err = metalink.Unmarshal(meta4Bytes, &meta4)
	if err != nil {
		return errors.Wrap(err, "unmarshaling metalink")
	}

	for _, file := range meta4.Files {
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
	}

	return nil
}
