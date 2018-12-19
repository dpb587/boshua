package clicommon

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/boshua/analysis"
	fileurl "github.com/dpb587/metalink/file/url/file"
	"github.com/dpb587/metalink/verification"
	"github.com/pkg/errors"
)

type ResultsCmd struct {
	Raw bool `long:"raw" description:"Show raw, unformatted analysis results"`
}

func (c *ResultsCmd) ExecuteAnalysis(downloaderGetter DownloaderGetter, analyzer analysis.AnalyzerName, loader AnalysisLoader, args []string) error {
	downloader, err := downloaderGetter()
	if err != nil {
		return errors.Wrap(err, "loading downloader")
	}

	artifact, err := loader()
	if err != nil {
		return errors.Wrap(err, "loading analysis")
	}

	tempfile, err := ioutil.TempFile("", "boshua-analysis-")
	if err != nil {
		log.Fatalf("creating temp file for download: %v", err)
	}

	defer os.Remove(tempfile.Name())

	file := artifact.MetalinkFile()

	progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
	progress.SetWriter(ioutil.Discard)

	err = downloader.TransferFile(
		file,
		fileurl.NewReference(tempfile.Name()),
		progress,
		verification.NewSimpleVerificationResultReporter(ioutil.Discard),
	)
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

	formatterArgs := append([]string{"analysis", "formatter", string(analyzer)}, args...)
	formatterCmd := exec.Command(os.Args[0], formatterArgs...)
	formatterCmd.Stdin = gzReader
	formatterCmd.Stdout = os.Stdout
	formatterCmd.Stderr = os.Stderr

	return formatterCmd.Run()
}
