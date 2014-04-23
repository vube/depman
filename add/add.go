// Package add implements functions for interactively adding a dependency
package add

// Copyright 2014 Vubeology, Inc.

import (
	"fmt"
	"strings"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/install"
	"github.com/vube/depman/util"
)

// Add interactively prompts the user for details of a dependency, adds it to deps.json, writes out the file, and installs the dependencies
func Add(deps dep.DependencyMap, name string) {
	_, exists := deps.Map[name]
	if exists {
		util.Fatal(colors.Red("Dependency '" + name + "'' is already defined, pick another name."))
	}

	util.Print(colors.Blue("Adding: ") + name)

	d := new(dep.Dependency)
	d.Type = promptType("Type", "git, git-clone, hg, bzr")

	if d.Type == dep.TypeGitClone {
		d.Repo = promptString("Repo", "git url")
	} else {
		d.Repo = promptString("Repo", "go import")
	}

	d.Version = promptString("Version", "hash, branch, or tag")

	if d.Type == dep.TypeGitClone {
		d.Alias = promptString("Alias", "where to install the repo")
	}

	deps.Map[name] = d

	for name, d := range deps.Map {
		err := d.SetupVCS(name)
		if err != nil {
			delete(deps.Map, name)
		}
	}

	err := deps.Write()
	if err != nil {
		util.Fatal(colors.Red("Error Writing " + deps.Path + ": " + err.Error()))
	}

	install.Install(deps)

	return
}

// promptString prompts the user with question and returns a string answer
func promptString(question string, details string) (answer string) {
	fmt.Print(colors.Blue(question) + " (" + details + "): ")
	fmt.Scanln(&answer)
	return
}

// promptType prompts the user with question, checks that the answer is a valid dep type and then return it
func promptType(question string, details string) (t string) {
	for {
		t = promptString(question, details)
		t = strings.TrimSpace(t)
		switch t {
		case dep.TypeBzr, dep.TypeGit, dep.TypeHg, dep.TypeGitClone:
			return
		default:
			util.Print(colors.Red("Invalid Type, try again..."))
		}
	}
}
