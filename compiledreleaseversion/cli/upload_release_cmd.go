package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
)

type UploadReleaseCmd struct {
	*CmdOpts `no-flag:"true"`

	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/upload-release")

	resInfo, err := c.getCompiledRelease()
	if err != nil {
		log.Fatalf("requesting compiled version info: %v", err)
	} else if resInfo == nil {
		log.Fatalf("no compiled release available")
	}

	var sha1 checksum.Checksum

	for _, cs := range metalinkutil.HashesToChecksums(resInfo.Data.Artifact.Hashes) {
		if cs.Algorithm().Name() == "sha1" {
			sha1 = &cs

			break
		}
	}

	if c.Cmd {
		fmt.Printf("bosh upload-release --name=%s --version=%s \\\n", c.CompiledReleaseOpts.Release.Name, c.CompiledReleaseOpts.Release.Version)
		fmt.Printf("  %s \\\n", resInfo.Data.Artifact.URLs[0].URL)
		fmt.Printf("  --sha1=%s \\\n", strings.TrimPrefix(sha1.String(), "sha1:"))
		fmt.Printf("  --stemcell=%s/%s\n", c.CompiledReleaseOpts.OS.Name, c.CompiledReleaseOpts.OS.Version)

		return nil
	}

	cmd := exec.Command(
		"bosh",
		"upload-release",
		fmt.Sprintf("--name=%s", c.CompiledReleaseOpts.Release.Name),
		fmt.Sprintf("--version=%s", c.CompiledReleaseOpts.Release.Version),
		resInfo.Data.Artifact.URLs[0].URL,
		fmt.Sprintf("--sha1=%s \\\n", strings.TrimPrefix(sha1.String(), "sha1:")),
		fmt.Sprintf("--stemcell=%s/%s\n", c.CompiledReleaseOpts.OS.Name, c.CompiledReleaseOpts.OS.Version),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
