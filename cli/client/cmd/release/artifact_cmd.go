package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/util/metalinkutil"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/transfer"
	"github.com/dpb587/metalink/verification/hash"
)

type ArtifactCmd struct {
	*CmdOpts `no-flag:"true"`

	Download *string `long:"download" description:"Download the release artifact" value-name:"OPTIONAL-PATH" optional:"true" optional-value:"default"`
	Format   string  `long:"format" description:"Output format for the release reference" value-name:"json|metalink|tsv" default:"tsv"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/artifact")

	client := c.AppOpts.GetClient()

	res, err := client.GetReleaseVersion(releaseversion.Reference{
		Name:      c.ReleaseOpts.Release.Name,
		Version:   c.ReleaseOpts.Release.Version,
		Checksums: checksum.ImmutableChecksums{c.ReleaseOpts.ReleaseChecksum.ImmutableChecksum},
	})
	if err != nil {
		return fmt.Errorf("fetching: %v", err)
	}

	if c.Download != nil {
		logger := boshlog.NewLogger(boshlog.LevelError)
		fs := boshsys.NewOsFileSystem(logger)

		urlLoader := urldefaultloader.New(fs)
		metaurlLoader := metaurl.NewLoaderFactory()

		file := res.Data.Artifact
		target := *c.Download

		if target == "default" {
			target = file.Name
		}

		local, err := urlLoader.Load(metalink.URL{URL: target})
		if err != nil {
			return fmt.Errorf("loading download destination: %v", err)
		}

		progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)

		return transfer.NewVerifiedTransfer(metaurlLoader, urlLoader, hash.StrongestVerification).TransferFile(file, local, progress)
	}

	if c.Format == "json" {
		output := map[string]string{
			"file": res.Data.Artifact.Name,
			"url":  res.Data.Artifact.URLs[0].URL,
		}
		for _, cs := range metalinkutil.HashesToChecksums(res.Data.Artifact.Hashes) {
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
				res.Data.Artifact,
			},
			Generator: "bosh-compiled-releases/0.0.0",
		}

		meta4Bytes, err := metalink.Marshal(meta4)
		if err != nil {
			log.Fatalf("marshalling response: %v", err)
		}

		fmt.Printf("%s\n", meta4Bytes)
	} else if c.Format == "tsv" {
		fmt.Printf("file\t%s\n", res.Data.Artifact.Name)
		fmt.Printf("url\t%s\n", res.Data.Artifact.URLs[0].URL)
		for _, cs := range metalinkutil.HashesToChecksums(res.Data.Artifact.Hashes) {
			fmt.Printf("%s\n", strings.Replace(cs.String(), ":", "\t", 1))
		}
	} else {
		return errors.New("invalid format")
	}

	return nil
}
