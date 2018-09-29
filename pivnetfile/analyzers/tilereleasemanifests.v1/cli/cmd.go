package cli

type Cmd struct {
	ReleasesCmd ReleasesCmd `command:"releases" description:"Show releases in the tile"`
}
