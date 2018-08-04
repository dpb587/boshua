package clicommon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
)

type UploadReleaseCmd struct {
	Cmd bool `long:"cmd" description:"Show the command instead of running it"`
}

type UploadReleaseOpts struct {
	Name     string
	Version  string
	Stemcell string
}

func (c *UploadReleaseCmd) ExecuteArtifact(loader artifact.Loader, opts UploadReleaseOpts) error {
	artifact, err := loader()
	if err != nil {
		log.Fatal(err)
	}

	// TODO verify URLs[0]
	url := artifact.MetalinkFile().URLs[0].URL

	var idArgs = []string{}
	var args = []string{}

	for _, cs := range metalinkutil.HashesToChecksums(artifact.MetalinkFile().Hashes) {
		if cs.Algorithm().Name() == "sha1" {
			args = append(args, fmt.Sprintf("--sha1=%s", strings.TrimPrefix(cs.String(), "sha1:")))

			break
		}
	}

	var effectiveOpts UploadReleaseOpts

	switch artifact := artifact.(type) {
	case releaseversion.Artifact:
		effectiveOpts.Name = artifact.Name
		effectiveOpts.Version = artifact.Version
	case compilation.Artifact:
		artifactRef := artifact.Reference().(compilation.Reference)

		effectiveOpts.Name = artifactRef.ReleaseVersion.Name
		effectiveOpts.Version = artifactRef.ReleaseVersion.Version
		effectiveOpts.Stemcell = fmt.Sprintf("%s/%s", artifactRef.OSVersion.Name, artifactRef.OSVersion.Version)
	}

	if opts.Name != "" {
		effectiveOpts.Name = opts.Name
	}

	if opts.Version != "" {
		effectiveOpts.Version = opts.Version
	}

	if opts.Stemcell != "" {
		effectiveOpts.Stemcell = opts.Stemcell
	}

	if effectiveOpts.Name != "" {
		idArgs = append(idArgs, fmt.Sprintf("--name=%s", effectiveOpts.Name))
	}

	if effectiveOpts.Version != "" {
		idArgs = append(idArgs, fmt.Sprintf("--version=%s", effectiveOpts.Version))
	}

	if effectiveOpts.Stemcell != "" {
		args = append(args, fmt.Sprintf("--stemcell=%s", effectiveOpts.Stemcell))
	}

	if c.Cmd {
		fmt.Printf("bosh upload-release %s", strings.Join(idArgs, " "))
		fmt.Printf(" \\\n  %s", url)

		for _, arg := range args {
			fmt.Printf(" \\\n  %s", arg)
		}

		fmt.Printf("\n")

		return nil
	}

	cmd := exec.Command("bosh", append([]string{"upload-release"}, append(append(idArgs, url), args...)...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
