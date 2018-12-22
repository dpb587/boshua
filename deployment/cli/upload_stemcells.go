package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"code.cloudfoundry.org/workpool"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/deployment/manifest"
	"github.com/dpb587/boshua/osversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UploadStemcellsCmd struct {
	setter.AppConfig `no-flag:"true"`

	IaaS       string `long:"stemcell-iaas" description:"The stemcell IaaS"`
	Hypervisor string `long:"stemcell-hypervisor" description:"The stemcell hypervisor"`
	Flavor     string `long:"stemcell-flavor" description:"The stemcell flavor (e.g. 'light')"`

	Parallel int `long:"parallel" description:"Maximum number of parallel operations" default:"3"`

	clicommon.UploadStemcellCmd
}

func (c *UploadStemcellsCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "deployment/upload-stemcells"})

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "reading stdin")
	}

	man, err := manifest.Parse(bytes, osversion.Reference{})
	if err != nil {
		return errors.Wrap(err, "parsing manifest")
	}

	index, err := c.Config.GetStemcellIndex("default")
	if err != nil {
		return errors.Wrap(err, "loading index")
	}

	var parallelize []func()

	requirements := man.StemcellRequirements()

	for reqIdx := range requirements {
		req := requirements[reqIdx]

		parallelize = append(parallelize, func() {
			f := req.FilterParams()

			if f.IaaS == "" && c.IaaS != "" {
				f.IaaSExpected = true
				f.IaaS = c.IaaS
			}

			if f.Hypervisor == "" && c.Hypervisor != "" {
				f.HypervisorExpected = true
				f.Hypervisor = c.Hypervisor
			}

			if f.Flavor == "" && c.Flavor != "" {
				f.FlavorExpected = true
				f.Flavor = c.Flavor
			}

			parallelLog := func(msg string) {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s] %s\n", time.Now().Format("15:04:05"), req.Slug(), msg))
			}

			result, err := stemcellversiondatastore.GetArtifact(index, f)
			if err != nil {
				parallelLog(errors.Wrap(err, "skipped: error: getting stemcell artifact").Error())

				return
			}

			err = c.UploadStemcellCmd.ExecuteArtifact(
				c.Config.GetDownloader,
				func() (artifact.Artifact, error) {
					return result, nil
				},
				clicommon.UploadStemcellOpts{},
			)
			if err != nil {
				log.Fatalf(fmt.Errorf("skipped: error: uploading stemcell: %v", err).Error())

				return
			}

			parallelLog("uploaded stemcell")
		})
	}

	pool, err := workpool.NewThrottler(c.Parallel, parallelize)
	if err != nil {
		log.Fatalf(fmt.Errorf("parallelizing: %v", err).Error())
	}
	pool.Work()

	return nil
}
