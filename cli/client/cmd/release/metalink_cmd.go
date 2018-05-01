package release

import (
	"fmt"
	"log"

	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
)

type MetalinkCmd struct {
	*CmdOpts `no-flag:"true"`

	Format string `long:"format" description:"Output format for metalink"`
}

func (c *MetalinkCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/metalink")

	client := c.AppOpts.GetClient()

	res, err := client.GetReleaseVersion(releaseversion.Reference{
		Name:      c.ReleaseOpts.Release.Name,
		Version:   c.ReleaseOpts.Release.Version,
		Checksums: checksum.ImmutableChecksums{c.ReleaseOpts.ReleaseChecksum.ImmutableChecksum},
	})
	if err != nil {
		return fmt.Errorf("fetching: %v", err)
	}

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

	return nil
}
