package cli

type Cmd struct {
	LsCmd        LsCmd        `command:"ls" description:"Show an ls-style list of files"`
	Sha1sumCmd   Sha1sumCmd   `command:"sha1sum" alias:"shasum" description:"Show sha1 checksums"`
	Sha256sumCmd Sha256sumCmd `command:"sha256sum" description:"Show sha256 checksums"`
	Sha512sumCmd Sha512sumCmd `command:"sha512sum" description:"Show sha512 checksums"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.LsCmd.Execute(extra)
}
