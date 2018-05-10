package clicommon

import (
	"compress/gzip"
	"io"
	"os/exec"
	// "io"
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
	"github.com/pkg/errors"
)

type ResultsCmd struct {
	Raw bool `long:"raw" description:"Show raw, unformatted analysis results"`
}

func (c *ResultsCmd) ExecuteAnalysis(analyzer string, loader AnalysisLoader, args []string) error {
	artifact, err := loader()
	if err != nil {
		return errors.Wrap(err, "finding analysis")
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

	file := artifact.ArtifactMetalinkFile()

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

	if c.Raw {
		_, err = io.Copy(os.Stdout, gzReader)
		if err != nil {
			log.Fatalf("piping results: %v", err)
		}

		return nil
	}

	formatterArgs := append([]string{"analysis", "formatter", analyzer}, args...)
	formatterCmd := exec.Command(os.Args[0], formatterArgs...)
	formatterCmd.Stdin = gzReader
	formatterCmd.Stdout = os.Stdout
	formatterCmd.Stderr = os.Stderr

	return formatterCmd.Run()
}
