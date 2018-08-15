package cli

type Cmd struct {
	PropertiesCmd PropertiesCmd `command:"properties" description:"Show the job properties"`
	SpecCmd       SpecCmd       `command:"spec" description:"Show the job or release manifests"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.SpecCmd.Execute(extra)
}
