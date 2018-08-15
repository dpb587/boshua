package main

import (
	"os"

	"github.com/dpb587/boshua/cli/app"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/main/boshua/cmd"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var defaultServer string

func main() {
	c := cmd.New(app.App)
	c.Opts.DefaultServer = defaultServer

	var parser = flags.NewParser(c, flags.Default)
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if command == nil {
			// TODO something more specific?
			return nil
		}

		if cs, ok := command.(setter.Setter); ok {
			config, err := c.Opts.GetConfig()
			if err != nil {
				return errors.Wrap(err, "loading config")
			}

			cs.SetConfig(config)
		}

		return command.Execute(args)
	}

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
