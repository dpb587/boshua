package clicommon

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/transfer"
	"github.com/dpb587/metalink/verification/hash"
	"github.com/pkg/errors"
)

type ArtifactCmd struct {
	Download *string `long:"download" description:"Download the release artifact" value-name:"OPTIONAL-PATH" optional:"true" optional-value:"default"`
	Format   string  `long:"format" description:"Output format for the release reference" value-name:"json|metalink|tsv" default:"tsv"`
}

func (c *ArtifactCmd) ExecuteArtifact(loader artifact.Loader) error {
	artifact, err := loader()
	if err != nil {
		log.Fatal(err)
	}

	artifactMetalinkFile := artifact.MetalinkFile()

	if c.Download != nil {
		logger := boshlog.NewLogger(boshlog.LevelError)
		fs := boshsys.NewOsFileSystem(logger)

		urlLoader := urldefaultloader.New(fs)
		metaurlLoader := metaurl.NewLoaderFactory()
		metaurlLoader.Add(boshreleasesource.Loader{})

		target := *c.Download

		if target == "default" {
			target = artifactMetalinkFile.Name
		}

		targetPath, err := filepath.Abs(target)
		if err != nil {
			return errors.Wrap(err, "finding download path")
		}

		local, err := urlLoader.Load(metalink.URL{URL: targetPath})
		if err != nil {
			return errors.Wrap(err, "loading download destination")
		}

		progress := pb.New64(int64(artifactMetalinkFile.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)

		return transfer.NewVerifiedTransfer(metaurlLoader, urlLoader, hash.StrongestVerification).TransferFile(artifactMetalinkFile, local, progress)
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

		for _, url := range artifactMetalinkFile.MetaURLs {
			fmt.Printf("metaurl\t%s\t%s\n", url.URL, url.MediaType)
		}

		for _, cs := range metalinkutil.HashesToChecksums(artifactMetalinkFile.Hashes) {
			fmt.Printf("%s\n", strings.Replace(cs.String(), ":", "\t", 1))
		}
	} else {
		return errors.New("invalid format")
	}

	return nil
}
