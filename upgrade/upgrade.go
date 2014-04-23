package upgrade

// Copyright 2013-2014 Vubeology, Inc.

import (
	"github.com/vube/depman/dep"
	"github.com/vube/depman/install"
	"github.com/vube/depman/util"
)

// Self upgrades this version of depman to the latest on the master branch
func Self() {
	deps := dep.New()
	d := new(dep.Dependency)
	d.Repo = "github.com/vube/depman"
	d.Version = "master"
	d.Type = "git"
	d.SetupVCS("depman")

	deps.Map["depman"] = d

	install.Recurse = false
	install.Install(deps)
	install.Recurse = true
	util.RunCommand("go install github.com/vube/depman")
}
