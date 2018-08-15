package cli

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type App struct {
	Name   string
	Semver string
	Commit string
	Built  time.Time
}

func MustApp(name, semver, commit, built string) App {
	app, err := NewApp(name, semver, commit, built)
	if err != nil {
		panic(err)
	}

	return app
}

func NewApp(name, semver, commit, built string) (App, error) {
	if name == "" {
		name = "unknown"
	}

	if semver == "" {
		semver = "0.0.0+dev"
	}

	if commit == "" {
		commit = "unknown"
	}

	var builtval time.Time
	var err error

	if built == "" {
		builtval = time.Now()
	} else {
		builtval, err = time.Parse(time.RFC3339, built)
		if err != nil {
			return App{}, errors.Wrap(err, "failed to parse version time")
		}
	}

	return App{name, semver, commit, builtval}, nil
}

func (a App) Slug() string {
	return fmt.Sprintf("%s/%s", a.Name, a.Semver)
}

func (a App) String() string {
	return fmt.Sprintf("%s/%s (commit %s, built %s)", a.Name, a.Semver, a.Commit, a.Built.Format(time.RFC3339))
}
