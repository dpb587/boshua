package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ExportReleaseCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	Args ExportReleaseCmdArgs `positional-args:"true" required:"true"`
}

type ExportReleaseCmdArgs struct {
	Local string `positional-arg-name:"PATH" description:"Path to save the exported release"`
}

func (c *ExportReleaseCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/compilation/export-release"})

	manifestFile, err := ioutil.TempFile("", "boshua-export-release")
	if err != nil {
		return errors.Wrap(err, "creating deployment manifest")
	}

	defer os.Remove(manifestFile.Name())

	deploymentName := fmt.Sprintf(
		"%s-%s-on-%s-stemcell-%s",
		c.CompiledReleaseOpts.ReleaseOpts.Name,
		c.CompiledReleaseOpts.ReleaseOpts.Version,
		c.CompiledReleaseOpts.StemcellOpts.OS,
		c.CompiledReleaseOpts.StemcellOpts.Version,
	)

	_, err = manifestFile.Write([]byte(fmt.Sprintf(`name: "%s"
instance_groups: []
releases:
- name: "%s"
  version: "%s"
stemcells:
- alias: default
  os: "%s"
  version: "%s"
update:
  canaries: 1
  canary_watch_time: 1
  max_in_flight: 1
  update_watch_time: 1
`,
		deploymentName,
		c.CompiledReleaseOpts.ReleaseOpts.Name,
		c.CompiledReleaseOpts.ReleaseOpts.Version,
		c.CompiledReleaseOpts.StemcellOpts.OS,
		c.CompiledReleaseOpts.StemcellOpts.Version,
	)))
	if err != nil {
		return errors.Wrap(err, "writing deployment manifest")
	}

	err = manifestFile.Close()
	if err != nil {
		return errors.Wrap(err, "closing deployment manifest")
	}

	cmd := exec.Command("bosh", "-n", fmt.Sprintf("-d=%s", deploymentName), "deploy", manifestFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "deploying")
	}

	cmd = exec.Command(
		"bosh",
		fmt.Sprintf("-d=%s", deploymentName),
		"export-release",
		fmt.Sprintf("%s/%s", c.CompiledReleaseOpts.ReleaseOpts.Name, c.CompiledReleaseOpts.ReleaseOpts.Version),
		fmt.Sprintf("%s/%s", c.CompiledReleaseOpts.StemcellOpts.OS, c.CompiledReleaseOpts.StemcellOpts.Version),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "exporting release")
	}

	exportedPath, err := filepath.Glob("*.tgz")
	if err != nil {
		return errors.Wrap(err, "finding exported release")
	} else if len(exportedPath) != 1 {
		return errors.New("expected an exported release tarball")
	}

	// TODO ew
	sys := boshsys.NewOsFileSystem(boshlog.NewLogger(boshlog.LevelNone))
	err = sys.CopyFile(exportedPath[0], c.Args.Local) // TODO Rename failed for some reason
	if err != nil {
		return errors.Wrap(err, "moving tarball")
	}

	err = os.Remove(exportedPath[0])
	if err != nil {
		return errors.Wrap(err, "removing initial tarball")
	}

	cmd = exec.Command("bosh", "-n", fmt.Sprintf("-d=%s", deploymentName), "delete-deployment")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "deleting deployment")
	}

	return nil
}
