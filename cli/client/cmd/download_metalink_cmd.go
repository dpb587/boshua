package cmd

import (
	"fmt"
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
		return fmt.Errorf("reading metalink: %v", err)
	}

	var meta4 metalink.Metalink

	err = metalink.Unmarshal(meta4Bytes, &meta4)
	if err != nil {
		return fmt.Errorf("unmarshaling metalink: %v", err)
	}

	for _, file := range meta4.Files {
		localPath := file.Name

		if c.Args.TargetDir != nil {
			localPath = filepath.Join(*c.Args.TargetDir, localPath)
		}

		local, err := urlLoader.Load(metalink.URL{URL: localPath})
		if err != nil {
			return fmt.Errorf("loading download destination: %v", err)
		}

		progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)

		err = transfer.NewVerifiedTransfer(metaurlLoader, urlLoader, hash.StrongestVerification).TransferFile(file, local, progress)
		if err != nil {
			return fmt.Errorf("failed transferring: %v", err)
		}
	}

	return nil
}
