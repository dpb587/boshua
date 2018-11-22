package clicommon

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/boshua/artifact"
  "github.com/dpb587/metalink/file/url/file"
	"github.com/dpb587/metalink/verification"
	"github.com/pkg/errors"
)

type DownloadCmd struct {
	Rename string          `long:"rename" description:"Override the downloaded file name"`
	Args   DownloadCmdArgs `positional-args:"true"`
}

type DownloadCmdArgs struct {
	TargetDir *string `positional-arg-name:"TARGET-DIR" description:"Directory to download files (default: .)"`
}

func (c *DownloadCmd) ExecuteArtifact(downloaderGetter DownloaderGetter, loader artifact.Loader) error {
	downloader, err := downloaderGetter()
	if err != nil {
		log.Fatal(err)
	}

	artifact, err := loader()
	if err != nil {
		log.Fatal(err)
	}

	artifactMetalinkFile := artifact.MetalinkFile()

	localPath := artifactMetalinkFile.Name

	if c.Rename != "" {
		localPath = c.Rename
	}

	if c.Args.TargetDir != nil {
		localPath = filepath.Join(*c.Args.TargetDir, localPath)
	}

	fullPath, err := filepath.Abs(localPath)
	if err != nil {
		return errors.Wrap(err, "finding output file")
	}

	progress := pb.New64(int64(artifactMetalinkFile.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
	if artifactMetalinkFile.Size == 0 {
		progress.SetWriter(bytes.NewBuffer(nil))
	}

	err = downloader.TransferFile(
		artifactMetalinkFile,
		file.NewReference(fullPath),
		progress,
		verification.NewSimpleVerificationResultReporter(os.Stdout),
	)
	if err != nil {
		return errors.Wrap(err, "failed transferring")
	}

	return nil
}
