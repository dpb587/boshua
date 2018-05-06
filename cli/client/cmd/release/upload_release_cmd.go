package release

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/releaseversion"
)

type UploadReleaseCmd struct {
	*CmdOpts `no-flag:"true"`

	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/upload-release")

	client := c.AppOpts.GetClient()

	res, err := client.GetReleaseVersion(releaseversion.Reference{
		Name:      c.ReleaseOpts.Release.Name,
		Version:   c.ReleaseOpts.Release.Version,
		Checksums: checksum.ImmutableChecksums{c.ReleaseOpts.ReleaseChecksum.ImmutableChecksum},
	})
	if err != nil {
		return fmt.Errorf("fetching: %v", err)
	}

	if c.Cmd {
		fmt.Printf("bosh upload-release --name=%s --version=%s \\\n", c.ReleaseOpts.Release.Name, c.ReleaseOpts.Release.Version)
		fmt.Printf("  %s \\\n", res.Data.Artifact.URLs[0].URL)
		fmt.Printf("  --sha1=%s\n", strings.TrimPrefix(c.ReleaseOpts.ReleaseChecksum.ImmutableChecksum.String(), "sha1:"))

		return nil
	}

	cmd := exec.Command(
		"bosh",
		"upload-release",
		fmt.Sprintf("--name=%s", c.ReleaseOpts.Release.Name),
		fmt.Sprintf("--version=%s", c.ReleaseOpts.Release.Version),
		res.Data.Artifact.URLs[0].URL,
		fmt.Sprintf("--sha1=%s", strings.TrimPrefix(c.ReleaseOpts.ReleaseChecksum.ImmutableChecksum.String(), "sha1:")),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
