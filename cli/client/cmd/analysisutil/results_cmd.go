package analysisutil

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
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

type ResultsCmd struct{}

func (c *ResultsCmd) ExecuteAnalysis(loader AnalysisLoader) error {
	resInfo, err := loader()
	if err != nil {
		log.Fatal(err)
	} else if resInfo == nil {
		log.Fatalf("no analysis available")
	}

	tempfile, err := ioutil.TempFile("", "boshua-analysis-")
	if err != nil {
		log.Fatalf("creating temp file for download: %v", err)
	}

	defer os.Remove(tempfile.Name())

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)
	metaurlLoader := metaurl.NewLoaderFactory()

	file := resInfo.Data.Artifact

	local, err := urlLoader.Load(metalink.URL{URL: tempfile.Name()})
	if err != nil {
		log.Fatalf("loading download destination: %v", err)
	}

	progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
	progress.SetWriter(ioutil.Discard)

	err = transfer.NewVerifiedTransfer(metaurlLoader, urlLoader, hash.StrongestVerification).TransferFile(file, local, progress)
	if err != nil {
		log.Fatalf("downloading results: %v", err)
	}

	gzReader, err := gzip.NewReader(tempfile)
	if err != nil {
		log.Fatalf("starting gzip: %v", err)
	}

	_, err = io.Copy(os.Stdout, gzReader)
	if err != nil {
		log.Fatalf("piping results: %v", err)
	}

	return nil
}
