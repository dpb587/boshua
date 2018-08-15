package cli

type Cmd struct {
	PackagesCmd PackagesCmd `command:"packages" description:"Show a simple list of package versions"`
	ContentsCmd ContentsCmd `command:"contents" description:"Show the original contents of packages.txt"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.PackagesCmd.Execute(extra)
}
