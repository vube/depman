package upgrade

// Copyright 2013 Vubeology, Inc.

import (
	"vube/depman/dep"
	"vube/depman/install"
	"vube/depman/util"
)

// Self upgrades this version of depman to the latest on the master branch
func Self() {
	deps := dep.New()
	d := dep.Dependency{}
	d.Repo = "vube/depman"
	d.Version = "master"
	d.Type = "git"

	deps.Map["depman"] = d

	install.Recurse = false
	install.Install(deps)
	install.Recurse = true
	util.RunCommand("go install vube/depman")
}
