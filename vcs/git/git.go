package git

// Copyright 2013 Vubeology, Inc.

import (
	"os/exec"
	"strings"
	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/util"
)

// LastCommit retrieves the version number of the last commit on branch
// Assumes that the current working directory is in the git repo
func LastCommit(d dep.Dependency, branch string) (hash string) {
	c := exec.Command("git", "log", "-1", "--format=%h")
	out, err := c.CombinedOutput()

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("git log -1 --format=%h"))
		util.PrintIndent(colors.Red(string(out)))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	hash = strings.Replace(string(out), "\n", "", -1)
	return
}

//GetHead - Render a revspec to a commit ID
func GetHead(d dep.Dependency) (hash string, err error) {
	var pwd string

	pwd = util.Pwd()
	util.Cd(d.Path())

	c := exec.Command("git", "rev-parse", d.Version)
	{
		var out_bytes []byte
		out_bytes, err = c.CombinedOutput()
		hash = strings.TrimSuffix(string(out_bytes), "\n")
	}

	util.Cd(pwd)

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("git rev-parse " + d.Version))
		util.PrintIndent(colors.Red(string(hash)))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	return
}

// IsBranch determines if a version (branch, commit hash, tag) is a branch (i.e. can we pull from the remote).
// Assumes we are already in a sub directory of the repo
func IsBranch(name string) (result bool) {
	c := exec.Command("git", "branch", "-r")
	out, err := c.CombinedOutput()

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("git branch -r"))
		util.PrintIndent(colors.Red(string(out)))
		util.PrintIndent(colors.Red(err.Error()))
		return false
	}

	// get the string version but also strip the trailing newline
	stringOut := string(out[0 : len(out)-1])

	lines := strings.Split(stringOut, "\n")
	for _, val := range lines {
		// for "origin/HEAD -> origin/master"
		arr := strings.Split(val, " -> ")
		remoteBranch := arr[0]

		// for normal "origin/develop"
		arr = strings.Split(remoteBranch, "/")
		branch := arr[1]
		if branch == name {
			return true
		}
	}

	return
}

// CloneFetch will clone d.Repo into d.Path() if d.Path does not exist, otherwise it will cd to d.Path() and run git fetch
func CloneFetch(d dep.Dependency) (result int) {
	if !util.Exists(d.Path()) {
		result = util.RunCommand("git clone " + d.Repo + " " + d.Path())
	} else {
		result = util.Cd(d.Path())
		result += util.RunCommand("git fetch")
	}
	return
}
