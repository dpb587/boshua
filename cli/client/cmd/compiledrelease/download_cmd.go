package compiledrelease

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/transfer"
	"github.com/dpb587/metalink/verification/hash"
)

type DownloadCmd struct {
	*CmdOpts `no-flag:"true"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`

	Quiet bool `long:"quiet" description:"Suppress informational output"`
}

func (c *DownloadCmd) Execute(_ []string) error {
	apiclient := client.New(http.DefaultClient, c.AppOpts.Server)

	releaseRef := models.ReleaseRef{
		Name:     c.CompiledReleaseOpts.Release.Name,
		Version:  c.CompiledReleaseOpts.Release.Version,
		Checksum: models.Checksum(c.CompiledReleaseOpts.ReleaseChecksum.String()),
	}
	stemcellRef := models.StemcellRef{
		OS:      c.CompiledReleaseOpts.Stemcell.OS,
		Version: c.CompiledReleaseOpts.Stemcell.Version,
	}

	var resInfo *models.CRVInfoResponse
	var err error

	if c.RequestAndWait {
		resInfo, err = client.RequestAndWait(apiclient, releaseRef, stemcellRef)
	} else {
		resInfo, err = apiclient.CompiledReleaseVersionInfo(models.CRVInfoRequest{
			Data: models.CRVInfoRequestData{
				Release:  releaseRef,
				Stemcell: stemcellRef,
			},
		})
	}

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
