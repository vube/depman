// Package update provides functions to update a dependency to the latest version
package update

// Copyright 2013 Vubeology, Inc.

import (
	"vube/depman/colors"
	"vube/depman/dep"
	"vube/depman/install"
	"vube/depman/util"
	"vube/depman/vcs"
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
	vcs.Checkout(d, false)
	v, err := vcs.LastCommit(d, branch)
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
