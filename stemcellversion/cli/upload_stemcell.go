package cli

import (
	"fmt"
	"os"
	"os/exec"
)

type UploadStemcellCmd struct {
	*CmdOpts `no-flag:"true"`

	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

func (c *UploadStemcellCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/upload-stemcell")

	artifact, err := c.getStemcell()
	if err != nil {
		return fmt.Errorf("finding compiled stemcell: %v", err)
	}

	if c.Cmd {
		fmt.Printf("bosh upload-stemcell --name=%s --version=%s \\\n", c.StemcellOpts.Stemcell.Name, c.StemcellOpts.Stemcell.Version)
		fmt.Printf("  %s \\\n", artifact.ArtifactMetalinkFile().URLs[0].URL)
		// fmt.Printf("  --sha1=%s\n", strings.TrimPrefix(c.StemcellOpts.StemcellChecksum.ImmutableChecksum.String(), "sha1:")) // TODO sha1-find

		return nil
	}

	cmd := exec.Command(
		"bosh",
		"upload-stemcell",
		fmt.Sprintf("--name=%s", c.StemcellOpts.Stemcell.Name),
		fmt.Sprintf("--version=%s", c.StemcellOpts.Stemcell.Version),
		artifact.ArtifactMetalinkFile().URLs[0].URL,
		// fmt.Sprintf("--sha1=%s", strings.TrimPrefix(c.StemcellOpts.StemcellChecksum.ImmutableChecksum.String(), "sha1:")), // TODO sha1-find
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
