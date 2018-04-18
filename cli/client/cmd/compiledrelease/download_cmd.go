package compiledrelease

import (
	"fmt"
	"log"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/transfer"
	"github.com/dpb587/metalink/verification/hash"
)

type DownloadCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *DownloadCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiled-release/download")

	resInfo, err := c.CompiledReleaseOpts.GetCompiledReleaseVersion(c.AppOpts.GetClient())
	if err != nil {
		log.Fatalf("requesting compiled version info: %v", err)
	} else if resInfo == nil {
		log.Fatalf("no compiled release available")
	}

	meta4 := createMetalink(resInfo)

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)
	metaurlLoader := metaurl.NewLoaderFactory()

	file := meta4.Files[0]

	local, err := urlLoader.Load(metalink.URL{URL: file.Name})
	if err != nil {
		return fmt.Errorf("loading download destination: %v", err)
	}

	progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)

	return transfer.NewVerifiedTransfer(metaurlLoader, urlLoader, hash.StrongestVerification).TransferFile(file, local, progress)
}
