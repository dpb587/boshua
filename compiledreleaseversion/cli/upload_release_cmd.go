package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/pkg/errors"
)

type UploadReleaseCmd struct {
	*CmdOpts `no-flag:"true"`

	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/upload-release")

	artifact, err := c.getCompiledRelease()
	if err != nil {
		return errors.Wrap(err, "finding compiled release")
	}

	var sha1 checksum.Checksum

	for _, cs := range metalinkutil.HashesToChecksums(artifact.ArtifactMetalinkFile().Hashes) {
		if cs.Algorithm().Name() == "sha1" {
			sha1 = &cs

			break
		}
	}

	if c.Cmd {
		fmt.Printf("bosh upload-release --name=%s --version=%s \\\n", c.CompiledReleaseOpts.Release.Name, c.CompiledReleaseOpts.Release.Version)
		fmt.Printf("  %s \\\n", artifact.ArtifactMetalinkFile().URLs[0].URL)
		fmt.Printf("  --sha1=%s \\\n", strings.TrimPrefix(sha1.String(), "sha1:"))
		fmt.Printf("  --stemcell=%s/%s\n", c.CompiledReleaseOpts.OS.Name, c.CompiledReleaseOpts.OS.Version)

		return nil
	}

	cmd := exec.Command(
		"bosh",
		"upload-release",
		fmt.Sprintf("--name=%s", c.CompiledReleaseOpts.Release.Name),
		fmt.Sprintf("--version=%s", c.CompiledReleaseOpts.Release.Version),
		artifact.ArtifactMetalinkFile().URLs[0].URL,
		fmt.Sprintf("--sha1=%s \\\n", strings.TrimPrefix(sha1.String(), "sha1:")),
		fmt.Sprintf("--stemcell=%s/%s\n", c.CompiledReleaseOpts.OS.Name, c.CompiledReleaseOpts.OS.Version),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
