// Package install provides functions to recursively install dependencies
// Cleaning of existing changes in dependency repositories is controlled by the --clean flag
package install

// Copyright 2013 Vubeology, Inc.

import (
	"flag"
	"fmt"
	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/timelock"
	"github.com/vube/depman/util"
	"time"
)

var (
	clean bool
)

// Whether to install recursively
var Recurse = true

func init() {
	flag.BoolVar(&clean, "clean", false, "Remove changes to code in dependencies")
}

// Install a DependencyMap
func Install(deps dep.DependencyMap) int {
	util.Print(colors.Blue("Installing:"))
	set := make(map[string]string)
	return recursiveInstall(deps, set)
}

// recursively install a DependencyMap
func recursiveInstall(deps dep.DependencyMap, set map[string]string) (result int) {
	for name, d := range deps.Map {
		start := time.Now()

		if duplicate(*d, set) {
			continue
		}

		stale := timelock.IsStale(d.Repo)

		util.PrintDep(name, d.Version, d.Repo, stale)

		subPath := d.Path()
		result += d.VCS.Clone(d)
		result += util.Cd(subPath)

		if clean {
			d.VCS.Clean(d)
		}

		if stale {
			util.VerboseIndent(" # repo is stale, pulling")
			result += d.VCS.Pull(d)
		}

		result += d.VCS.Checkout(d)

		util.VerboseIndent(fmt.Sprintf("# time to install: %.3fs", time.Since(start).Seconds()))

		// Recursive
		depsFile := util.UpwardFind(subPath, dep.DepsFile)
		if depsFile != "" && Recurse {
			subDeps, err := dep.Read(depsFile)
			if err != nil {
				util.Print(colors.Yellow("Error reading deps from '" + subDeps.Path + "': " + err.Error()))
			} else {
				util.IncreaseIndent()
				result += recursiveInstall(subDeps, set)
				util.DecreaseIndent()
			}
		}
	}
	return
}

// Check for duplicate dependency
// if same name and same version, skip
// if same name and different version, exit
// if different name, add to set, don't skip
func duplicate(d dep.Dependency, set map[string]string) (skip bool) {
	version, installed := set[d.Repo]
	if installed && version != d.Version {
		util.Print(colors.Red("ERROR    : Duplicate dependency with different versions detected"))
		util.Print(colors.Red("Repo     : " + d.Repo))
		util.Fatal(colors.Red("Versions : " + d.Version + "\t" + version))
	} else if installed {
		util.VerboseIndent(colors.Yellow("Skipping previously installed dependency: ") + d.Repo)
		skip = true
	} else {
		set[d.Repo] = d.Version
	}
	return
}

// Mock sets clean to true for testing
func Mock() {
	clean = true
}
