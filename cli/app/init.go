package app

import "github.com/dpb587/boshua/cli"

var name, semver, commit, built string

var App cli.App

func init() {
	App = cli.MustApp(name, semver, commit, built)
}
