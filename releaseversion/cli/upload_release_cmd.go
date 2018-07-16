package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type UploadReleaseCmd struct {
	*CmdOpts `no-flag:"true"`

	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/upload-release")

	artifact, err := c.getRelease()
	if err != nil {
		return errors.Wrap(err, "finding compiled release")
	}

	if c.Cmd {
		fmt.Printf("bosh upload-release --name=%s --version=%s \\\n", artifact.Name, artifact.Version)
		fmt.Printf("  %s \\\n", artifact.MetalinkFile().URLs[0].URL)
		fmt.Printf("  --sha1=%s\n", strings.TrimPrefix(artifact.PreferredChecksum().String(), "sha1:"))

		return nil
	}

	cmd := exec.Command(
		"bosh",
		"upload-release",
		fmt.Sprintf("--name=%s", artifact.Name),
		fmt.Sprintf("--version=%s", artifact.Version),
		artifact.MetalinkFile().URLs[0].URL,
		fmt.Sprintf("--sha1=%s", strings.TrimPrefix(artifact.PreferredChecksum().String(), "sha1:")),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
