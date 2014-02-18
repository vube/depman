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

	// record the old version
	oldVersion := d.Version

	// temporarily use the branch
	d.Version = branch

	pwd := util.Pwd()
	util.Cd(d.Path())
	d.VCS.Checkout(d)
	d.VCS.Update(d)

	// get the last commit on the newly checked out branch
	v, err := d.VCS.LastCommit(d, branch)
	if err != nil {
		util.Fatal(err)
	}

	// set the version to be the last commit
	d.Version = v

	util.PrintIndent(colors.Blue(name) + " (" + oldVersion + " --> " + d.Version + ")")

	util.Cd(pwd)
	deps.Map[name] = d
	deps.Write()

	install.Install(deps)

}
