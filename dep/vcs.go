// Package vcs provides functions to various Version Control Systems
package vcs

// Copyright 2013 Vubeology, Inc.

import (
	"errors"
	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/util"
	"github.com/vube/depman/vcs/bzr"
	"github.com/vube/depman/vcs/git"
	"github.com/vube/depman/vcs/hg"
)

// Checkout uses the appropriate VCS to checkout the specified version of the code
func Checkout(d dep.Dependency, clean bool) (result int) {
	switch d.Type {
	case dep.TypeGit, dep.TypeGitClone:
		if clean {
			util.PrintIndent(colors.Red("Cleaning:") + colors.Blue(" git reset --hard HEAD"))
			util.RunCommand("git reset --hard HEAD")
			util.RunCommand("git clean -f")
		}
		if util.RunCommand("git checkout " + d.Version) != 0 {
			util.RunCommand("git fetch")
			util.RunCommand("git checkout " + d.Version)
		}

		if git.IsBranch(d.Version) {
			util.RunCommand("git pull origin " + d.Version)
		}
	case dep.TypeBzr:
		util.RunCommand("bzr up --revision " + d.Version)
	case dep.TypeHg:
		if clean {
			util.PrintIndent(colors.Red("Cleaning:") + colors.Blue(" hg up --clean "+d.Version))
			util.RunCommand("hg up --clean " + d.Version)
		} else {
			util.RunCommand("hg up " + d.Version)
		}
	default:
		util.PrintIndent(colors.Red(d.Repo + ": Unknown repository type (" + d.Type + "), skipping..."))
		util.PrintIndent(colors.Red("Valid Repository types: " + dep.TypeGit + ", " + dep.TypeHg + ", " + dep.TypeBzr + ", " + dep.TypeGitClone))
		result += 1
	}
	return
}

// LastCommit uses the appropriate VCS to retrieve the version number of the last commit on branch
func LastCommit(d dep.Dependency, branch string) (hash string, err error) {
	switch d.Type {
	case dep.TypeGit, dep.TypeGitClone:
		if !git.IsBranch(branch) {
			err = errors.New("Branch '" + branch + "' is not a valid branch")
		}
		hash = git.LastCommit(d, branch)
	case dep.TypeBzr:
		hash = bzr.LastCommit(d, branch)
	case dep.TypeHg:
		hash = hg.LastCommit(d, branch)
	default:
		util.PrintIndent(colors.Red(d.Repo + ": Unknown repository type (" + d.Type + ")"))
		util.PrintIndent(colors.Red("Valid Repository types: " + dep.TypeGit + ", " + dep.TypeHg + ", " + dep.TypeBzr + ", " + dep.TypeGitClone))
		util.Fatal("")
	}
	return
}

//GetHead - Render a revspec to a commit ID
func GetHead(d dep.Dependency) (to_return string, err error) {
	switch d.Type {
		case dep.TypeGit, dep.TypeGitClone:
			to_return, err = git.GetHead(d)
		case dep.TypeHg:
			to_return, err = hg.GetHead(d)
		case dep.TypeBzr:
			to_return, err = bzr.GetHead(d)
		default:
			util.PrintIndent(colors.Red(d.Repo + ": Unknown repository type (" + d.Type + ")"))
			util.PrintIndent(colors.Red("Valid Repository types: " + dep.TypeGit + ", " + dep.TypeHg + ", " + dep.TypeBzr + ", " + dep.TypeGitClone))
			util.Fatal("")
	}
	return
}
