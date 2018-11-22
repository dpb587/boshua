package clicommon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum/algorithm"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/url/file"
	"github.com/dpb587/metalink/verification"
	"github.com/pkg/errors"
)

type ArtifactCmd struct {
	Download *string `long:"download" description:"Download the release artifact" value-name:"OPTIONAL-PATH" optional:"true" optional-value:"default"`
	Format   string  `long:"format" description:"Output format for the release reference" value-name:"json|metalink|tsv" default:"tsv"`
}

func (c *ArtifactCmd) ExecuteArtifact(downloaderGetter DownloaderGetter, loader artifact.Loader) error {
	artifact, err := loader()
	if err != nil {
		log.Fatal(err)
	}

	artifactMetalinkFile := artifact.MetalinkFile()

	if c.Download != nil {
		downloader, err := downloaderGetter()
		if err != nil {
			log.Fatal(err)
		}

		target := *c.Download

		if target == "default" {
			target = artifactMetalinkFile.Name
		}

		targetPath, err := filepath.Abs(target)
		if err != nil {
			return errors.Wrap(err, "finding download path")
		}

		progress := pb.New64(int64(artifactMetalinkFile.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)

		return downloader.TransferFile(
			artifactMetalinkFile,
			file.NewReference(targetPath),
			progress,
			verification.NewSimpleVerificationResultReporter(os.Stdout),
		)
	}

	if c.Format == "json" {
		output := map[string]string{
			"file": artifactMetalinkFile.Name,
		}

		for _, url := range artifactMetalinkFile.URLs { // TODO only first?
			output["url"] = url.URL

			break
		}

		for _, metaurl := range artifactMetalinkFile.MetaURLs { // TODO only first?
			output["metaurl"] = metaurl.URL

			break
		}

		for _, cs := range metalinkutil.HashesToChecksums(artifactMetalinkFile.Hashes) {
			output[cs.Algorithm().Name()] = fmt.Sprintf("%x", cs.Data())
		}

		outputBytes, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatalf("marshalling response: %v", err)
		}

		fmt.Printf("%s\n", outputBytes)
	} else if c.Format == "metalink" {
		meta4 := metalink.Metalink{
			Files: []metalink.File{
				artifactMetalinkFile,
			},
			Generator: "boshua/0.0.0",
		}

		meta4Bytes, err := metalink.MarshalXML(meta4)
		if err != nil {
			log.Fatalf("marshalling response: %v", err)
		}

		fmt.Printf("%s\n", meta4Bytes)
	} else if c.Format == "tsv" {
		fmt.Printf("file\t%s\n", artifactMetalinkFile.Name)

		for _, url := range artifactMetalinkFile.URLs {
			fmt.Printf("url\t%s\n", url.URL)
		}

		for _, cs := range metalinkutil.HashesToChecksums(artifactMetalinkFile.Hashes) {
			switch cs.Algorithm().Name() {
			case algorithm.SHA1, algorithm.SHA256:
				fmt.Printf("%s\n", strings.Replace(cs.String(), ":", "\t", 1))
			}
		}

		for _, metaurl := range artifactMetalinkFile.MetaURLs {
			fmt.Printf("metaurl\t%s\t%s\n", metaurl.URL, metaurl.MediaType)
		}

		for _, label := range artifact.GetLabels() {
			fmt.Printf("label\t%s\n", label)
		}
	} else {
		return errors.New("invalid format")
	}

	return nil
}
