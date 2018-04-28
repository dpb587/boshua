package compiledrelease

import (
	"fmt"
	"log"

	"github.com/dpb587/metalink"
)

type MetalinkCmd struct {
	*CmdOpts `no-flag:"true"`

	Format string `long:"format" description:"Output format for metalink"`
}

func (c *MetalinkCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiled-release/metalink")

	resInfo, err := c.CompiledReleaseOpts.GetCompiledReleaseVersion(c.AppOpts.GetClient())
	if err != nil {
		log.Fatalf("requesting compiled version info: %v", err)
	} else if resInfo == nil {
		log.Fatalf("no compiled release available")
	}

	meta4 := metalink.Metalink{
		Files: []metalink.File{
			resInfo.Data,
		},
		Generator: "bosh-compiled-releases/0.0.0",
	}

	meta4Bytes, err := metalink.Marshal(meta4)
	if err != nil {
		log.Fatalf("marshalling response: %v", err)
	}

	fmt.Printf("%s\n", meta4Bytes)

	return nil
}
