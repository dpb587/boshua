package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/pkg/errors"
)

type UploadReleaseCmd struct {
	*CmdOpts `no-flag:"true"`

	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/upload-release")

	artifact, err := c.CompiledReleaseOpts.Artifact()
	if err != nil {
		return errors.Wrap(err, "finding compiled release")
	}

	artifactRef := artifact.Reference().(compilation.Reference)

	var sha1 checksum.Checksum

	for _, cs := range metalinkutil.HashesToChecksums(artifact.MetalinkFile().Hashes) {
		if cs.Algorithm().Name() == "sha1" {
			sha1 = &cs

			break
		}
	}

	if c.Cmd {
		fmt.Printf("bosh upload-release --name=%s --version=%s \\\n", artifactRef.ReleaseVersion.Name, artifactRef.ReleaseVersion.Version)
		fmt.Printf("  %s \\\n", artifact.MetalinkFile().URLs[0].URL)
		fmt.Printf("  --sha1=%s \\\n", strings.TrimPrefix(sha1.String(), "sha1:"))
		fmt.Printf("  --stemcell=%s/%s\n", artifactRef.OSVersion.Name, artifactRef.OSVersion.Version)

		return nil
	}

	cmd := exec.Command(
		"bosh",
		"upload-release",
		fmt.Sprintf("--name=%s", artifactRef.ReleaseVersion.Name),
		fmt.Sprintf("--version=%s", artifactRef.ReleaseVersion.Version),
		artifact.MetalinkFile().URLs[0].URL,
		fmt.Sprintf("--sha1=%s \\\n", strings.TrimPrefix(sha1.String(), "sha1:")),
		fmt.Sprintf("--stemcell=%s/%s\n", artifactRef.OSVersion.Name, artifactRef.OSVersion.Version),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
