// Package add implements functions for adding a dependency
package add

// Copyright 2013 Vubeology, Inc.

import (
	"fmt"
	"strings"
	"vube/depman/colors"
	"vube/depman/dep"
	"vube/depman/install"
	"vube/depman/util"
)

// Add interactively prompts the user for details of a dependency, adds it to deps.json, and writes out the file
func Add(deps dep.DependencyMap, name string) (result int) {
	var cont = true
	_, exists := deps.Map[name]
	if exists {
		util.Fatal(colors.Red("Dependency '" + name + "'' is already defined, pick another name."))
	}

	util.Print(colors.Blue("Adding: ") + name)

	for cont {
		d := dep.Dependency{}
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
		cont = promptBool("Add another", "y/N")
	}

	err := deps.Write()
	if err != nil {
		util.Fatal(colors.Red("Error Writing " + deps.Path + ": " + err.Error()))
	}

	result = install.Install(deps)

	return
}

// Promt the user with question and return a string answer
func promptString(question string, details string) (answer string) {
	fmt.Print(colors.Blue(question) + " (" + details + "): ")
	fmt.Scanln(&answer)
	return
}

// Prompt the user with question and return a bool answer
func promptBool(question string, details string) (answer bool) {
	str := promptString(question, details)

	switch str {
	case "y", "yes":
		answer = true
	default:
		answer = false
	}

	return
}

// Prompt the user with question check that the answer is a valid dep type and then return it
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
