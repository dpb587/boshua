package compiledrelease

import (
	"fmt"
	"log"
	"strings"

	"github.com/dpb587/boshua/util/metalinkutil"
	yaml "gopkg.in/yaml.v2"
)

type OpsFileCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *OpsFileCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiled-release/ops-file")

	resInfo, err := c.getCompiledRelease()
	if err != nil {
		log.Fatalf("requesting compiled version info: %v", err)
	} else if resInfo == nil {
		log.Fatalf("no compiled release available")
	}

	opsBytes, err := yaml.Marshal([]map[string]interface{}{
		{
			"path": fmt.Sprintf("/releases/name=%s?", c.CompiledReleaseOpts.Release.Name),
			"type": "replace",
			"value": map[string]interface{}{
				"name":    c.CompiledReleaseOpts.Release.Name,
				"version": c.CompiledReleaseOpts.Release.Version,
				"sha1":    strings.TrimPrefix(metalinkutil.HashToChecksum(resInfo.Data.Artifact.Hashes[0]).String(), "sha1:"), // TODO .Preferred()
				"url":     resInfo.Data.Artifact.URLs[0].URL,
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
