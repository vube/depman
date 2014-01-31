// Package update provides functions to update a dependency to the latest version
package update

// Copyright 2013 Vubeology, Inc.

import (
	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/install"
	"github.com/vube/depman/util"
)

// Update rewrites Dependency name in deps.json to use the last commit in branch as version
func Update(deps dep.DependencyMap, name string, branch string) {
	util.Print(colors.Blue("Updating:"))

	d, ok := deps.Map[name]
	if !ok {
		util.Fatal(colors.Red("Dependency Name '" + name + "' not found in deps.json"))
	}

	oldVersion := d.Version

	pwd := util.Pwd()
	util.Cd(d.Path())
	d.VCS.Checkout(d)
	d.VCS.Pull(d)
	v, err := d.VCS.LastCommit(d, branch)
	if err != nil {
		util.Fatal(err)
	}
	d.Version = v

	util.PrintIndent(colors.Blue(name) + " (" + oldVersion + " --> " + d.Version + ")")

	util.Cd(pwd)
	deps.Map[name] = d
	deps.Write()

	install.Install(deps)

}
