package compiledrelease

import (
	"fmt"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type OpsFileCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *OpsFileCmd) Execute(_ []string) error {
	resInfo, err := c.CompiledReleaseOpts.GetCompiledReleaseVersion(c.AppOpts.GetClient())
	if err != nil {
		log.Fatalf("requesting compiled version info: %v", err)
	} else if resInfo == nil {
		log.Fatalf("no compiled release available")
	}

	opsBytes, err := yaml.Marshal([]map[string]interface{}{
		{
			"path": fmt.Sprintf("/releases/name=%s?", resInfo.Data.Release.Name),
			"type": "replace",
			"value": map[string]interface{}{
				"name":    resInfo.Data.Release.Name,
				"version": resInfo.Data.Release.Version,
				"sha1":    strings.TrimPrefix(string(resInfo.Data.Tarball.Checksums[0]), "sha1:"), // TODO .Preferred()
				"url":     resInfo.Data.Tarball.URL,
				"stemcell": map[string]string{
					"os":      resInfo.Data.Stemcell.OS,
					"version": resInfo.Data.Stemcell.Version,
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
