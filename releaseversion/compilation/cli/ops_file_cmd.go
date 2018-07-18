package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type OpsFileCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *OpsFileCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiled-release/ops-file")

	artifact, err := c.CompiledReleaseOpts.Artifact()
	if err != nil {
		return errors.Wrap(err, "finding compiled release")
	}

	var releaseName, releaseVersion string

	if c.CompiledReleaseOpts.ReleaseOpts.NameVersion != nil {
		nv := c.CompiledReleaseOpts.ReleaseOpts.NameVersion

		releaseName = nv.Name
		releaseVersion = nv.Version
	} else {
		releaseName = c.CompiledReleaseOpts.ReleaseOpts.Name
		releaseVersion = c.CompiledReleaseOpts.ReleaseOpts.Version
	}

	opsBytes, err := yaml.Marshal([]map[string]interface{}{
		{
			"path": fmt.Sprintf("/releases/name=%s?", releaseName),
			"type": "replace",
			"value": map[string]interface{}{
				"name":    releaseName,
				"version": releaseVersion,
				"sha1":    strings.TrimPrefix(metalinkutil.HashToChecksum(artifact.MetalinkFile().Hashes[0]).String(), "sha1:"), // TODO .Preferred()
				"url":     artifact.MetalinkFile().URLs[0].URL,
				"stemcell": map[string]string{
					"os":      c.CompiledReleaseOpts.OS.Name,
					"version": c.CompiledReleaseOpts.OS.Version,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("marshaling ops: %v", err)
	}

	fmt.Printf("%s\n", opsBytes)

	return nil
}
