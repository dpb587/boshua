package cli

type Cmd struct {
	ContentsCmd ContentsCmd `command:"contents" description:"Show the original contents of stemcell.MF"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ContentsCmd.Execute(extra)
}
